package main

import "net/http"

// RequestParams for http request
type RequestParams struct {
	URL      string            `json:"url"`
	Method   string            `json:"method"`
	Header   map[string]string `json:"header"`
	Body     interface{}       `json:"body"`
	jsonBody string
	request  *http.Request
}

// APIConfig which stores the api configuration
type APIConfig struct {
	concurrentConnections int
	duration              int
	timeOut               int
	finalStatus           chan *APIStatus
	interrupt             int32
	params                []*RequestParams
}

func newAPIConfig(goroutines, duration, timeOut int, finalStatusChan chan *APIStatus, params []*RequestParams) *APIConfig {
	a := &APIConfig{
		concurrentConnections: goroutines,
		duration:              duration,
		timeOut:               timeOut,
		finalStatus:           finalStatusChan,
		params:                params,
	}
	return a
}
