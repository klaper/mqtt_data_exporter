package prom

import "github.com/prometheus/client_golang/prometheus"

func (metrics *Metrics) RegisterGauge(key string, name string, description string, labelNames []string) bool {
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
	metrics.gaugeCleaner.RegisterGauge(key, metrics.gauges[key].metric)
	err := prometheus.Register(metrics.gauges[key].metric)
	return err == nil
}

func (metrics *Metrics) GaugeSet(key string, deviceName string, labels map[string]string, value float64) {
	counter, found := metrics.gauges[key]
	if !found {
		return
	}
	labelValues := metrics.appendRestrictedToValues(deviceName, labels)
	counter.metric.WithLabelValues(metrics.prepareLabelValues(counter.labels, labelValues)...).Set(value)
	metrics.gaugeCleaner.updateMetric(key, labelValues)
}

type gaugeWithMetadata struct {
	metric *prometheus.GaugeVec
	labels []string
}
