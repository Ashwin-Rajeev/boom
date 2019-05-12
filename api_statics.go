package main

import (
	"net/http"
	"sync/atomic"
	"time"
)

// APIStatus shows the current api status.
type APIStatus struct {
	TotalDuration     time.Duration
	MinRequestTime    time.Duration
	MaxRequestTime    time.Duration
	NumberOfRequests  int
	TotalResponseSize int64
	ErrorCount        int
}

func headerSize(headers http.Header) (result int64) {
	result = 0
	for k, v := range headers {
		result += int64(len(k) + len(": \r\n"))
		for _, s := range v {
			result += int64(len(s))
		}
	}
	result += int64(len("\r\n"))
	return result
}

func findMaxRequestTime(t1, t2 time.Duration) time.Duration {
	if t1 > t2 {
		return t1
	}
	return t2
}

func findMinRequestTime(t1, t2 time.Duration) time.Duration {
	if t1 < t2 {
		return t1
	}
	return t2
}

func (conf *APIConfig) stop() {
	atomic.StoreInt32(&conf.interrupt, 1)
}
