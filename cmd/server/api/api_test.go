package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nivanov045/silver-octo-train/cmd/server/service"
	"github.com/nivanov045/silver-octo-train/cmd/server/storage"
	"github.com/stretchr/testify/assert"
)

func Test_api_requestMetricsHandler(t *testing.T) {
	type args struct {
		r string
		m string
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want
	}{
		{
			name: "correct counter request",
			args: args{
				r: "/value/counter/testCounter/100",
				m: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "correct gauge request",
			args: args{
				r: "/value/gauge/testGauge/100",
				m: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "1-part request",
			args: args{
				r: "/value/counter/",
				m: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "request with incorrect method",
			args: args{
				r: "/value/counter/",
				m: http.MethodPut,
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "request invalid type",
			args: args{
				r: "/value/unknown/testCounter/100",
				m: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusNotImplemented,
			},
		},
		{
			name: "request invalid value",
			args: args{
				r: "/value/counter/testCounter/none",
				m: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "request invalid value",
			args: args{
				r: "/value/gauge/testGauge/none",
				m: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	storage := storage.New()
	serv := service.New(storage)
	a := api{serv}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.m, tt.args.r, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(a.requestMetricsHandler)
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
