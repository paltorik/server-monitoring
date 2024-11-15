package metrics

import (
	"server-monitoring/model"

	"time"

	"github.com/shirou/gopsutil/disk"
)

type DiskUsageInfo struct {
	Used        float64 // Использованное пространство в GB
	UsedPercent float64 // Процент использованного пространства
	Free        float64 // Свободное пространство в GB
}

type DickMertic struct {
	model.BaseMetric
}

func (m *DickMertic) Configure(base model.BaseMetric) {
	m.BaseMetric = base
}

func (m *DickMertic) GetName() string {
	return m.Name
}

func (m *DickMertic) GetType() string {
	return m.MetricType
}

func (m *DickMertic) IsActive() bool {
	return m.Active
}

func (m *DickMertic) GetPeriod() time.Duration {
	return m.Period * time.Second
}

func (m *DickMertic) GetRetentionPeriod() time.Duration {
	return m.RetentionPeriod * time.Second
}
func (m *DickMertic) GetFilePath() string {
	return m.FilePath
}

// Реализация метода Collect, который будет возвращать данные о CPU
func (m *DickMertic) Collect() (model.BaseMetricValue, error) {
	usage, err := disk.Usage("/")
	if err != nil {
		return model.BaseMetricValue{}, err
	}
	percentages := DiskUsageInfo{

		Used:        float64(usage.Used) / 1024 / 1024 / 1024, // Преобразуем байты в GB
		UsedPercent: usage.UsedPercent,
		Free:        float64(usage.Free) / 1024 / 1024 / 1024, // Преобразуем байты в GB
	}

	return model.BaseMetricValue{
		Value:      percentages,
		ExecutedAt: time.Now(),
	}, nil
}
