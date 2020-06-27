package prom

import (
	"crypto"
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
	"testing"
	"time"
)

var inputKey = "key"
var inputLabels = map[string]string{"label": "value"}
var startDate = time.Date(2020, 02, 15, 17, 46, 32, 0, time.Local)

type testClock struct {
	current time.Time
}

func (tc testClock) Now() time.Time                      { return tc.current }
func (testClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
func (tc *testClock) MoveForward(duration time.Duration) { tc.current = tc.current.Add(duration) }
func (tc *testClock) Reset()                             { tc.current = startDate }

type testGaugeVec struct {
	deleted bool
}

func (tgv *testGaugeVec) Delete(labels prometheus.Labels) bool {
	tgv.deleted = true
	return true
}
func (testGaugeVec) GetMetricWith(labels prometheus.Labels) (prometheus.Gauge, error) {
	return prometheus.NewGauge(prometheus.GaugeOpts{Name: "test"}), nil
}

func Test_gaugeCleaner_updateMetric(t *testing.T) {
	//given
	cleaner := &gaugeCleaner{input: make(chan metricsDescriptor, 1), clk: &testClock{current: startDate}, timeout: 1}
	defer cleaner.Close()

	//when
	cleaner.updateMetric(inputKey, inputLabels)

	//then
	if len(cleaner.input) != 1 {
		t.Error("updateMetrics(): expected to publish input on channel")
	}

	result := <-cleaner.input

	if result.key != inputKey {
		t.Errorf("updateMetrics(): expected key to be: \"%s\" but was \"%s\"", inputKey, result.key)
	}

	if !reflect.DeepEqual(result.labels, inputLabels) {
		t.Errorf("updateMetrics(): expected labels to be: \"%+v\" but was \"%+v\"", inputLabels, result.labels)
	}
}

func Test_gaugeCleaner_updateMetric_shouldSkipWhenTimeoutIsZero(t *testing.T) {
	//given
	cleaner := &gaugeCleaner{input: make(chan metricsDescriptor, 1), clk: &testClock{current: startDate}}
	defer cleaner.Close()

	//when
	cleaner.updateMetric(inputKey, inputLabels)

	//then
	if len(cleaner.input) != 0 {
		t.Error("updateMetrics(): expected to skip publishing on timeout 0")
	}
}

func Test_gaugeCleaner_receive(t *testing.T) {
	//given
	hash := calculateHash(inputKey, inputLabels)
	input := metricsDescriptor{key: inputKey, labels: inputLabels}

	cleaner := NewGaugeCleaner("")
	cleaner.clk = &testClock{current: startDate}

	cleaner.input <- input
	cleaner.Close()

	//when
	cleaner.receive()

	//then
	if updated, ok := cleaner.updates[hash]; !ok {
		t.Errorf("receive(): expected to set value for \"%x\"", hash)
	} else if updated.updateTime != startDate {
		t.Errorf("receive(): expected to set update time \"%s\" for \"%x\"", startDate, hash)
	} else if !reflect.DeepEqual(updated.metricsDescriptor, input) {
		t.Errorf("receive(): expected to set descriptor to \"%v\" for \"%x\"", input, hash)
	}
}

func Test_gaugeCleaner_receive_loop(t *testing.T) {
	//given
	input := metricsDescriptor{key: inputKey, labels: inputLabels}

	cleaner := NewGaugeCleaner("")
	cleaner.clk = &testClock{current: startDate}

	cleaner.input <- input
	cleaner.input <- input
	cleaner.Close()

	//when
	cleaner.receive()

	//then
	if len(cleaner.input) > 0 {
		t.Errorf("receive(): expected to consume 2 elements")
	}
}

func Test_gaugeCleaner_RegisterGauge(t *testing.T) {
	//given
	gauge := &testGaugeVec{}
	cleaner := NewGaugeCleaner("")
	defer cleaner.Close()
	cleaner.clk = &testClock{current: startDate}

	//when
	cleaner.RegisterGauge(inputKey, gauge)

	//then
	if result, ok := cleaner.metrics[inputKey]; !ok {
		t.Errorf("RegisterGauge(): expected to register gauge [%+v] in metrics as %s", gauge, inputKey)
	} else if result != gauge {
		t.Errorf("RegisterGauge(): expected to register gauge [%+v] but [%+v] was registered", gauge, result)
	}
}

func Test_gaugeCleaner_clean(t *testing.T) {
	//given
	gauge := &testGaugeVec{deleted: false}
	hash := calculateHash(inputKey, inputLabels)
	cleaner := NewGaugeCleanerWithTimeout(1, "")
	defer cleaner.Close()
	clock := &testClock{current: startDate}
	cleaner.clk = clock

	cleaner.metrics[inputKey] = gauge
	cleaner.updates[hash] = metricsUpdateDescriptor{
		updateTime:        clock.Now(),
		metricsDescriptor: metricsDescriptor{key: inputKey, labels: inputLabels},
	}

	//when
	clock.MoveForward(cleaner.timeout + (2 * time.Second))
	//and
	cleaner.clean()

	//then
	if !gauge.deleted {
		t.Error("clean(): gauge should be deleted")
	}
	if _, ok := cleaner.updates[hash]; ok {
		t.Error("clean(): updates should be deleted")
	}
}

func Test_gaugeCleaner_clean_shouldNotDeleteWithZeroTimeout(t *testing.T) {
	//given
	gauge := &testGaugeVec{deleted: false}
	cleaner := NewGaugeCleaner("")
	defer cleaner.Close()
	clock := &testClock{current: startDate}
	cleaner.clk = clock

	cleaner.metrics[inputKey] = gauge
	cleaner.updates[calculateHash(inputKey, inputLabels)] = metricsUpdateDescriptor{
		updateTime:        clock.Now(),
		metricsDescriptor: metricsDescriptor{key: inputKey, labels: inputLabels},
	}

	//when
	clock.MoveForward(cleaner.timeout + (2 * time.Second))
	//and
	cleaner.clean()

	//then
	if gauge.deleted {
		t.Error("clean(): gauge should not be deleted")
	}
}

func TestCalculateHash(t *testing.T) {
	//given
	key := "key"
	labels := map[string]string{
		"a": "b",
		"c": "d",
		"e": "a",
	}
	expectedResult := string(crypto.SHA1.New().Sum([]byte("key b d a")))

	//when
	result := calculateHash(key, labels)

	//then
	if result != expectedResult {
		t.Errorf("calculateHash() == \"%x\", want \"%x\"", result, expectedResult)
	}
}

func TestCalculateHash_order(t *testing.T) {
	//given
	key := "key"
	labelsFirst := map[string]string{
		"a": "b",
		"c": "d",
		"e": "a",
	}
	labelsSecond := map[string]string{
		"e": "a",
		"c": "d",
		"a": "b",
	}

	//when
	firstResult := calculateHash(key, labelsFirst)
	secondResult := calculateHash(key, labelsSecond)

	//then
	if firstResult != secondResult {
		t.Errorf("calculateHash(firstLabels) != calculateHash(secondLabels)")
	}
}
func BenchmarkCalculateHash(b *testing.B) {
	//BenchmarkCalculateHash-8   	  616411	      1906 ns/op
	key := "key"
	labels := map[string]string{
		"a": "b",
		"c": "d",
		"e": "f",
	}

	for i := 0; i < b.N; i++ {
		calculateHash(key, labels)
	}
}
