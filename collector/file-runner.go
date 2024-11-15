package collector

import (
	"fmt"
	"time"
)

func StartMetricWriteWorker(fileWriter *FileWriter, writeInterval time.Duration, stopChan chan struct{}) {
	go func() {
		ticker := time.NewTicker(writeInterval)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				fmt.Println("Воркер записи данных остановлен")
				return
			case <-ticker.C:
				// Запись данных из буфера в файлы
				if err := fileWriter.WriteToFile(); err != nil {
					fmt.Printf("Ошибка записи данных в файл: %v\n", err)
				}
			}
		}
	}()
}

func StartMetricCleanWorker(fileWriter *FileWriter, cleanInterval time.Duration, stopChan chan struct{}) {
	go func() {
		ticker := time.NewTicker(cleanInterval)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				fmt.Println("Воркер очистки данных остановлен")
				return
			case <-ticker.C:
				// Очистка старых данных
				if err := fileWriter.CleanOldData(); err != nil {
					fmt.Printf("Ошибка очистки старых данных: %v\n", err)
				}
			}
		}
	}()
}
