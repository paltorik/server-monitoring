package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"server-monitoring/metrics"
	"server-monitoring/model"
)

// Регистрируем типы метрик
var metricFactories = make(map[string]func() metrics.Metric)

// Регистрация всех метрик
func init() {
	metricFactories["cpu"] = func() metrics.Metric { return &metrics.CPUMetric{} }
	metricFactories["memory"] = func() metrics.Metric { return &metrics.MemoryMetric{} }
	metricFactories["disk"] = func() metrics.Metric { return &metrics.DickMertic{} }
	metricFactories["docker"] = func() metrics.Metric { return &metrics.DockerMetric{} }
}

// Фабрика для создания метрики по типу
func CreateMetric(baseMetric model.BaseMetric) (metrics.Metric, error) {
	// Ищем зарегистрированный тип метрики
	if factory, found := metricFactories[baseMetric.MetricType]; found {
		// Создаём экземпляр метрики через фабрику
		metric := factory()
		metric.Configure(baseMetric)
		return metric, nil
	}

	// Если метрика не найдена, возвращаем nil и ошибку
	return nil, fmt.Errorf("неизвестный тип метрики: %s", baseMetric.MetricType)
}

func LoadConfig(filePath string) ([]metrics.Metric, error) {
	// Определяем путь к файлу в директории storage
	//filePath := filepath.Join("storage", filename)

	// Читаем содержимое файла целиком
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %v", err)
	}

	// Десериализуем JSON в массив BaseMetric
	var baseMetrics []model.BaseMetric
	if err := json.Unmarshal(data, &baseMetrics); err != nil {
		return nil, fmt.Errorf("ошибка десериализации JSON: %v", err)
	}

	// Создаём метрики на основе десериализованных конфигураций
	var createdMetrics []metrics.Metric
	for _, baseMetric := range baseMetrics {
		metric, err := CreateMetric(baseMetric)
		if err != nil {
			fmt.Printf("Ошибка создания метрики %s: %v\n", baseMetric.Name, err)
			continue
		}
		createdMetrics = append(createdMetrics, metric)
	}

	return createdMetrics, nil
}
func ReadMetricsFromFile(metric metrics.Metric) ([]model.BaseMetricValue, error) {
	// Открываем файл
	file, err := os.Open(metric.GetFilePath())
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла: %v", err)
	}
	defer file.Close()

	// Массив для хранения данных из файла
	var stats []model.BaseMetricValue

	// Декодируем JSON из файла в массив
	decoder := json.NewDecoder(file)
	for {
		var stat model.BaseMetricValue
		if err := decoder.Decode(&stat); err != nil {
			// Если достигнут конец файла (EOF), завершаем чтение
			if err.Error() == "EOF" {
				break
			}
			// Иначе возвращаем ошибку
			return nil, fmt.Errorf("ошибка при декодировании JSON: %v", err)
		}
		// Добавляем декодированную запись в массив
		stats = append(stats, stat)
	}

	return stats, nil
}
