package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	var contracts []Contract
	contractBytes, err := ioutil.ReadFile("./contract.json")
	if err != nil {
		log.Fatal("failed to read")
	}
	err = json.Unmarshal(contractBytes, &contracts)
	if err != nil {
		log.Fatal("failed to unmarshal contracts", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print("Failed to read request body")
		}

		fmt.Println("----------------Request----------------")
		fmt.Println(fmt.Sprintf("%s %s", r.Method, r.RequestURI))
		r.Header.Write(os.Stdout)
		fmt.Println(string(bodyBytes))
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

		for _, contract := range contracts {
			if contract.HttpOverrideForwardedRequest != nil {
				if !strings.EqualFold(contract.HttpRequest.Path, r.URL.Path) {
					continue
				}
				httpRequestContract := contract.HttpOverrideForwardedRequest.HttpRequest
				newUrl := fmt.Sprintf("%s://%s%s", httpRequestContract.Protocol, httpRequestContract.Host, r.RequestURI)

				fmt.Println("----------------Forwarding--------------")
				fmt.Println(newUrl)
				fmt.Println(contract.HttpRequest.Body.JSON)

				request, _ := http.NewRequest(r.Method, newUrl, bytes.NewReader(bodyBytes))
				for key, values := range r.Header {
					for _, value := range values {
						request.Header.Add(key, value)
					}
				}
				for key, value := range httpRequestContract.Headers {
					request.Header.Set(key, value)
				}
				request.Header.Del("Accept-Encoding")

				client := http.Client{
					Timeout: time.Duration(5) * time.Second,
				}
				response, _ := client.Do(request)
				responseBodyBytes, _ := ioutil.ReadAll(response.Body)
				w.WriteHeader(response.StatusCode)
				w.Write(responseBodyBytes)
				for key, values := range response.Header {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}
				fmt.Println("----------------Response----------------")
				fmt.Println(response.StatusCode)
				fmt.Println(string(responseBodyBytes))
				fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

				return
			}
			if !strings.EqualFold(contract.HttpRequest.Path, r.URL.Path) {
				continue
			}
			if !strings.EqualFold(contract.HttpRequest.Method, r.Method) {
				continue
			}
			if !IsEqualJson([]byte(contract.HttpRequest.Body.JSON), bodyBytes) {
				continue
			}

			headerMatch := true
			for key, value := range contract.HttpRequest.Headers {
				var regexValue = regexp.MustCompile(value)
				if !regexValue.MatchString(r.Header.Get(key)) {
					headerMatch = false
					break
				}
			}
			if !headerMatch {
				continue
			}

			queryStringMatch := true
			for k, v := range contract.HttpRequest.QueryString {
				if !regexp.MustCompile(v).MatchString(r.URL.Query().Get(k)) {
					queryStringMatch = false
					break
				}
			}
			if !queryStringMatch {
				continue
			}

			fmt.Println("----------------Matching----------------")
			fmt.Println(fmt.Sprintf("%s %s", contract.HttpRequest.Method, contract.HttpRequest.Path))
			fmt.Println(contract.HttpRequest.Body.JSON)
			fmt.Println("----------------Response----------------")
			fmt.Println(contract.HttpResponse.StatusCode)
			fmt.Println(contract.HttpResponse.Body)
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

			w.WriteHeader(contract.HttpResponse.StatusCode)
			w.Write([]byte(contract.HttpResponse.Body))
			w.Header().Set("Content-Type", "application/json")
			return
		}

		fmt.Println("------------------NotFound------------------")
		w.WriteHeader(404)
	})
	http.ListenAndServe(":8080", nil)
}
