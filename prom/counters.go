package prom

import "github.com/prometheus/client_golang/prometheus"

func (metrics *Metrics) Inc(key string, deviceName string, labels map[string]string) {
	counter, found := metrics.counters[key]
	if !found {
		return
	}
	var completedLabels = metrics.prepareLabelValues(counter.labels, metrics.appendRestrictedToValues(deviceName, labels))
	counter.metric.WithLabelValues(completedLabels...).Inc()
}

func (metrics *Metrics) prepareCounter(key string, name string, description string, labelNames []string) bool {
	_, ok := metrics.counters[key]
	if ok {
		return false
	}
	labels := prepareLabelNames(labelNames)
	metrics.counters[key] = counterWithMetadata{
		metric: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metrics.prefixName(name),
				Help: description,
			}, labels),
		labels: labels,
	}
	return true
}

type counterWithMetadata struct {
	metric *prometheus.CounterVec
	labels []string
}
