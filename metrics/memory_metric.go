package metrics

import (
	"server-monitoring/model"

	"time"

	"github.com/shirou/gopsutil/mem"
)

type MemoryUsageInfo struct {
	Used        float64 // Использованное пространство в GB
	UsedPercent float64 // Процент использованного пространства
	Free        float64 // Свободное пространство в GB
}

type MemoryMetric struct {
	model.BaseMetric
}

func (m *MemoryMetric) Configure(base model.BaseMetric) {
	m.BaseMetric = base
}

func (m *MemoryMetric) GetName() string {
	return m.Name
}

func (m *MemoryMetric) GetType() string {
	return m.MetricType
}

func (m *MemoryMetric) IsActive() bool {
	return m.Active
}

func (m *MemoryMetric) GetPeriod() time.Duration {
	return m.Period * time.Second
}

func (m *MemoryMetric) GetRetentionPeriod() time.Duration {
	return m.RetentionPeriod * time.Second
}
func (m *MemoryMetric) GetFilePath() string {
	return m.FilePath
}

// Реализация метода Collect, который будет возвращать данные о CPU
func (m *MemoryMetric) Collect() (model.BaseMetricValue, error) {
	vmStat, err := mem.VirtualMemory()

	if err != nil {
		return model.BaseMetricValue{}, err
	}
	percentages := MemoryUsageInfo{

		Used:        float64(vmStat.Used) / 1024 / 1024 / 1024, // Преобразуем байты в GB
		UsedPercent: vmStat.UsedPercent,
		Free:        float64(vmStat.Free) / 1024 / 1024 / 1024, // Преобразуем байты в GB
	}

	return model.BaseMetricValue{
		Value:      percentages,
		ExecutedAt: time.Now(),
	}, nil
}
