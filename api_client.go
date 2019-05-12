package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// var (
// 	requestDuration time.Duration
// 	responseSize    int
// )

func (conf *APIConfig) request() {
	status := &APIStatus{
		TotalDuration:  time.Minute,
		MinRequestTime: time.Minute,
		MaxRequestTime: time.Minute,
	}
	client, err := newHTTPClient(conf.timeOut)
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	for time.Since(start).Seconds() <= float64(conf.duration) && atomic.LoadInt32(&conf.interrupt) == 0 {
		reqDuration, respSize := run(
			client,
			conf.method,
			conf.url,
			conf.body,
			conf.header,
		)
		if respSize > 0 {
			status.NumberOfRequests++
			status.TotalResponseSize += int64(respSize)
			status.TotalDuration += reqDuration
			status.MaxRequestTime = findMaxRequestTime(reqDuration, status.MaxRequestTime)
			status.MinRequestTime = findMinRequestTime(reqDuration, status.MinRequestTime)
		} else {
			status.ErrorCount++
		}
	}
	conf.finalStatus <- status
}

func run(httpClient *http.Client, method, url, requestBody string, header map[string]string) (requestDuration time.Duration, responseSize int) {
	var buffer io.Reader
	requestDuration = -1
	responseSize = -1
	if len(requestBody) > 0 {
		buffer = bytes.NewBufferString(requestBody)
	}
	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		fmt.Println("[Info] An error occured while creating a new http request", err)
		return
	}

	for headerKey, headerValue := range header {
		req.Header.Add(headerKey, headerValue)
	}

	start := time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("[Info] An error occured while creating the request", err)
	}
	if resp == nil {
		fmt.Println("[Info] empty response")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[Info] An error occured while reading  response body", err)
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		requestDuration = time.Since(start)
		responseSize = len(body) + int(headerSize(resp.Header))
	} else {
		fmt.Println("[Info] Got status code", resp.StatusCode, "from", resp.Header, "content", string(body), req)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	return
}

func newHTTPClient(timeOut int) (*http.Client, error) {
	client := &http.Client{}
	client.Transport = &http.Transport{
		ResponseHeaderTimeout: time.Millisecond * time.Duration(timeOut),
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("url redirection not allowed")
	}
	client.CloseIdleConnections()
	return client, nil
}
