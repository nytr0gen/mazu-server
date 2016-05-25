package main

import (
	"crypto/tls"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const ACTION_TIMEOUT = time.Second * 90

func main() {
	rand.Seed(time.Now().UnixNano())

	client := &http.Client{
		Timeout: ACTION_TIMEOUT,
		Transport: &http.Transport{
			ResponseHeaderTimeout: ACTION_TIMEOUT,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestURL := "http://" + r.Host + r.URL.String()

		if rand.Intn(100) > 80 && (strings.HasSuffix(r.URL.Path, ".jpg") ||
			strings.HasSuffix(r.URL.Path, ".jpeg") ||
			strings.HasSuffix(r.URL.Path, ".gif") ||
			strings.HasSuffix(r.URL.Path, ".png")) {

			http.ServeFile(w, r, "./maz.jpg")
		} else if r.Host == "localhost" || strings.HasPrefix(r.Host, "localhost:") {
			w.Write([]byte("localhost"))
		} else {
			request, err := http.NewRequest(r.Method, requestURL, r.Body)
			if err != nil {
				log.Println(err)
			}
			request.Header = r.Header
			response, err := client.Do(request)
			if err != nil {
				log.Println(err)
			}

			for k, vals := range response.Header {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}

			_, err = io.Copy(w, response.Body)
			if err != nil {
				log.Println(err)
			}
		}
	})

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}
