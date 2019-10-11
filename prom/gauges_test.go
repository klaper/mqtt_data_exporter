package prom

import (
	"testing"
)

func TestMetrics_RegisterGauge_Count(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	initialLen := len(metrics.gauges)

	//when
	metrics.RegisterGauge(firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//then
	if len(metrics.gauges)-1 != initialLen {
		t.Errorf("New metric should be added [expected: %d, actual: %d]", initialLen+1, len(metrics.gauges))
	}
}

func TestMetrics_RegisterGauge_Key(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)

	//when
	metrics.RegisterGauge(firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//then
	if _, ok := metrics.gauges[firstInputMetricsKey]; !ok {
		t.Errorf("Element \"%s\" was not found on metrics list", firstInputMetricsKey)
	}
}

func TestMetrics_RegisterGauge_MetricExists(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	metrics.RegisterGauge(firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//when
	ok := metrics.RegisterMetric(
		GAUGE,
		firstInputMetricsKey,
		"TestMetrics_RegisterGauge_MetricExists",
		inputMetricsDescription,
		inputLabelNames,
	)

	//then
	if ok {
		t.Errorf("RegisterMetric() = %v, want %v", ok, false)
	}
}

func TestMetrics_RegisterGauge_MetricAdded(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	metrics.RegisterGauge(firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//when
	ok := metrics.RegisterMetric(
		GAUGE,
		secondInputMetricsKey,
		"TestMetrics_RegisterGauge_MetricAdded",
		inputMetricsDescription,
		inputLabelNames,
	)

	//then
	if !ok {
		t.Errorf("RegisterMetric() = %v, want %v", ok, true)
	}
}
