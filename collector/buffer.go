package collector

import (
	"server-monitoring/model"
	"sync"
)

// Структура для данных по метрике
type MetricData struct {
	mu     sync.Mutex
	values []model.BaseMetricValue
}

// Буфер для хранения данных метрик с мьютексами по типам метрик
type MetricBuffer struct {
	metricsData map[string]*MetricData
}

func NewMetricBuffer() *MetricBuffer {
	return &MetricBuffer{
		metricsData: make(map[string]*MetricData), // Инициализация карты

	}
}

// Добавление данных в буфер для конкретного типа метрики
func (mb *MetricBuffer) AddData(metricType string, data model.BaseMetricValue) {
	mb.metricsDataLock(metricType)
	defer mb.metricsDataUnlock(metricType)

	// Получаем или создаем данные для метрики
	metricData, _ := mb.metricsData[metricType]

	// Добавляем новые данные в массив
	metricData.values = append(metricData.values, data)
}

// Получение данных из буфера для конкретного типа метрики
func (mb *MetricBuffer) GetData(metricType string) []model.BaseMetricValue {

	data, exists := mb.metricsData[metricType]
	if !exists {
		return nil
	}
	return data.values
}

// Очистка данных в буфере для конкретного типа метрики
func (mb *MetricBuffer) ClearData(metricType string) {
	mb.metricsDataLock(metricType)
	defer mb.metricsDataUnlock(metricType)
	mb.metricsData[metricType].values = []model.BaseMetricValue{}

}

// Блокировка мьютекса для конкретного типа метрики
func (mb *MetricBuffer) metricsDataLock(metricType string) {
	metricData, exists := mb.metricsData[metricType]
	if !exists {
		// Создаем структуру для мьютекса, если она не существует
		metricData = &MetricData{}
		mb.metricsData[metricType] = metricData
	}
	metricData.mu.Lock()
}

// Разблокировка мьютекса для конкретного типа метрики
func (mb *MetricBuffer) metricsDataUnlock(metricType string) {
	metricData := mb.metricsData[metricType]
	metricData.mu.Unlock()
}
