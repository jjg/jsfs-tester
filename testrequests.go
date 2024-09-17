package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

func connectTest(proto string, name string, port int) (bool, int64, error) {
	startTime := time.Now()

	reqUrl := fmt.Sprintf("%s://%s:%d", proto, name, port)
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return false, time.Since(startTime).Milliseconds(), err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, time.Since(startTime).Milliseconds(), err
	}

	log.Printf("Got a response code (%d), server is up!", res.StatusCode)

	// if we got this far, success!
	return true, time.Since(startTime).Milliseconds(), nil
}

func postTest(proto string, name string, port int, fileName string) (bool, int64, error) {
	startTime := time.Now()

	// TODO: Don't hard-code the access key
	// TODO: Offer to randomize the body to bypass dedupe
	reqUrl := fmt.Sprintf("%s://%s:%d/%s?access_key=foo", proto, name, port, fileName)
	jsonBody := []byte(`{"message": "I am a teapot"}`)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, reqUrl, bodyReader)
	if err != nil {
		return false, time.Since(startTime).Milliseconds(), err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, time.Since(startTime).Milliseconds(), err
	}

	log.Printf("POST got a response code %d", res.StatusCode)

	// TODO: Might need to check for variations of 2xx
	if res.StatusCode == 200 {
		return true, time.Since(startTime).Milliseconds(), nil
	} else {
		return false, time.Since(startTime).Milliseconds(), nil
	}
}
