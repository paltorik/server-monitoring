package metrics

import (
	"server-monitoring/model"
	"time"
)

// Интерфейс для метрик
type Metric interface {
	GetName() string
	GetType() string
	Configure(base model.BaseMetric)
	IsActive() bool
	GetPeriod() time.Duration
	GetRetentionPeriod() time.Duration
	GetFilePath() string
	Collect() (model.BaseMetricValue, error) // Метод для сбора данных метрики
}
