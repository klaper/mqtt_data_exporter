package prom

import "testing"

func TestMetrics_RegisterCounter_Count(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	initialLen := len(metrics.counters)

	//when
	metrics.RegisterCounter(firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//then
	if len(metrics.counters)-1 != initialLen {
		t.Errorf("New metric should be added [expected: %d, actual: %d]", initialLen+1, len(metrics.counters))
	}
}

func TestMetrics_RegisterCounter_Key(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)

	//when
	metrics.RegisterCounter(firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//then
	if _, ok := metrics.counters[firstInputMetricsKey]; !ok {
		t.Errorf("Element \"%s\" was not found on metrics list", firstInputMetricsKey)
	}
}

func TestMetrics_RegisterCounter_MetricExists(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	metrics.RegisterCounter(firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//when
	ok := metrics.RegisterCounter(
		firstInputMetricsKey,
		"TestMetrics_RegisterCounter_MetricExists",
		inputMetricsDescription,
		inputLabelNames,
		)

	//then
	if ok {
		t.Errorf("RegisterMetric() = %v, want %v", ok, false)
	}
}

func TestMetrics_RegisterCounter_MetricAdded(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	metrics.RegisterCounter(firstInputMetricsKey, firstInputMetricsName, inputMetricsDescription, inputLabelNames)

	//when
	ok := metrics.RegisterCounter(
		secondInputMetricsKey,
		"TestMetrics_RegisterCounter_MetricAdded",
		inputMetricsDescription,
		inputLabelNames,
		)

	//then
	if !ok {
		t.Errorf("RegisterMetric() = %v, want %v", ok, true)
	}
}
