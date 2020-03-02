package main

import (
	"reflect"
	"testing"
)

func Test_newAPIConfig(t *testing.T) {
	type args struct {
		goroutines      int
		duration        int
		timeOut         int
		finalStatusChan chan *APIStatus
		params          []*RequestParams
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
				duration:        10,
				timeOut:         1000,
				finalStatusChan: nil,
				params: []*RequestParams{
					{
						URL:    "www.sample.com",
						Method: "POST",
						Header: nil,
					},
				},
			},
			want: &APIConfig{
				concurrentConnections: 10,
				duration:              10,
				timeOut:               1000,
				finalStatus:           nil,
				interrupt:             0,
				params: []*RequestParams{
					&RequestParams{
						URL:    "www.sample.com",
						Method: "POST",
						Header: nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newAPIConfig(tt.args.goroutines, tt.args.duration, tt.args.timeOut, tt.args.finalStatusChan, tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAPIConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
