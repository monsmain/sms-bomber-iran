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
