package prom

import (
	"github.com/klaper_/mqtt_data_exporter/naming"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

var restrictedSet = map[string]struct{}{
	"device":        {},
	"group":         {},
	"friendly_name": {},
}

type Metrics struct {
	counters          map[string]*prometheus.CounterVec
	namer             *naming.Namer
	metricsNamePrefix string
}

func NewMetrics(metricsNamePrefix string, namer *naming.Namer) *Metrics {
	return &Metrics{
		counters:          make(map[string]*prometheus.CounterVec),
		namer:             namer,
		metricsNamePrefix: strings.Trim(metricsNamePrefix, "_"),
	}
}

func (metrics *Metrics) prefixName(name string) string {
	return metrics.metricsNamePrefix + "_" + strings.Trim(name, "_")
}

func (metrics *Metrics) RegisterMetric(key string, name string, description string, labelNames []string) bool {
	_, ok := metrics.counters[key]
	if ok {
		return false
	}
	metrics.counters[key] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metrics.prefixName(name),
			Help: description,
		}, labelNames)
	return true
}

func processLabels(labelNames []string) []string {
	resultSet := make(map[string]struct{})
	for label := range labelNames {
		l := labelNames[label]
		if _, ok := restrictedSet[l]; ok {
			l = "module_" + l
		}
		if _, ok := resultSet[l]; ok {
			continue
		}
		resultSet[l] = struct{}{}
	}
	result := make([]string, 0, len(resultSet))
	for k := range resultSet {
		result = append(result, k)
	}
	for k := range restrictedSet {
		result = append(result, k)
	}
	return result
}
