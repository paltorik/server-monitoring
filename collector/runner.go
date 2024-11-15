package collector

import (
	"fmt"
	"server-monitoring/metrics"
	"time"
)

type MetricsCollector struct {
	metrics []metrics.Metric
	buffer  *MetricBuffer
	stopCh  chan struct{} // Канал для остановки процесса сбора
}

func NewMetricsCollector(metrics []metrics.Metric, newBuffer *MetricBuffer) *MetricsCollector {
	return &MetricsCollector{
		metrics: metrics,
		buffer:  newBuffer,
		stopCh:  make(chan struct{}),
	}
}

func (mc *MetricsCollector) Start() {
	for _, metric := range mc.metrics {
		if metric.IsActive() {
			go mc.collectMetric(metric) // Запускаем сбор данных для каждой активной метрики в отдельной горутине
		}
	}
}

// collectMetric запускает сбор метрики с учётом её периодичности
func (mc *MetricsCollector) collectMetric(metric metrics.Metric) {
	ticker := time.NewTicker(metric.GetPeriod()) // Создаем тикер для заданного интервала
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Запускаем сбор метрики
			value, err := metric.Collect()
			if err != nil {
				fmt.Printf("Ошибка сбора метрики %s: %v\n", metric.GetName(), err)
				continue
			}

			mc.buffer.AddData(metric.GetType(), value)

		case <-mc.stopCh:
			// Останавливаем сбор данных, если получен сигнал остановки
			return
		}
	}
}

// Stop останавливает процесс сбора всех метрик
func (mc *MetricsCollector) Stop() {
	close(mc.stopCh) // Закрытие канала остановит все горутины
}
