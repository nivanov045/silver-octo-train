package service

import (
	"reflect"
	"testing"

	"github.com/nivanov045/silver-octo-train/cmd/server/storage"
)

func Test_service_ParseAndSave(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "correct gauge",
			args:    args{"gauge/tmp/1123.456"},
			wantErr: false,
		},
		{
			name:    "correct counter",
			args:    args{"counter/tmp/1123"},
			wantErr: false,
		},
		{
			name:    "0 parts",
			args:    args{""},
			wantErr: true,
		},
		{
			name:    "1 part",
			args:    args{"gauge"},
			wantErr: true,
		},
		{
			name:    "2 parts",
			args:    args{"gauge/testGauge"},
			wantErr: true,
		},
		{
			name:    "4 parts",
			args:    args{"gauge/testGauge/1123.456/tmp"},
			wantErr: true,
		},
		{
			name:    "unknown type",
			args:    args{"strangetype/tmp/1123"},
			wantErr: true,
		},
		{
			name:    "wrong gauge value",
			args:    args{"gauge/tmp/qewrt"},
			wantErr: true,
		},
		{
			name:    "wrong counter value",
			args:    args{"counter/tmp/qewrt"},
			wantErr: true,
		},
		{
			name:    "autotest 1",
			args:    args{"counter/testCounter/100"},
			wantErr: false,
		},
		{
			name:    "autotest 1",
			args:    args{"counter/testCounter/100"},
			wantErr: false,
		},
		{
			name:    "autotest 2",
			args:    args{"counter"},
			wantErr: true,
		},
		{
			name:    "autotest 3",
			args:    args{"gauge/testGauge/100"},
			wantErr: false,
		},
	}
	ser := service{storage.New()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ser.ParseAndSave(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("service.ParseAndSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_ParseAndGet(t *testing.T) {
	tests := []struct {
		name    string
		set     string
		args    string
		want    string
		wantErr bool
	}{
		{
			name:    "gauge correct",
			set:     "gauge/testSetGet1/1.23",
			args:    "gauge/testSetGet1",
			want:    "1.23",
			wantErr: false,
		},
		{
			name:    "counter correct",
			set:     "counter/testSetGet2/2345",
			args:    "counter/testSetGet2",
			want:    "2345",
			wantErr: false,
		},
		{
			name:    "unexisted",
			args:    "counter/testSetGet3",
			want:    "",
			wantErr: true,
		},
	}
	ser := service{storage.New()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ser.ParseAndSave(tt.set)
			got, err := ser.ParseAndGet(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.ParseAndGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("service.ParseAndGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetKnownMetrics(t *testing.T) {
	tests := []struct {
		name string
		ser  *service
		want []string
		set  []string
	}{
		{
			name: "correct",
			ser: &service{
				storage: storage.New(),
			},
			want: []string{"TestMetricC", "TestMetricG"},
			set:  []string{"gauge/TestMetricG/2345.1234", "counter/TestMetricC/123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, val := range tt.set {
				tt.ser.ParseAndSave(val)
			}
			if got := tt.ser.GetKnownMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetKnownMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
