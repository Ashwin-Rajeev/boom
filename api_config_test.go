package main

import (
	"reflect"
	"testing"
)

func Test_newAPIConfig(t *testing.T) {
	type args struct {
		goroutines      int
		url             string
		method          string
		header          map[string]string
		duration        int
		body            string
		timeOut         int
		finalStatusChan chan *APIStatus
	}
	tests := []struct {
		name string
		args args
		want *APIConfig
	}{
		{
			name: "test1",
			args: args{
				goroutines:      10,
				url:             "www.sample.com",
				method:          "POST",
				header:          nil,
				duration:        10,
				body:            "sample",
				timeOut:         1000,
				finalStatusChan: nil,
			},
			want: &APIConfig{
				concurrentConnections: 10,
				url:                   "www.sample.com",
				method:                "POST",
				header:                nil,
				duration:              10,
				body:                  "sample",
				timeOut:               1000,
				finalStatus:           nil,
				interrupt:             0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newAPIConfig(tt.args.goroutines, tt.args.url, tt.args.method, tt.args.header, tt.args.duration, tt.args.body, tt.args.timeOut, tt.args.finalStatusChan); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAPIConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
