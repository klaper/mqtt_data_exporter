package prom

import "github.com/prometheus/client_golang/prometheus"

func (metrics *Metrics) RegisterGauge( key string, name string, description string, labelNames []string) bool {
	_, ok := metrics.gauges[key]
	if ok {
		return false
	}
	labels := prepareLabelNames(labelNames)
	metrics.gauges[key] = gaugeWithMetadata{
		metric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: metrics.prefixName(name),
				Help: description,
			}, labels),
		labels: labels,
	}
	err := prometheus.Register(metrics.gauges[key].metric)
	return err == nil
}

func (metrics *Metrics) Set(key string, deviceName string, labels map[string]string, value float64) {
	counter, found := metrics.gauges[key]
	if !found {
		return
	}
	var completedLabels = metrics.prepareLabelValues(counter.labels, metrics.appendRestrictedToValues(deviceName, labels))
	counter.metric.WithLabelValues(completedLabels...).Set(value)
}

type gaugeWithMetadata struct {
	metric *prometheus.GaugeVec
	labels []string
}
