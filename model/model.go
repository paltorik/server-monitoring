package model

import (
	"time"
)

// Структура для данных метрики

// Базовая метрика - это метрика с базовой информацией
type BaseMetric struct {
	Name            string        `json:"name"`             // Имя метрики
	MetricType      string        `json:"metric_type"`      // Тип метрики (например, "cpu", "memory", "status")
	Active          bool          `json:"active"`           // Статус активности метрики (включена или выключена)
	Period          time.Duration `json:"period"`           // Периодичность сбора метрики
	RetentionPeriod time.Duration `json:"retention_period"` // Периодичность сбора метрики
	FilePath        string        `json:"file"`             // Периодичность сбора метрики
}

type BaseMetricValue struct {
	ExecutedAt time.Time   `json:"executed_at"`
	Value      MetricValue `json:"value"`
}

// Кастомное значение метрики
type MetricValue interface{}

// Структура, которая хранит значения метрики
type Metric struct {
	BaseMetric
	Values []MetricValue `json:"value"` // Значение метрики (число, объект или массив)
}

// Структура для серверной информации
type Server struct {
	Name        string    `json:"name"`         // Название сервера
	OS          string    `json:"os"`           // Операционная система
	RAMTotal    float64   `json:"ram_total"`    // Общее количество оперативной памяти
	DiskSize    float64   `json:"disk_size"`    // Размер диска
	Status      string    `json:"status"`       // Состояние: prod, dev, local
	LastChecked time.Time `json:"last_checked"` // Последний момент обновления
	Metrics     []Metric  `json:"metrics"`      // Метрики, которые собираются для этого сервера
}

type Config struct {
	Metrics []BaseMetric `json:"metrics"` // Список всех метрик
}
