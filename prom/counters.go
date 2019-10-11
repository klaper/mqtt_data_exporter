package prom

import "github.com/prometheus/client_golang/prometheus"

func (metrics *Metrics) RegisterCounter(key string, name string, description string, labelNames []string) bool {
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
	err := prometheus.Register(metrics.counters[key].metric)
	return err == nil
}

func (metrics *Metrics) CounterInc(key string, deviceName string, labels map[string]string) {
	counter, found := metrics.counters[key]
	if !found {
		return
	}
	var completedLabels = metrics.prepareLabelValues(counter.labels, metrics.appendRestrictedToValues(deviceName, labels))
	counter.metric.WithLabelValues(completedLabels...).Inc()
}

type counterWithMetadata struct {
	metric *prometheus.CounterVec
	labels []string
}
