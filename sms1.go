package main

import (
    "context"
    "errors"
    "flag"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"
    "os/signal"
    "regexp"
    "sync"
    "syscall"
    "time"
)

// OTPService defines an OTP API endpoint
type OTPService struct {
    Name        string
    URL         string
    ContentType string
    BuildBody   func(phone string) (string, error)
}

func main() {
    // Command-line flags
    phone := flag.String("phone", "", "Phone number in format 09xxxxxxxxx")
    count := flag.Int("n", 1, "Number of OTP requests per service")
    rate := flag.Int("rate", 1, "Requests per second per service")
    flag.Parse()

    // Validate phone
    if err := validatePhone(*phone); err != nil {
        log.Fatalf("invalid phone: %v", err)
    }

    // Setup services
    services := []OTPService{
        {
            Name:        "SnappFood",
            URL:         "https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass",
            ContentType: "application/x-www-form-urlencoded",
            BuildBody: func(phone string) (string, error) {
                v := url.Values{}
                v.Set("cellphone", phone)
                // add static query params
                v.Set("lat", "35.774")
                v.Set("long", "51.418")
                v.Set("client", "WEBSITE")
                return v.Encode(), nil
            },
        },
        {
            Name:        "Mobinnet",
            URL:         "https://my.mobinnet.ir/api/account/SendRegisterVerificationCode",
            ContentType: "application/json",
            BuildBody: func(phone string) (string, error) {
                return fmt.Sprintf(`{"cellNumber":"%s"}`, phone), nil
            },
        },
    }

    // Shared HTTP client
    client := &http.Client{Timeout: 10 * time.Second}

    // Context for cancellation
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    var wg sync.WaitGroup
    // Rate limiter tickers
    tickers := make([]*time.Ticker, len(services))
    for i := range services {
        tickers[i] = time.NewTicker(time.Second / time.Duration(*rate))
    }

    for i, svc := range services {
        for j := 0; j < *count; j++ {
            wg.Add(1)
            <-tickers[i].C
            go func(s OTPService, idx, iter int) {
                defer wg.Done()

                body, err := s.BuildBody(*phone)
                if err != nil {
                    log.Printf("[%s] build body error: %v", s.Name, err)
                    return
                }

                req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.URL, strings.NewReader(body))
                if err != nil {
                    log.Printf("[%s] new request error: %v", s.Name, err)
                    return
                }
                req.Header.Set("Content-Type", s.ContentType)

                resp, err := client.Do(req)
                if err != nil {
                    log.Printf("[%s] do request error: %v", s.Name, err)
                    return
                }
                defer resp.Body.Close()

                if resp.StatusCode >= 200 && resp.StatusCode < 300 {
                    log.Printf("[%s][%d] success status %d", s.Name, iter+1, resp.StatusCode)
                } else {
                    log.Printf("[%s][%d] unexpected status %d", s.Name, iter+1, resp.StatusCode)
                }
            }(svc, i, j)
        }
    }

    // Wait for all requests or cancellation
    wg.Wait()
    log.Println("All done.")
}

// validatePhone checks Iranian mobile format
func validatePhone(p string) error {
    re := regexp.MustCompile(`^09\d{9}$`)
    if !re.MatchString(p) {
        return errors.New("must match ^09xxxxxxxxx$")
    }
    return nil
}
