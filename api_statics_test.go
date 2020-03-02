package main

import (
	"net/http"
	"testing"
	"time"
)

func Test_headerSize(t *testing.T) {
	type args struct {
		headers http.Header
	}
	tests := []struct {
		name       string
		args       args
		wantResult int64
	}{
		{
			name: "with header",
			args: args{
				http.Header{
					"header": []string{"application/json"},
				},
			},
			wantResult: 28,
		},
		{
			name: "without header",
			args: args{
				http.Header{},
			},
			wantResult: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := headerSize(tt.args.headers); gotResult != tt.wantResult {
				t.Errorf("headerSize() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_findMaxRequestTime(t *testing.T) {
	tm := func(d string) time.Duration {
		dr, _ := time.ParseDuration(d)
		return dr
	}
	type args struct {
		t1 time.Duration
		t2 time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "test1",
			args: args{
				t1: tm("10s"),
				t2: tm("5s"),
			},
			want: tm("10s"),
		},
		{
			name: "test2",
			args: args{
				t1: tm("5s"),
				t2: tm("10s"),
			},
			want: tm("10s"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findMaxRequestTime(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("findMaxRequestTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findMinRequestTime(t *testing.T) {
	tm := func(d string) time.Duration {
		dr, _ := time.ParseDuration(d)
		return dr
	}
	type args struct {
		t1 time.Duration
		t2 time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "test1",
			args: args{
				t1: tm("10s"),
				t2: tm("5s"),
			},
			want: tm("5s"),
		},
		{
			name: "test2",
			args: args{
				t1: tm("5s"),
				t2: tm("10s"),
			},
			want: tm("5s"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findMinRequestTime(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("findMinRequestTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIConfig_stop(t *testing.T) {
	type fields struct {
		interrupt int32
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test1",
			fields: fields{
				interrupt: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &APIConfig{
				interrupt: tt.fields.interrupt,
			}
			conf.stop()
		})
	}
}
