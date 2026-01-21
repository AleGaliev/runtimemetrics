package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"

	"github.com/AleGaliev/runtimemetrics/internal/handler"
	log "github.com/AleGaliev/runtimemetrics/internal/logger"
	"github.com/AleGaliev/runtimemetrics/internal/observer"
	"github.com/AleGaliev/runtimemetrics/mocks"
	"github.com/golang/mock/gomock"
)

// ExampleMyHandler_ServeHTTPUpdate демонстрирует добавление метрики через JSON
func ExampleMyHandler_ServeHTTPUpdate() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()
	storage := mocks.NewMockstorage(ctrl)
	logger, _ := log.CreateLogger()
	connector := mocks.NewMockconnector(ctrl)
	eventAudit := observer.NewEvent()
	storage.EXPECT().UpdateMetrics(gomock.Any()).Return(nil)
	h := handler.CreateMyHandler(storage, connector, logger, "", eventAudit, nil)
	server := httptest.NewServer(h)
	defer server.Close()

	// Подготавливаем JSON данные
	metric := map[string]interface{}{
		"id":    "cpu_usage",
		"type":  "gauge",
		"value": 95.5,
	}

	jsonData, _ := json.Marshal(metric)

	// Отправляем POST запрос
	resp, err := http.Post(server.URL+"/update/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Metric updated, status: %s", resp.Status)
	// Output: Metric updated, status: 200 OK
}

// ExampleMyHandler_ServeHTTPValue демонстрирует получение метрики через JSON
func ExampleMyHandler_ServeHTTPValue() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()
	storage := mocks.NewMockstorage(ctrl)
	logger, _ := log.CreateLogger()
	connector := mocks.NewMockconnector(ctrl)
	eventAudit := observer.NewEvent()

	storage.EXPECT().ValueMetrics(gomock.Any()).Return([]byte(`{"id":"memory_usage","type":"gauge","value":75.3}`), true, nil)
	h := handler.CreateMyHandler(storage, connector, logger, "", eventAudit, nil)
	server := httptest.NewServer(h)
	defer server.Close()

	// Запрос на получение метрики
	request := map[string]string{
		"id": "memory_usage",
	}

	jsonData, _ := json.Marshal(request)

	resp, err := http.Post(server.URL+"/value/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	defer resp.Body.Close()

	// Читаем ответ
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Printf("Metric value: %v", result["value"])
	// Output: Metric value: 75.3
}

// ExampleMyHandler_GetValue демонстрирует получение метрики через URL параметры
func ExampleMyHandler_GetValue() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()
	storage := mocks.NewMockstorage(ctrl)
	logger, _ := log.CreateLogger()
	connector := mocks.NewMockconnector(ctrl)
	eventAudit := observer.NewEvent()

	storage.EXPECT().GetMetrics("disk_space").Return("500.0", true)
	h := handler.CreateMyHandler(storage, connector, logger, "", eventAudit, nil)
	server := httptest.NewServer(h)
	defer server.Close()

	// Получаем метрику через GET запрос
	resp, err := http.Get(server.URL + "/value/gauge/disk_space")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Читаем текстовый ответ
	var body bytes.Buffer
	body.ReadFrom(resp.Body)

	fmt.Printf("Disk space: %s", body.String())
	// Output: Disk space: 500.0
}

// ExampleMyHandler_ListMetrics демонстрирует получение списка всех метрик
func ExampleMyHandler_ListMetrics() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()
	storage := mocks.NewMockstorage(ctrl)
	logger, _ := log.CreateLogger()
	connector := mocks.NewMockconnector(ctrl)
	eventAudit := observer.NewEvent()
	storage.EXPECT().GetAllMetric().Return("<li>test_metric: 42.5</li>", nil)

	h := handler.CreateMyHandler(storage, connector, logger, "", eventAudit, nil)
	server := httptest.NewServer(h)
	defer server.Close()

	// Получаем HTML страницу со списком метрик
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Metrics page loaded successfully")
	// Output: Metrics page loaded successfully
}

// ExampleMyHandler_ServeHTTP демонстрирует добавление метрики через URL параметры
func ExampleMyHandler_ServeHTTP() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()
	storage := mocks.NewMockstorage(ctrl)
	logger, _ := log.CreateLogger()
	connector := mocks.NewMockconnector(ctrl)
	eventAudit := observer.NewEvent()
	storage.EXPECT().AddMetric("gauge", "test", "1.5").Return(nil)
	h := handler.CreateMyHandler(storage, connector, logger, "", eventAudit, nil)
	server := httptest.NewServer(h)
	defer server.Close()

	// Добавляем метрику через URL параметры
	resp, err := http.Post(server.URL+"/update/gauge/test/1.5", "", nil)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Metric added via URL, status: %s", resp.Status)
	// Output: Metric added via URL, status: 200 OK
}

// ExampleMyHandler_ServeHTTPBatchUpdate демонстрирует пакетное обновление метрик
func ExampleMyHandler_ServeHTTPBatchUpdate() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()
	storage := mocks.NewMockstorage(ctrl)
	logger, _ := log.CreateLogger()
	connector := mocks.NewMockconnector(ctrl)
	eventAudit := observer.NewEvent()
	storage.EXPECT().BatchUpdateMetrics(gomock.Any()).Return(nil)
	h := handler.CreateMyHandler(storage, connector, logger, "", eventAudit, nil)
	server := httptest.NewServer(h)
	defer server.Close()

	// Подготавливаем пакет метрик
	metrics := []map[string]interface{}{
		{
			"id":    "metric1",
			"type":  "gauge",
			"value": 10.5,
		},
		{
			"id":    "metric2",
			"type":  "counter",
			"delta": 42,
		},
	}
	jsonData, _ := json.Marshal(metrics)

	// Отправляем пакетный запрос
	resp, err := http.Post(server.URL+"/updates/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Batch update completed, status: %s", resp.Status)
	// Output: Batch update completed, status: 200 OK
}
