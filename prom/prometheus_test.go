package prom

import (
	"github.com/klaper_/mqtt_data_exporter/devices"
	"reflect"
	"testing"
)

var (
	inputDeviceName          = "SomeDeviceName"
	firstInputMetricsKey     = "SomeMetricKey"
	secondInputMetricsKey    = "OtherMetricKey"
	firstInputMetricsName    = "SomeMetricName"
	secondInputMetricsName   = "OtherMetricName"
	inputMetricsDescription  = "SomeMetricDescription"
	inputLabelNames          = []string{"label1", "label2"}
	processLabelsBenchResult []string
	prefixNameBenchResult    string
	restrictedLabels         = []string{"device", "group", "friendly_name"}
	expectedDeviceName       = "deviceName"
	expectedDeviceGroup      = "deviceGroup"
	expectedDeviceFName      = "deviceFName"
)

//BenchmarkMetrics_prefixName-8      	 5568180	       189 ns/op
func BenchmarkMetrics_prefixName(b *testing.B) {
	metrics := NewMetrics("_prefix_", nil)
	var r string
	for i := 0; i < b.N; i++ {
		r = metrics.prefixName("_n_a_m_e_")
	}
	prefixNameBenchResult = r
}

//BenchmarkMetrics_processLabels-8   	 4051521	       300 ns/op
func BenchmarkMetrics_prepareLabelNames(b *testing.B) {
	var r []string
	for i := 0; i < b.N; i++ {
		r = prepareLabelNames([]string{"name", "device", "friendly", "name"})
	}
	processLabelsBenchResult = r
}

func TestMetrics_prefixName(t *testing.T) {
	type fields struct {
		counters          map[string]counterWithMetadata
		namer             DevicePropertiesProvider
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
				counters:           tt.fields.counters,
				propertiesProvider: tt.fields.namer,
				metricsNamePrefix:  tt.fields.metricsNamePrefix,
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
		namer             DevicePropertiesProvider
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

func Test_prepareLabelValues(t *testing.T) {
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
			if got := prepareLabelNames(tt.args.labelNames); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareLabelNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_prepareLabelValues(t *testing.T) {
	type fields struct {
		counters          map[string]counterWithMetadata
		namer             DevicePropertiesProvider
		metricsNamePrefix string
	}
	type args struct {
		labelNames  []string
		labelValues map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			"empty labelNames",
			fields{},
			args{
				labelNames:  []string{},
				labelValues: map[string]string{},
			},
			[]string{},
		},
		{
			"not empty",
			fields{},
			args{
				labelNames:  []string{"2", "1", "4", "3"},
				labelValues: map[string]string{"1": "a", "2": "b", "3": "c", "4": "d", "5": "e"},
			},
			[]string{"b", "a", "d", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := &Metrics{
				counters:           tt.fields.counters,
				propertiesProvider: tt.fields.namer,
				metricsNamePrefix:  tt.fields.metricsNamePrefix,
			}
			if got := metrics.prepareLabelValues(tt.args.labelNames, tt.args.labelValues); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareLabelValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_prepareLabelValues_missingValue(t *testing.T) {
	inputLabelNames := []string{"a"}
	inputLabelValues := map[string]string{}
	want := []string{""}

	metrics := &Metrics{
		counters:           nil,
		propertiesProvider: nil,
		metricsNamePrefix:  "",
	}

	if got := metrics.prepareLabelValues(inputLabelNames, inputLabelValues); !reflect.DeepEqual(got, want) {
		t.Errorf("prepareLabelValues() = %v, want %v", got, want)
	}
}

func TestMetrics_appendRestrictedToValues(t *testing.T) {
	var namingService DevicePropertiesProvider = TestNamingService{}
	type fields struct {
		counters          map[string]counterWithMetadata
		namer             DevicePropertiesProvider
		metricsNamePrefix string
	}
	type args struct {
		deviceName string
		labels     map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			"empty values",
			fields{namer: namingService},
			args{deviceName: inputDeviceName, labels: map[string]string{}},
			map[string]string{
				"device":        expectedDeviceName,
				"group":         expectedDeviceGroup,
				"friendly_name": expectedDeviceFName,
			},
		},
		{
			"empty values",
			fields{namer: namingService},
			args{deviceName: inputDeviceName, labels: map[string]string{"new": "way"}},
			map[string]string{
				"device":        expectedDeviceName,
				"group":         expectedDeviceGroup,
				"friendly_name": expectedDeviceFName,
				"new":           "way",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := &Metrics{
				counters:           tt.fields.counters,
				propertiesProvider: tt.fields.namer,
				metricsNamePrefix:  tt.fields.metricsNamePrefix,
			}
			if got := metrics.appendRestrictedToValues(tt.args.deviceName, tt.args.labels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendRestrictedToValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_appendRestrictedToValues_WrongDevice(t *testing.T) {
	deviceName := "wrongDeviceName"
	expected := map[string]string{"device": deviceName, "group": "", "friendly_name": deviceName}

	metrics := &Metrics{propertiesProvider: TestNamingService{}}
	result := make(map[string]string)

	if got := metrics.appendRestrictedToValues(deviceName, result); !reflect.DeepEqual(got, expected) {
		t.Errorf("appendRestrictedToValues() = %v, want %v", got, expected)
	}
}

type TestNamingService struct{}

func (t TestNamingService) GetProperties(deviceName string) (*devices.Properties, bool) {
	if deviceName == inputDeviceName {
		return &devices.Properties{Name: expectedDeviceFName, Device: expectedDeviceName, Group: expectedDeviceGroup}, true
	}
	return nil, false
}
