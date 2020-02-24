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

func (conf *APIConfig) request() {
	status := &APIStatus{
		TotalDuration:  time.Minute,
		MinRequestTime: time.Minute,
		StatusCodes:    &StatusCodes{},
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
			status,
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

func run(httpClient *http.Client, method, url, requestBody string, header map[string]string, s *APIStatus) (requestDuration time.Duration, responseSize int) {
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
		s.StatusCodes.TwoXX++
	} else if resp.StatusCode == http.StatusContinue || resp.StatusCode == http.StatusSwitchingProtocols ||
		resp.StatusCode == http.StatusProcessing || resp.StatusCode == http.StatusEarlyHints {
		s.StatusCodes.OneXX++
	} else if resp.StatusCode == http.StatusMultipleChoices || resp.StatusCode == http.StatusMovedPermanently ||
		resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther ||
		resp.StatusCode == http.StatusNotModified {
		s.StatusCodes.ThreeXX++
	} else if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusPaymentRequired || resp.StatusCode == http.StatusForbidden ||
		resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed ||
		resp.StatusCode == http.StatusNotAcceptable || resp.StatusCode == http.StatusProxyAuthRequired ||
		resp.StatusCode == http.StatusRequestTimeout || resp.StatusCode == http.StatusContinue {
		s.StatusCodes.FourXX++
	} else if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusNotImplemented ||
		resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout || resp.StatusCode == http.StatusHTTPVersionNotSupported {
		s.StatusCodes.FiveXX++
	} else {
		s.StatusCodes.Others++
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
	return client, nil
}
