package service

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nivanov045/silver-octo-train/cmd/server/storage"
	"github.com/nivanov045/silver-octo-train/internal/metrics"
)

func Test_service_ParseAndSave(t *testing.T) {
	type args struct {
		name       string
		valueInt   int64
		valueFloat float64
		mType      string
	}
	tests := []struct {
		name string
		data args
	}{
		{
			name: "gauge correct",
			data: args{
				name:       "testSetGet1",
				valueInt:   0,
				valueFloat: 1.23,
				mType:      "gauge",
			},
		},
		{
			name: "counter correct",
			data: args{
				name:       "testSetGet2",
				valueInt:   2345,
				valueFloat: 0,
				mType:      "counter",
			},
		},
	}
	ser := service{storage.New()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := metrics.MetricsInterface{
				ID:    tt.data.name,
				MType: tt.data.mType,
				Delta: &tt.data.valueInt,
				Value: &tt.data.valueFloat,
			}
			if v.MType == "counter" {
				v.Value = nil
			} else {
				v.Delta = nil
			}
			marshal, err := json.Marshal(v)
			assert.NoError(t, err)
			err = ser.ParseAndSave(string(marshal))
		})
	}
}

func Test_service_ParseAndGet(t *testing.T) {
	type args struct {
		name       string
		valueInt   int64
		valueFloat float64
		mType      string
	}
	tests := []struct {
		name string
		data args
	}{
		{
			name: "gauge correct",
			data: args{
				name:       "testSetGet1",
				valueInt:   0,
				valueFloat: 1.23,
				mType:      "gauge",
			},
		},
		{
			name: "counter correct",
			data: args{
				name:       "testSetGet2",
				valueInt:   2345,
				valueFloat: 0,
				mType:      "counter",
			},
		},
	}
	ser := service{storage.New()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := metrics.MetricsInterface{
				ID:    tt.data.name,
				MType: tt.data.mType,
				Delta: &tt.data.valueInt,
				Value: &tt.data.valueFloat,
			}
			if v.MType == "counter" {
				v.Value = nil
			} else {
				v.Delta = nil
			}
			marshal, err := json.Marshal(v)
			assert.NoError(t, err)
			err = ser.ParseAndSave(string(marshal))
			assert.NoError(t, err)
			marshalGet, err := json.Marshal(metrics.MetricsInterface{
				ID:    tt.data.name,
				MType: tt.data.mType,
				Delta: nil,
				Value: nil,
			})
			got, err := ser.ParseAndGet(string(marshalGet))
			assert.NoError(t, err)
			assert.Equal(t, got, string(marshal))
		})
	}
}

func Test_service_GetKnownMetrics(t *testing.T) {
	type args struct {
		name       string
		valueInt   int64
		valueFloat float64
		mType      string
	}
	tests := []struct {
		name string
		ser  *service
		want []string
		set  []args
	}{
		{
			name: "correct",
			ser: &service{
				storage: storage.New(),
			},
			want: []string{"TestMetricC", "TestMetricG"},
			set: []args{
				{
					name:       "TestMetricG",
					valueInt:   0,
					valueFloat: 2345.1234,
					mType:      "gauge",
				}, {
					name:       "TestMetricC",
					valueInt:   123,
					valueFloat: 0,
					mType:      "counter",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, val := range tt.set {
				marshal, err := json.Marshal(metrics.MetricsInterface{
					ID:    val.name,
					MType: val.mType,
					Delta: &val.valueInt,
					Value: &val.valueFloat,
				})
				assert.NoError(t, err)
				tt.ser.ParseAndSave(string(marshal))
			}
			if got := tt.ser.GetKnownMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetKnownMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
