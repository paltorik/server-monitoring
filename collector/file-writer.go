package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"server-monitoring/metrics"
	"server-monitoring/model"
	"sync"
	"time"
)

type FileMutex struct {
	mu sync.Mutex
}

type FileWriter struct {
	metricBuf *MetricBuffer
	metrics   []metrics.Metric
	mutexData map[string]*FileMutex
}

// Создание нового FileWriter
func NewFileWriter(metrics []metrics.Metric, metricBuf *MetricBuffer) *FileWriter {
	return &FileWriter{
		metrics:   metrics,
		metricBuf: metricBuf,
		mutexData: make(map[string]*FileMutex),
	}
}

// Метод для записи данных в файл для конкретной метрики
func (fw *FileWriter) WriteToFile() error {

	for _, metric := range fw.metrics {

		data := fw.metricBuf.GetData(metric.GetType())

		if data == nil || len(data) == 0 {

			continue // Если данных нет, переходим к следующей метрике
		}
		// Получаем путь к файлу для сохранения данных
		filePath := metric.GetFilePath()

		fw.mutextDataLock(metric.GetType())
		// Открываем или создаем файл для записи
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fw.mutextDataUnlock(metric.GetType())
			return fmt.Errorf("не удалось открыть файл %s для записи: %v", filePath, err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		for _, value := range data {
			if err := encoder.Encode(value); err != nil {
				file.Close()
				fw.mutextDataUnlock(metric.GetType())
				return fmt.Errorf("не удалось записать данные в файл %s: %v", filePath, err)
			}
		}

		// Закрываем файл
		file.Close()

		// Разблокируем доступ
		fw.mutextDataUnlock(metric.GetType())

		// Удаляем данные из буфера только после успешной записи
		fw.metricBuf.ClearData(metric.GetType())
	}
	return nil
}

func (fw *FileWriter) mutextDataLock(metricType string) {
	_, exists := fw.mutexData[metricType]
	if !exists {
		fw.mutexData[metricType] = &FileMutex{}
	}
	fw.mutexData[metricType].mu.Lock()
}

func (fw *FileWriter) mutextDataUnlock(metricType string) {
	mutexData := fw.mutexData[metricType]
	mutexData.mu.Unlock()
}

func (fw *FileWriter) CleanOldData() error {
	for _, metric := range fw.metrics {
		// Получаем путь к файлу и период хранения данных
		filePath := metric.GetFilePath()
		retentionPeriod := metric.GetRetentionPeriod()

		// Блокируем доступ к файлу
		fw.mutextDataLock(metric.GetType())

		// Открываем файл для чтения
		file, err := os.Open(filePath)
		if err != nil {
			fw.mutextDataUnlock(metric.GetType())
			return fmt.Errorf("не удалось открыть файл %s для очистки: %v", filePath, err)
		}

		// Считываем все записи из файла
		var allData []model.BaseMetricValue
		decoder := json.NewDecoder(file)
		for {
			var value model.BaseMetricValue
			if err := decoder.Decode(&value); err != nil {
				if err == io.EOF {
					break
				}
				file.Close()
				fw.mutextDataUnlock(metric.GetType())
				return fmt.Errorf("ошибка чтения из файла %s: %v", filePath, err)
			}
			allData = append(allData, value)
		}
		file.Close()

		// Вычисляем границу для удаления старых данных
		cutoffTime := time.Now().Add(-retentionPeriod)

		// Отбираем только актуальные данные
		var filteredData []model.BaseMetricValue
		for _, value := range allData {
			if value.ExecutedAt.After(cutoffTime) {
				filteredData = append(filteredData, value)
			}
		}

		// Перезаписываем файл только с актуальными данными
		file, err = os.Create(filePath)
		if err != nil {
			fw.mutextDataUnlock(metric.GetType())
			return fmt.Errorf("не удалось открыть файл %s для перезаписи: %v", filePath, err)
		}

		encoder := json.NewEncoder(file)
		for _, value := range filteredData {
			if err := encoder.Encode(value); err != nil {
				file.Close()
				fw.mutextDataUnlock(metric.GetType())
				return fmt.Errorf("не удалось записать данные в файл %s: %v", filePath, err)
			}
		}

		// Закрываем файл и снимаем блокировку
		file.Close()
		fw.mutextDataUnlock(metric.GetType())
	}
	return nil
}
