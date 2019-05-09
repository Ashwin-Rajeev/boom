package main

import "time"

// APIStatus shows the current api status.
type APIStatus struct {
	TotalDuration     time.Duration
	MinRequestTime    time.Duration
	MaxRequestTime    time.Duration
	NumberOfRequests  int
	TotalResponseSize int64
	ErrorCount        int
}
