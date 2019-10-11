package prom

import (
	"github.com/klaper_/mqtt_data_exporter/naming"
	"strings"
)

type MetricType int

const (
	COUNTER MetricType = 0
	GAUGE   MetricType = 1
)

type NamingService interface {
	TranslateDevice(deviceName string) (*naming.NamerDevice, bool)
}

type Metrics struct {
	counters          map[string]counterWithMetadata
	gauges            map[string]gaugeWithMetadata
	NamingService     NamingService
	metricsNamePrefix string
}

func NewMetrics(metricsNamePrefix string, namingService NamingService) *Metrics {
	return &Metrics{
		counters:          make(map[string]counterWithMetadata),
		gauges:            make(map[string]gaugeWithMetadata),
		NamingService:     namingService,
		metricsNamePrefix: strings.Trim(metricsNamePrefix, "_"),
	}
}

func (metrics *Metrics) RegisterMetric(metricsType MetricType, key string, name string, description string, labelNames []string) bool {
	switch metricsType {
	case COUNTER:
		return metrics.prepareCounter(key, name, description, labelNames)
	case GAUGE:
		return metrics.prepareGauge(key, name, description, labelNames)
	}
	return false
}

func (metrics *Metrics) prefixName(name string) string {
	return metrics.metricsNamePrefix + "_" + strings.Trim(name, "_")
}

func (metrics *Metrics) prepareLabelValues(labelNames []string, labelValues map[string]string) []string {
	result := make([]string, len(labelNames))
	for i, label := range labelNames {
		result[i] = labelValues[label]
	}
	return result
}

func (metrics *Metrics) appendRestrictedToValues(deviceName string, labels map[string]string) map[string]string {
	deviceInfo, ok := metrics.NamingService.TranslateDevice(deviceName)
	result := make(map[string]string)
	for k, v := range labels {
		result[k] = v
	}
	if !ok {
		deviceInfo = &naming.NamerDevice{Name: deviceName, Group: "", Device: deviceName}
	}
	result["device"] = deviceInfo.Device
	result["group"] = deviceInfo.Group
	result["friendly_name"] = deviceInfo.Name
	return result
}

func prepareLabelNames(labelNames []string) []string {
	result := restrictedLabelNames
	for label := range labelNames {
		l := labelNames[label]
		if contains(restrictedLabelNames, l) {
			l = "module_" + l
		}
		if contains(result, l) {
			continue
		}
		result = append(result, l)
	}
	return result
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

var restrictedLabelNames = []string{"device", "group", "friendly_name"}

