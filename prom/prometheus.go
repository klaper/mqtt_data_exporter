package prom

import (
	"github.com/klaper_/mqtt_data_exporter/devices"
	"strings"
)

type DevicePropertiesProvider interface {
	GetProperties(deviceName string) (*devices.Properties, bool)
}

type Metrics struct {
	counters           map[string]counterWithMetadata
	gauges             map[string]gaugeWithMetadata
	propertiesProvider DevicePropertiesProvider
	metricsNamePrefix  string
}

func NewMetrics(metricsNamePrefix string, propertiesProvider DevicePropertiesProvider) *Metrics {
	return &Metrics{
		counters:           make(map[string]counterWithMetadata),
		gauges:             make(map[string]gaugeWithMetadata),
		propertiesProvider: propertiesProvider,
		metricsNamePrefix:  strings.Trim(metricsNamePrefix, "_"),
	}
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
	deviceInfo, ok := metrics.propertiesProvider.GetProperties(deviceName)
	result := make(map[string]string)
	for k, v := range labels {
		result[k] = v
	}
	if !ok {
		deviceInfo = &devices.Properties{Name: deviceName, Group: "", Device: deviceName, Sensors: make(map[string]string, 0)}
	}
	result["device"] = deviceInfo.Device
	result["group"] = deviceInfo.Group
	result["friendly_name"] = deviceInfo.Name
	if alias, ok := getSensorAlias(labels, deviceInfo); ok {
		result["sensor_alias"] = alias
	}

	return result
}

func getSensorAlias(labels map[string]string, deviceInfo *devices.Properties) (string, bool) {
	if name, has_sensor := labels["sensor_name"]; has_sensor {
		if sensor, has_alias := deviceInfo.Sensors[name]; has_alias {
			return sensor, true
		}
		return name, true
	}
	return "", false
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

var restrictedLabelNames = []string{"device", "group", "friendly_name", "sensor_alias"}
