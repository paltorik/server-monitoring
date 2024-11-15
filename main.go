package main

import (
	"fmt"
	"os"
	"os/signal"
	"server-monitoring/collector"
	"server-monitoring/loader"
	"server-monitoring/metrics"

	"github.com/gin-gonic/gin"

	"time"
)

func main() {

	metrics, _ := loader.LoadConfig("storage/metrics.json")

	buffer := collector.NewMetricBuffer()

	workerMertic := collector.NewMetricsCollector(metrics, buffer)

	fileWriter := collector.NewFileWriter(metrics, buffer)

	stopWriteWorker := make(chan struct{})
	stopCleanWorker := make(chan struct{})

	collector.StartMetricWriteWorker(fileWriter, 20*time.Second, stopWriteWorker)
	collector.StartMetricCleanWorker(fileWriter, 120*time.Second, stopCleanWorker)
	workerMertic.Start()
	go startAPIServer(metrics)

	fmt.Println("Нажмите Ctrl+C для выхода...")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	close(stopWriteWorker)
	close(stopCleanWorker)
	workerMertic.Stop()

}
func startAPIServer(metrics []metrics.Metric) {
	r := gin.Default()

	// Пример маршрута для получения всех метрик
	r.GET("/metrics", func(c *gin.Context) {
		// Здесь логика получения списка метрик
		// Например, можно вызвать workerMetric.GetMetrics() и вернуть их в ответе
		c.JSON(200, gin.H{
			"metrics": metrics, // Пример, можно использовать реальные метрики
		})
	})

	// Пример маршрута для получения статистики по метрике
	r.GET("/metrics/:type", func(c *gin.Context) {
		metricType := c.Param("type")
		metric, _ := findMetricByType(metrics, metricType)
		data, _ := loader.ReadMetricsFromFile(metric)

		c.JSON(200, gin.H{
			"metrics": metric,
			"data":    data,
		})
	})

	// Запуск сервера на порту 8080
	if err := r.Run(":8080"); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
func findMetricByType(metrics []metrics.Metric, metricType string) (metrics.Metric, bool) {
	for _, metric := range metrics {
		if metric.GetType() == metricType {
			return metric, true
		}
	}
	return nil, false // Если метрика не найдена, возвращаем false
}
