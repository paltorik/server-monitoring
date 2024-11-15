package metrics

import (
	"server-monitoring/model"

	"time"

	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerInfo struct {
	Name   string
	Status string
}

type DockerMetric struct {
	model.BaseMetric
}

func (m *DockerMetric) Configure(base model.BaseMetric) {
	m.BaseMetric = base
}

func (m *DockerMetric) GetName() string {
	return m.Name
}

func (m *DockerMetric) GetType() string {
	return m.MetricType
}

func (m *DockerMetric) IsActive() bool {
	return m.Active
}

func (m *DockerMetric) GetPeriod() time.Duration {
	return m.Period * time.Second
}

func (m *DockerMetric) GetRetentionPeriod() time.Duration {
	return m.RetentionPeriod * time.Second
}
func (m *DockerMetric) GetFilePath() string {
	return m.FilePath
}

// Реализация метода Collect, который будет возвращать данные о CPU
func (m *DockerMetric) Collect() (model.BaseMetricValue, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.43"), client.FromEnv)
	if err != nil {
		return model.BaseMetricValue{}, err
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return model.BaseMetricValue{}, err
	}
	// Получаем список всех контейнеров, включая остановленные

	// Создаем список DockerMetric для записи в лог
	var metrics []DockerInfo

	// Проходим по всем контейнерам и записываем их статус
	for _, container := range containers {
		status := "running"
		if container.State == "exited" {
			status = "stopped"
		}

		metrics = append(metrics, DockerInfo{
			Name:   container.Names[0],
			Status: status,
		})
	}
	cli.Close()
	return model.BaseMetricValue{
		Value:      metrics,
		ExecutedAt: time.Now(),
	}, nil

}
