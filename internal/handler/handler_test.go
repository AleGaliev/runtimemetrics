package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AleGaliev/runtimemetrics/internal/logger"
	"github.com/AleGaliev/runtimemetrics/internal/observer"
	"github.com/AleGaliev/runtimemetrics/internal/service/crypto"
	"github.com/AleGaliev/runtimemetrics/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyHandler_GetPing(t *testing.T) {
	tests := []struct {
		connectError   error
		name           string
		expectedStatus int
	}{
		{
			name:           "successful connection",
			connectError:   nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "connection failed",
			connectError:   fmt.Errorf("connection error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := mocks.NewMockstorage(ctrl)
			mockConnector := mocks.NewMockconnector(ctrl)

			// Настраиваем ожидания
			mockConnector.EXPECT().Connect().Return(tt.connectError)
			logServer, _ := logger.CreateLogger()

			handler := CreateMyHandler(mockStorage, mockConnector, logServer, "", observer.NewEvent(), &crypto.Crypto{}, "")

			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestMyHandler_ServeHTTPUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockstorage(ctrl)
	mockConnector := mocks.NewMockconnector(ctrl)

	tests := []struct {
		updateError    error
		name           string
		body           string
		contentType    string
		method         string
		expectedStatus int
	}{
		{
			name:           "successful update",
			body:           `{"id": "test123","type": "gauge","value": 1.5}`,
			updateError:    nil,
			expectedStatus: http.StatusOK,
			contentType:    "application/json",
			method:         http.MethodPost,
		},
		{
			name:           "update error",
			body:           `{"id":"test","type":"gauge","value":1.5}`,
			updateError:    fmt.Errorf("update error"),
			expectedStatus: http.StatusBadRequest,
			contentType:    "application/json",
			method:         http.MethodPost,
		},
		{
			name:           "wrong content type",
			body:           `{"id":"test","type":"gauge","value":1.5}`,
			updateError:    nil,
			expectedStatus: http.StatusBadRequest,
			contentType:    "text/plain",
			method:         http.MethodPost,
		},
		{
			name:           "wrong method",
			body:           `{"id":"test","type":"gauge","value":1.5}`,
			updateError:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
			contentType:    "application/json",
			method:         http.MethodGet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем ожидание только для успешных случаев и случаев с ошибкой update
			if tt.expectedStatus == http.StatusOK ||
				(tt.expectedStatus == http.StatusBadRequest && tt.updateError != nil) {
				mockStorage.EXPECT().UpdateMetrics(gomock.Any()).Return(tt.updateError)
			}

			logServer, _ := logger.CreateLogger()

			handler := CreateMyHandler(mockStorage, mockConnector, logServer, "", observer.NewEvent(), &crypto.Crypto{}, "")

			req := httptest.NewRequest(tt.method, "/update/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestMyHandler_ServeHTTPBatchUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockstorage(ctrl)
	mockConnector := mocks.NewMockconnector(ctrl)

	tests := []struct {
		batchError     error
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "successful batch update",
			body:           `[{"id":"test1","type":"gauge","value":1.5},{"id":"test2","type":"counter","delta":1}]`,
			batchError:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "batch update error",
			body:           `[{"id":"test1","type":"gauge","value":1.5}]`,
			batchError:     fmt.Errorf("batch error"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.EXPECT().BatchUpdateMetrics(gomock.Any()).Return(tt.batchError)
			logServer, _ := logger.CreateLogger()
			handler := CreateMyHandler(mockStorage, mockConnector, logServer, "", observer.NewEvent(), &crypto.Crypto{}, "")

			req := httptest.NewRequest(http.MethodPost, "/updates/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestMyHandler_ServeHTTPValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockstorage(ctrl)
	mockConnector := mocks.NewMockconnector(ctrl)

	tests := []struct {
		valueError     error
		name           string
		body           string
		metrics        []byte
		expectedStatus int
		found          bool
	}{
		{
			name:           "successful get value",
			body:           `{"id":"test","type":"gauge"}`,
			metrics:        []byte(`{"id":"test","type":"gauge","value":1.5}`),
			found:          true,
			valueError:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "value not found",
			body:           `{"id":"test","type":"gauge"}`,
			metrics:        nil,
			found:          false,
			valueError:     nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "value error",
			body:           `{"id":"test","type":"gauge"}`,
			metrics:        nil,
			found:          false,
			valueError:     fmt.Errorf("value error"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.EXPECT().ValueMetrics(gomock.Any()).Return(tt.metrics, tt.found, tt.valueError)
			logServer, _ := logger.CreateLogger()
			handler := CreateMyHandler(mockStorage, mockConnector, logServer, "", observer.NewEvent(), &crypto.Crypto{}, "")

			req := httptest.NewRequest(http.MethodPost, "/value/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.JSONEq(t, string(tt.metrics), w.Body.String())
			}
		})
	}
}

func TestMyHandler_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockstorage(ctrl)
	mockConnector := mocks.NewMockconnector(ctrl)

	tests := []struct {
		addMetricError error
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "successful add metric",
			url:            "/update/gauge/test/1.5",
			addMetricError: nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "add metric error",
			url:            "/update/gauge/test/1.5",
			addMetricError: fmt.Errorf("add error"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid url",
			url:            "/update/gauge/test",
			addMetricError: nil,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.Count(tt.url, "/") >= 4 && tt.addMetricError != nil {
				mockStorage.EXPECT().AddMetric("gauge", "test", "1.5").Return(tt.addMetricError)
			} else if strings.Count(tt.url, "/") >= 4 {
				mockStorage.EXPECT().AddMetric("gauge", "test", "1.5").Return(nil)
			}
			logServer, _ := logger.CreateLogger()
			handler := CreateMyHandler(mockStorage, mockConnector, logServer, "", observer.NewEvent(), &crypto.Crypto{}, "")

			req := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestMyHandler_GetValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockstorage(ctrl)
	mockConnector := mocks.NewMockconnector(ctrl)

	tests := []struct {
		name           string
		metricName     string
		metricValue    string
		found          bool
		expectedStatus int
	}{
		{
			name:           "metric found",
			metricName:     "test_metric",
			metricValue:    "42.5",
			found:          true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "metric not found",
			metricName:     "unknown_metric",
			metricValue:    "",
			found:          false,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.EXPECT().GetMetrics(tt.metricName).Return(tt.metricValue, tt.found)
			logServer, _ := logger.CreateLogger()
			handler := CreateMyHandler(mockStorage, mockConnector, logServer, "", observer.NewEvent(), &crypto.Crypto{}, "")

			req := httptest.NewRequest(http.MethodGet, "/value/gauge/"+tt.metricName, nil)
			w := httptest.NewRecorder()

			// Используем chi router context для параметров
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", "gauge")
			rctx.URLParams.Add("name", tt.metricName)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.found {
				assert.Equal(t, tt.metricValue, w.Body.String())
			}
		})
	}
}

func TestMyHandler_ListMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockstorage(ctrl)
	mockConnector := mocks.NewMockconnector(ctrl)

	tests := []struct {
		getAllError    error
		name           string
		allMetrics     string
		expectedStatus int
	}{
		{
			name:           "successful list",
			allMetrics:     "<li>test_metric: 42.5</li>",
			getAllError:    nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "get all error",
			allMetrics:     "",
			getAllError:    fmt.Errorf("get all error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.EXPECT().GetAllMetric().Return(tt.allMetrics, tt.getAllError)
			logServer, _ := logger.CreateLogger()
			handler := CreateMyHandler(mockStorage, mockConnector, logServer, "", observer.NewEvent(), &crypto.Crypto{}, "")

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.getAllError == nil {
				assert.Contains(t, w.Body.String(), tt.allMetrics)
				assert.Contains(t, w.Body.String(), "Metrics List")
			}
		})
	}
}

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	successResponse(w, "")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Запрос обработан", response["message"])
}
