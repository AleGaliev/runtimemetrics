package filestore

import (
	"os"
	"reflect"
	"testing"
)

func TestWriteMetrics(t *testing.T) {
	type args struct {
		filename string
		data     []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "PositiveWriteMetrics",
			args: args{
				filename: "write_metrics.json",
				data:     []byte("Hello, World!"),
			},
			wantErr: false,
		},
		{
			name: "NegativeWriteMetrics",
			args: args{
				filename: "test/write_metrics.json",
				data:     []byte("Hello, World!"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFileStore(tt.args.filename)
			if err := f.WriteMetrics(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("WriteMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
			_ = os.Remove(tt.args.filename)
		})
	}
}

func TestReadMetrics(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "PositiveReadMetrics",
			args: args{
				filename: "read_metrics.json",
			},
			want:    []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name: "fileNotFounfReadMetrics",
			args: args{
				filename: "not_exists.json",
			},
			want:    nil,
			wantErr: true,
		},
	}
	f := NewFileStore("read_metrics.json")
	if err := f.WriteMetrics([]byte("Hello, World!")); err != nil {
		t.Fatalf("could not open metrics file: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFileStore(tt.args.filename)
			got, err := f.ReadMetrics()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadMetrics() got = %v, want %v", got, tt.want)
			}
			_ = os.Remove(tt.args.filename)
		})

	}
}
