package repository

import (
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"testing"
)

func TestConvertMemStatsInNameMetrics(t *testing.T) {
	type args struct {
		memStats runtime.MemStats
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertMemStatsInNameMetrics(tt.args.memStats); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertMemStatsInNameMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendMetrics_DoRequest(t *testing.T) {
	type fields struct {
		URL    map[string]url.URL
		Client *http.Client
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SendMetrics{
				URL:    tt.fields.URL,
				Client: tt.fields.Client,
			}
			s.DoRequest()
		})
	}
}

func TestSendMetrics_InitMetrics(t *testing.T) {
	type fields struct {
		URL    map[string]url.URL
		Client *http.Client
	}
	type args struct {
		mericsNameValue map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SendMetrics{
				URL:    tt.fields.URL,
				Client: tt.fields.Client,
			}
			s.InitMetrics(tt.args.mericsNameValue)
		})
	}
}

func Test_createRequest(t *testing.T) {
	type args struct {
		name    string
		metrics string
	}
	tests := []struct {
		name string
		args args
		want url.URL
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createRequest(tt.args.name, tt.args.metrics); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
