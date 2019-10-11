package prom

import "testing"

func TestMetrics_RegisterMetric_Counter_Count(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	initialLen := len(metrics.counters)

	//when
	metrics.RegisterMetric(COUNTER, inputMetricsKey, inputMetricsName, inputMetricsDescription, inputLabelNames)

	//then
	if len(metrics.counters)-1 != initialLen {
		t.Errorf("New metric should be added [expected: %d, actual: %d]", initialLen+1, len(metrics.counters))
	}
}

func TestMetrics_RegisterMetric_Counter_Key(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)

	//when
	metrics.RegisterMetric(COUNTER, inputMetricsKey, inputMetricsName, inputMetricsDescription, inputLabelNames)

	//then
	if _, ok := metrics.counters[inputMetricsKey]; !ok {
		t.Errorf("Element \"%s\" was not found on metrics list", inputMetricsKey)
	}
}

func TestMetrics_RegisterMetric_Counter_MetricExists(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	metrics.RegisterMetric(COUNTER, inputMetricsKey, inputMetricsName, inputMetricsDescription, inputLabelNames)

	//when
	ok := metrics.RegisterMetric(COUNTER, inputMetricsKey, inputMetricsName+"1", inputMetricsDescription+"1", inputLabelNames)

	//then
	if ok {
		t.Errorf("RegisterMetric() = %v, want %v", ok, false)
	}
}

func TestMetrics_RegisterMetric_Counter_MetricAdded(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	metrics.RegisterMetric(COUNTER, inputMetricsKey, inputMetricsName, inputMetricsDescription, inputLabelNames)

	//when
	ok := metrics.RegisterMetric(COUNTER, inputMetricsKey+"1", inputMetricsName+"1", inputMetricsDescription+"1", inputLabelNames)

	//then
	if !ok {
		t.Errorf("RegisterMetric() = %v, want %v", ok, true)
	}
}
