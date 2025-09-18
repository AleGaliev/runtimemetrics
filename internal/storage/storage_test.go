package storage

import (
	"testing"

	"github.com/AleGaliev/kubercontroller/internal/filestore"
	models "github.com/AleGaliev/kubercontroller/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestStorage_AddMetric(t *testing.T) {
	var (
		allocValue          = 123.45
		buckHashValue       = 67.89
		count         int64 = 3
	)

	type fields struct {
		Metrics map[string]models.Metrics
	}
	type args struct {
		myType string
		name   string
		value  string
	}
	var tests = []struct {
		name    string
		args    args
		parent  fields
		result  fields
		wantErr bool
	}{
		{
			name: "AddMetric guage",
			args: args{
				myType: models.Gauge,
				name:   "BuckHash",
				value:  "67.89",
			},
			wantErr: false,
			parent: fields{
				Metrics: map[string]models.Metrics{
					"Alloc": {
						ID:    "Alloc",
						MType: models.Gauge,
						Value: &allocValue,
					},
				},
			},
			result: fields{
				Metrics: map[string]models.Metrics{
					"Alloc": {
						ID:    "Alloc",
						MType: models.Gauge,
						Value: &allocValue,
					},
					"BuckHash": {
						ID:    "BuckHash",
						MType: models.Gauge,
						Value: &buckHashValue,
					},
				},
			},
		},
		{
			name: "AddMetric caunter",
			args: args{
				myType: models.Counter,
				name:   "tetsName",
				value:  "3",
			},
			wantErr: false,
			parent: fields{
				Metrics: map[string]models.Metrics{
					"Alloc": {
						ID:    "Alloc",
						MType: models.Gauge,
						Value: &allocValue,
					},
				},
			},
			result: fields{
				Metrics: map[string]models.Metrics{
					"Alloc": {
						ID:    "Alloc",
						MType: models.Gauge,
						Value: &allocValue,
					},
					"tetsName": {
						ID:    "tetsName",
						MType: models.Counter,
						Delta: &count,
					},
				},
			},
		},
		{
			name: "error check",
			args: args{
				myType: "testType",
				name:   "BuckHash",
				value:  "67.89",
			},
			wantErr: true,
			parent: fields{
				Metrics: map[string]models.Metrics{
					"Alloc": {
						ID:    "Alloc",
						MType: models.Gauge,
						Value: &allocValue,
					},
				},
			},
			result: fields{
				Metrics: map[string]models.Metrics{
					"Alloc": {
						ID:    "Alloc",
						MType: models.Gauge,
						Value: &allocValue,
					},
				},
			},
		},
	}
	f := filestore.NewFileStore("storage_test.json")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Metrics:       tt.parent.Metrics,
				StoreInterval: 5,
				FileStorage:   f,
			}
			if err := s.AddMetric(tt.args.myType, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("AddMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.result.Metrics, s.Metrics)
		})
	}
}

func TestStorage_GetAllMetric(t *testing.T) {
	var (
		allocValue       = 123.45
		count      int64 = 3
		metrics          = map[string]models.Metrics{
			"Alloc": {
				ID:    "Alloc",
				MType: models.Gauge,
				Value: &allocValue,
			},
			"tetsName": {
				ID:    "tetsName",
				MType: models.Counter,
				Delta: &count,
			},
		}
	)

	var tests = []struct {
		name string
		want string
	}{
		{
			name: "GetAllMetric",
			want: "<li> Alloc: 123.45</li><li> tetsName: 3</li>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Metrics: metrics,
			}
			if got, _ := s.GetAllMetric(); got != tt.want {
				t.Errorf("GetAllMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetMetrics(t *testing.T) {
	var (
		allocValue       = 123.45
		count      int64 = 3
		metrics          = map[string]models.Metrics{
			"Alloc": {
				ID:    "Alloc",
				MType: models.Gauge,
				Value: &allocValue,
			},
			"tetsName": {
				ID:    "tetsName",
				MType: models.Counter,
				Delta: &count,
			},
		}
	)

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
		ok   bool
	}{
		{
			name: "GetMetrics",
			args: args{
				name: "Alloc",
			},
			want: "123.45",
			ok:   true,
		},
		{
			name: "Negative",
			args: args{
				name: "Unometrics",
			},
			want: "",
			ok:   false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Metrics: metrics,
			}
			got, ok := s.GetMetrics(tt.args.name)
			if ok != tt.ok {
				t.Errorf("GetMetrics() got = %v, want %v", got, tt.want)
			}
			if got != tt.want {
				t.Errorf("GetMetrics() got = %v, want %v", got, tt.want)
			}
		})
	}
}
