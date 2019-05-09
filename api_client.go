package main

import (
	"errors"
	"net/http"
	"time"
)

func (*APIConfig) request() {

}

func run() {

}

func newHTTPClient(timeOut int) (*http.Client, error) {
	to := time.Second * time.Duration(timeOut)
	client := &http.Client{
		Timeout: to,
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("url redirection not allowed")
	}
	return client, nil
}
