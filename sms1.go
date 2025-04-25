package main

import (
    "context"
    "errors"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "os"
    "os/signal"
    "regexp"
    "strconv"
    "strings"
    "sync"
    "syscall"
    "time"
)

// OTPService defines configuration for an OTP endpoint
type OTPService struct {
    Name        string
    URL         string
    ContentType string
    BuildBody   func(phone string) (io.Reader, error)
    Rate        int           // requests per second
    Timeout     time.Duration // per-request timeout
    MaxRetries  int           // number of retry attempts
    Backoff     time.Duration // initial backoff duration
}

func main() {
    // Command-line flags
    phone := flag.String("phone", "", "Phone number in format 09xxxxxxxxx (required)")
    count := flag.Int("n", 1, "Number of OTP requests per service")
    flag.Parse()

    if *phone == "" {
        fmt.Println("Error: -phone flag is required")
        flag.Usage()
        os.Exit(1)
    }
    if err := validatePhone(*phone); err != nil {
        log.Fatalf("invalid phone: %v", err)
    }

    // Define services with rate-limits, timeouts, retries, backoff
    services := []OTPService{
        {
            Name:        "SnappFood",
            URL:         "https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass",
            ContentType: "application/x-www-form-urlencoded",
            BuildBody: func(phone string) (io.Reader, error) {
                v := url.Values{}
                v.Set("cellphone", phone)
                v.Set("lat", "35.774")
                v.Set("long", "51.418")
                v.Set("client", "WEBSITE")
                return strings.NewReader(v.Encode()), nil
            },
            Rate:       1,
            Timeout:    5 * time.Second,
            MaxRetries: 5,
            Backoff:    2 * time.Second,
        },
        {
            Name:        "Mobinnet",
            URL:         "https://my.mobinnet.ir/api/account/SendRegisterVerificationCode",
            ContentType: "application/json",
            BuildBody: func(phone string) (io.Reader, error) {
                payload := fmt.Sprintf(`{"cellNumber":"%s"}`, phone)
                return strings.NewReader(payload), nil
            },
            Rate:       1,
            Timeout:    15 * time.Second,
            MaxRetries: 3,
            Backoff:    3 * time.Second,
        },
    }

    // Shared HTTP client (no default timeout)
    client := &http.Client{}

    // Cancellation context
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    var wg sync.WaitGroup

    for _, svc := range services {
        ticker := time.NewTicker(time.Second / time.Duration(svc.Rate))
        defer ticker.Stop()
        for i := 0; i < *count; i++ {
            wg.Add(1)
            go func(s OTPService, iter int) {
                defer wg.Done()
                <-ticker.C
                sendWithRetry(ctx, client, s, *phone, iter)
            }(svc, i+1)
        }
    }

    wg.Wait()
    log.Println("All done.")
}

// sendWithRetry attempts the request with retries, backoff, and handles 429
func sendWithRetry(ctx context.Context, client *http.Client, s OTPService, phone string, iter int) {
    backoff := s.Backoff
    for attempt := 1; attempt <= s.MaxRetries; attempt++ {
        // per-request timeout
        reqCtx, cancel := context.WithTimeout(ctx, s.Timeout)
        bodyReader, err := s.BuildBody(phone)
        if err != nil {
            log.Printf("[%s][%d.%d] build body error: %v", s.Name, iter, attempt, err)
            cancel()
            return
        }
        req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, s.URL, bodyReader)
        if err != nil {
            log.Printf("[%s][%d.%d] new request error: %v", s.Name, iter, attempt, err)
            cancel()
            return
        }
        req.Header.Set("Content-Type", s.ContentType)

        resp, err := client.Do(req)
        if err != nil {
            log.Printf("[%s][%d.%d] do request error: %v", s.Name, iter, attempt, err)
            cancel()
            time.Sleep(backoff)
            backoff *= 2
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            log.Printf("[%s][%d] success on attempt %d, status %d", s.Name, iter, attempt, resp.StatusCode)
            cancel()
            return
        }
        if resp.StatusCode == 429 {
            // handle rate-limit
            ra := resp.Header.Get("Retry-After")
            secs, err := strconv.Atoi(ra)
            cancel()
            if err == nil {
                log.Printf("[%s][%d.%d] 429 received, retry after %d seconds", s.Name, iter, attempt, secs)
                time.Sleep(time.Duration(secs) * time.Second)
            } else {
                log.Printf("[%s][%d.%d] 429 received, backoff %v", s.Name, iter, attempt, backoff)
                time.Sleep(backoff)
                backoff *= 2
            }
            continue
        }
        log.Printf("[%s][%d] unexpected status on attempt %d: %d", s.Name, iter, attempt, resp.StatusCode)
        cancel()
        re
