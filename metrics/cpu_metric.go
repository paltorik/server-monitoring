package metrics

import (
	"server-monitoring/model"

	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type CPUMetric struct {
	model.BaseMetric
}

func (m *CPUMetric) Configure(base model.BaseMetric) {
	m.BaseMetric = base
}

func (m *CPUMetric) GetName() string {
	return m.Name
}

func (m *CPUMetric) GetType() string {
	return m.MetricType
}

func (m *CPUMetric) IsActive() bool {
	return m.Active
}

func (m *CPUMetric) GetPeriod() time.Duration {
	return m.Period * time.Second
}

func (m *CPUMetric) GetRetentionPeriod() time.Duration {
	return m.RetentionPeriod * time.Second
}
func (m *CPUMetric) GetFilePath() string {
	return m.FilePath
}

// Реализация метода Collect, который будет возвращать данные о CPU
func (m *CPUMetric) Collect() (model.BaseMetricValue, error) {
	percentages, err := cpu.Percent(0, false)

	if err != nil {
		return model.BaseMetricValue{}, err
	}

	if len(percentages) > 0 {
		return model.BaseMetricValue{
			Value:      percentages[0],
			ExecutedAt: time.Now(),
		}, nil
	}
	return model.BaseMetricValue{}, fmt.Errorf("не удалось получить данные о загрузке CPU")
}
