package prom

import (
	"github.com/klaper_/mqtt_data_exporter/naming"
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

var (
	inputMetricsKey          = "SomeMetricKey"
	inputMetricsName         = "SomeMetricName"
	inputMetricsDescription  = "SomeMetricDescription"
	inputLabelNames          = []string{"label1", "label2"}
	processLabelsBenchResult []string
	prefixNameBenchResult    string
	restrictedLabels         = []string{"device", "group", "friendly_name"}
)

func TestMetrics_RegisterMetric_Count(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	initialLen := len(metrics.counters)

	//when
	metrics.RegisterMetric(inputMetricsKey, inputMetricsName, inputMetricsDescription, inputLabelNames)

	//then
	if len(metrics.counters)-1 != initialLen {
		t.Errorf("New metric should be added [expected: %d, actual: %d]", initialLen+1, len(metrics.counters))
	}
}

func TestMetrics_RegisterMetric_Key(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)

	//when
	metrics.RegisterMetric(inputMetricsKey, inputMetricsName, inputMetricsDescription, inputLabelNames)

	//then
	if _, ok := metrics.counters[inputMetricsKey]; !ok {
		t.Errorf("Element \"%s\" was not found on metrics list", inputMetricsKey)
	}
}

func TestMetrics_RegisterMetric_MetricExists(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	metrics.RegisterMetric(inputMetricsKey, inputMetricsName, inputMetricsDescription, inputLabelNames)

	//when
	ok := metrics.RegisterMetric(inputMetricsKey, inputMetricsName+"1", inputMetricsDescription+"1", inputLabelNames)

	//then
	if ok {
		t.Errorf("RegisterMetric() = %v, want %v", ok, false)
	}
}

func TestMetrics_RegisterMetric_MetricAdded(t *testing.T) {
	//given
	metrics := NewMetrics("", nil)
	metrics.RegisterMetric(inputMetricsKey, inputMetricsName, inputMetricsDescription, inputLabelNames)

	//when
	ok := metrics.RegisterMetric(inputMetricsKey+"1", inputMetricsName+"1", inputMetricsDescription+"1", inputLabelNames)

	//then
	if !ok {
		t.Errorf("RegisterMetric() = %v, want %v", ok, true)
	}
}

func TestMetrics_prefixName(t *testing.T) {
	type fields struct {
		counters          map[string]*prometheus.CounterVec
		namer             *naming.Namer
		metricsNamePrefix string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"simple join with no underscore", fields{metricsNamePrefix: "prefix"}, args{name: "name"}, "prefix_name"},
		{"name with underscore postfix", fields{metricsNamePrefix: "prefix"}, args{name: "name_"}, "prefix_name"},
		{"name with underscore prefix", fields{metricsNamePrefix: "prefix"}, args{name: "_name"}, "prefix_name"},
		{"lots of with underscores", fields{metricsNamePrefix: "prefix"}, args{name: "_n_a_m_e_"}, "prefix_n_a_m_e"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := &Metrics{
				counters:          tt.fields.counters,
				namer:             tt.fields.namer,
				metricsNamePrefix: tt.fields.metricsNamePrefix,
			}
			if got := metrics.prefixName(tt.args.name); got != tt.want {
				t.Errorf("prefixName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMetrics_PrefixTrim(t *testing.T) {
	type args struct {
		metricsNamePrefix string
		namer             *naming.Namer
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"prefix with underscore postfix", args{metricsNamePrefix: "prefix_"}, "prefix"},
		{"prefix with underscore prefix", args{metricsNamePrefix: "_prefix"}, "prefix"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetrics(tt.args.metricsNamePrefix, tt.args.namer); got.metricsNamePrefix != tt.want {
				t.Errorf("NewMetrics() = %v, want %v", got.metricsNamePrefix, tt.want)
			}
		})
	}
}

func Test_processLabels(t *testing.T) {
	type args struct {
		labelNames []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"empty set", args{labelNames: []string{}}, restrictedLabels},
		{"adding sets", args{labelNames: []string{"Label1"}}, append(restrictedLabels, "Label1")},
		{"repeating input", args{labelNames: []string{"Label1", "Label1"}}, append(restrictedLabels, "Label1")},
		{
			"adding sets and renaming to module_",
			args{labelNames: []string{restrictedLabels[0]}},
			append(restrictedLabels, "module_"+restrictedLabels[0]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processLabels(tt.args.labelNames); !arrayEqual(got, tt.want) {
				t.Errorf("processLabels() = %v, want %v", got, tt.want)
			}
		})
	}
}

//BenchmarkMetrics_prefixName-8      	 5568180	       189 ns/op
func BenchmarkMetrics_prefixName(b *testing.B) {
	metrics := NewMetrics("_prefix_", nil)
	var r string
	for i := 0; i < b.N; i++ {
		r = metrics.prefixName("_n_a_m_e_")
	}
	prefixNameBenchResult = r
}

//BenchmarkMetrics_processLabels-8   	 1651156	       713 ns/op
func BenchmarkMetrics_processLabels(b *testing.B) {
	var r []string
	for i := 0; i < b.N; i++ {
		r = processLabels([]string{"name", "device", "friendly", "name"})
	}
	processLabelsBenchResult = r
}

func arrayEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if !contains(b, v) {
			return false
		}
	}
	return true
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
