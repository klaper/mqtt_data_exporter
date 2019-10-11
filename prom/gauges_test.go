package prom

import (
	"testing"
)

func TestMetrics_RegisterMetric_Gauge_Count(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	initialLen := len(metrics.gauges)

	//when
	metrics.RegisterMetric(GAUGE, firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//then
	if len(metrics.gauges)-1 != initialLen {
		t.Errorf("New metric should be added [expected: %d, actual: %d]", initialLen+1, len(metrics.gauges))
	}
}

func TestMetrics_RegisterMetric_Gauge_Key(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)

	//when
	metrics.RegisterMetric(GAUGE, firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//then
	if _, ok := metrics.gauges[firstInputMetricsKey]; !ok {
		t.Errorf("Element \"%s\" was not found on metrics list", firstInputMetricsKey)
	}
}

func TestMetrics_RegisterMetric_Gauge_MetricExists(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	metrics.RegisterMetric(GAUGE, firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//when
	ok := metrics.RegisterMetric(
		GAUGE,
		firstInputMetricsKey,
		"TestMetrics_RegisterMetric_Gauge_MetricExists",
		inputMetricsDescription,
		inputLabelNames,
	)

	//then
	if ok {
		t.Errorf("RegisterMetric() = %v, want %v", ok, false)
	}
}

func TestMetrics_RegisterMetric_Gauge_MetricAdded(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	metrics.RegisterMetric(GAUGE, firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//when
	ok := metrics.RegisterMetric(
		GAUGE,
		secondInputMetricsKey,
		"TestMetrics_RegisterMetric_Gauge_MetricAdded",
		inputMetricsDescription,
		inputLabelNames,
	)

	//then
	if !ok {
		t.Errorf("RegisterMetric() = %v, want %v", ok, true)
	}
}
