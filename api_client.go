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
		TotalDuration:  time.Millisecond,
		MinRequestTime: time.Hour,
		MaxRequestTime: time.Millisecond,
		StatusCodes:    &StatusCodes{},
	}
	client, err := newHTTPClient(conf.timeOut)
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	finalConf, err := prepareRequest(conf)
	if err != nil {
		log.Fatal(err)
	}
	for time.Since(start).Seconds() <= float64(finalConf.duration) && atomic.LoadInt32(&finalConf.interrupt) == 0 {
		for _, val := range finalConf.params {
			reqDuration, respSize := run(client, val.request, status)
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
	}
	finalConf.finalStatus <- status
}

func run(httpClient *http.Client, req *http.Request, s *APIStatus) (requestDuration time.Duration, responseSize int) {
	requestDuration = -1
	responseSize = -1
	start := time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("[Info] An error occurred while creating the request", err)
	}
	if resp == nil {
		fmt.Println("[Info] empty response")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[Info] An error occurred while reading  response body", err)
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		requestDuration = time.Since(start)
		responseSize = len(body) + int(headerSize(resp.Header))
		s.StatusCodes.TwoXX++
	} else if resp.StatusCode == http.StatusContinue || resp.StatusCode == http.StatusSwitchingProtocols ||
		resp.StatusCode == http.StatusProcessing {
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

func prepareRequest(conf *APIConfig) (*APIConfig, error) {
	for _, val := range conf.params {
		var buffer io.Reader
		if len(val.jsonBody) > 0 {
			buffer = bytes.NewBufferString(val.jsonBody)
		}
		req, err := http.NewRequest(val.Method, val.URL, buffer)
		if err != nil {
			fmt.Println("[Info] An error occurred while creating a new http request", err)
			return nil, err
		}
		for headerKey, headerValue := range val.Header {
			req.Header.Add(headerKey, headerValue)
		}
		val.request = req
	}
	return conf, nil
}
