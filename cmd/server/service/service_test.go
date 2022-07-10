package service

import (
	"testing"

	"github.com/nivanov045/silver-octo-train/cmd/server/storage"
)

func Test_service_ParseAndSet(t *testing.T) {
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
			if err := ser.ParseAndSet(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("service.ParseAndSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
