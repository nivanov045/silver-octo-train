package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nivanov045/silver-octo-train/cmd/server/service"
	"github.com/nivanov045/silver-octo-train/cmd/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_api_updateMetricsHandler(t *testing.T) {
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
				r: "/update/counter/testCounter/100",
				m: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "correct gauge request",
			args: args{
				r: "/update/gauge/testGauge/100",
				m: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "1-part request",
			args: args{
				r: "/update/counter/",
				m: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "request invalid type",
			args: args{
				r: "/update/unknown/testCounter/100",
				m: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusNotImplemented,
			},
		},
		{
			name: "request invalid value",
			args: args{
				r: "/update/counter/testCounter/none",
				m: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "request invalid value",
			args: args{
				r: "/update/gauge/testGauge/none",
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
			h := http.HandlerFunc(a.updateMetricsHandler)
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func Test_api_getMetricsHandler(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "correct",
		},
	}
	storage := storage.New()
	serv := service.New(storage)
	a := api{serv}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestSend := httptest.NewRequest(http.MethodPost, "/update/counter/TestMetrics/100", nil)
			wSend := httptest.NewRecorder()
			hSend := http.HandlerFunc(a.updateMetricsHandler)
			hSend.ServeHTTP(wSend, requestSend)
			request := httptest.NewRequest(http.MethodGet, "/value/counter/TestMetrics/", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(a.getMetricsHandler)
			h.ServeHTTP(w, request)
			result := w.Result()

			respBody, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			defer result.Body.Close()

			assert.Equal(t, http.StatusOK, result.StatusCode)
			assert.Equal(t, "100", strings.Trim(string(respBody), "\n"))
		})
	}
}

func Test_api_rootHandler(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "correct",
		},
	}
	storage := storage.New()
	serv := service.New(storage)
	a := api{serv}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestSend := httptest.NewRequest(http.MethodPost, "/update/counter/TestMetrics/100", nil)
			wSend := httptest.NewRecorder()
			hSend := http.HandlerFunc(a.updateMetricsHandler)
			hSend.ServeHTTP(wSend, requestSend)
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(a.rootHandler)
			h.ServeHTTP(w, request)
			result := w.Result()

			respBody, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			defer result.Body.Close()

			assert.Equal(t, http.StatusOK, result.StatusCode)
			assert.Equal(t, "TestMetrics", strings.Trim(string(respBody), "\n"))
		})
	}
}
