package sds

import (
	"fmt"

	"istio.io/istio/security/pkg/monitoring"
)

type MetricsManager struct {
	MetricsMap map[string]monitoring.Metric
}

func NewMetricsManager(metrics ...monitoring.Metric) *MetricsManager {
	newMetricsMap := map[string]monitoring.Metric{}
	for _, metric := range metrics {
		newMetricsMap[metric.Name()] = metric
		monitoring.MustRegisterViews(metric)
	}
	return &MetricsManager{
		newMetricsMap,
	}
}

func (m *MetricsManager) CreateGaugeMetric(name, description string, tags ...monitoring.Tag) {
	if _, found := m.MetricsMap[name]; !found {
		fmt.Println("HELLO. ADDED GAUGE " + name)
		m.MetricsMap[name] = monitoring.NewGauge(name, description, tags...)
		monitoring.MustRegisterViews(m.MetricsMap[name])
	}
}

func (m *MetricsManager) CreateSumMetric(name, description string, tags ...monitoring.Tag) {
	if _, found := m.MetricsMap[name]; !found {
		fmt.Println("HELLO. ADDED SUM " + name)
		m.MetricsMap[name] = monitoring.NewSum(name, description, tags...)
		monitoring.MustRegisterViews(m.MetricsMap[name])
	}
}

func (m *MetricsManager) RemoveMetrics(metricNames ...string) {
	fmt.Println("DELETED METRICS")
	for _, metricName := range metricNames {
		if metric, found := m.MetricsMap[metricName]; found {
			delete(m.MetricsMap, metricName)
			monitoring.UnregisterViews(metric)
		}
	}
}

func (m *MetricsManager) CreateSumMetricAndInc(name, description string, tags ...monitoring.Tag) {
	m.CreateSumMetric(name, description, tags...)
	m.MetricsMap[name].Increment()
}
