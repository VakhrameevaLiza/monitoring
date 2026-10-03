package handlers

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	metrics "monitoring/internal/model"
)

func Update(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Неверный метод запроса", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")

	if len(parts) != 4 || parts[0] != "update" {
		http.Error(res, "Неверный путь", http.StatusNotFound)
		return
	}

	metricType := parts[1]
	metricName := parts[2]

	if !isValidMetricType(metricType) {
		http.Error(res, "Неверный тип метрики", http.StatusBadRequest)
		return
	}
	if !isValidMetricName(metricName) {
		http.Error(res, "Пустое имя метрики невалидно", http.StatusBadRequest)
		return
	}

	metric := metrics.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	err := parseGaugeValue(parts[3], &metric)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(metric)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "application/json")
	res.Write(json)
}

func isValidMetricType(metricType string) bool {
	return (metricType == metrics.Gauge || metricType == metrics.Counter)
}

func isValidMetricName(metricName string) bool {
	return len(metricName) > 0
}

func parseGaugeValue(rawMetricValue string, metric *metrics.Metrics) error {
	if metric.MType == metrics.Counter {
		counter, err := strconv.ParseInt(rawMetricValue, 10, 64)
		if err != nil {
			return errors.New("Неверное знанение counter метрики: " + rawMetricValue)
		}
		metric.Delta = &counter
		return nil
	}

	if metric.MType == metrics.Gauge {
		counter, err := strconv.ParseFloat(rawMetricValue, 64)
		if err != nil {
			return errors.New("Неверное знанение gauge метрики: " + rawMetricValue)
		}
		if math.IsNaN(counter) || math.IsInf(counter, 0) {
			return errors.New("значение gauge должно быть конечным числом")
		}
		metric.Value = &counter
		return nil
	}

	return errors.New("Неизвестное значение типа метрики: " + metric.MType)
}
