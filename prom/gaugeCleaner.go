package prom

import (
	"crypto"
	"fmt"
	"github.com/klaper_/mqtt_data_exporter/logger"
	"github.com/prometheus/client_golang/prometheus"
	"sort"
	"sync"
	"time"
)

const (
	MetricsUpdatesCount = "cleaner_update_count"
	MetricsCleanTime    = "cleaner_clean_time"
	MetricsCleanCount   = "cleaner_clean_count"
)

var (
	MetricsUpdatesCounter prometheus.Counter
	MetricsCleanGauge     prometheus.Gauge
	MetricsCleanCounter   prometheus.Counter
)

type clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type realClock struct{}

func (realClock) Now() time.Time                         { return time.Now() }
func (realClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

type metricsDescriptor struct {
	key    string
	labels map[string]string
}

type metricsUpdateDescriptor struct {
	metricsDescriptor metricsDescriptor
	updateTime        time.Time
}

type gaugeVector interface {
	Delete(labels prometheus.Labels) bool
	GetMetricWith(labels prometheus.Labels) (prometheus.Gauge, error)
}

type gaugeCleaner struct {
	updates       map[string]metricsUpdateDescriptor
	metrics       map[string]gaugeVector
	input         chan metricsDescriptor
	clk           clock
	timeout       time.Duration
	lock          sync.Mutex
	metricsPrefix string
}

func NewGaugeCleaner(metricsPrefix string) *gaugeCleaner {
	return &gaugeCleaner{
		updates:       make(map[string]metricsUpdateDescriptor),
		metrics:       make(map[string]gaugeVector),
		input:         make(chan metricsDescriptor, 10),
		clk:           &realClock{},
		timeout:       0,
		metricsPrefix: metricsPrefix,
	}
}

func NewGaugeCleanerWithTimeout(timeout time.Duration, metricsPrefix string) *gaugeCleaner {
	return &gaugeCleaner{
		updates:       make(map[string]metricsUpdateDescriptor),
		metrics:       make(map[string]gaugeVector),
		input:         make(chan metricsDescriptor, 10),
		clk:           &realClock{},
		timeout:       timeout,
		metricsPrefix: metricsPrefix,
	}
}

func (gc *gaugeCleaner) Close() error {
	close(gc.input)
	return nil
}

func (gc *gaugeCleaner) RegisterGauge(key string, vector gaugeVector) {
	logger.Debug("gauge_cleaner", "Registering new gauge %+v under %s", vector, key)
	gc.metrics[key] = vector
}

func (gc *gaugeCleaner) Run() {
	gc.prepareMetrics()
	if gc.timeout == 0 {
		logger.Warn("gauge_cleaner", "Timeout is set to 0; disabling cleaner;")
		return
	}
	go gc.receive()
	go gc.cleaner()
}

func (gc *gaugeCleaner) prepareMetrics() {
	MetricsUpdatesCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: gc.metricsPrefix + "_" + MetricsUpdatesCount,
			Help: "Count of registered metrics update",
		})
	MetricsCleanGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: gc.metricsPrefix + "_" + MetricsCleanTime,
			Help: "Metrics clean time in seconds - lock time only",
		})
	MetricsCleanCounter = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: gc.metricsPrefix + "_" + MetricsCleanCount,
			Help: "Removed metrics count",
		})

	err := prometheus.Register(MetricsCleanCounter)
	logger.Debug("gauge_cleaner", "%s register error: %v", MetricsCleanCount, err)
	err = prometheus.Register(MetricsUpdatesCounter)
	logger.Debug("gauge_cleaner", "%s register error: %v", MetricsUpdatesCount, err)
	err = prometheus.Register(MetricsCleanGauge)
	logger.Debug("gauge_cleaner", "%s register error: %v", MetricsCleanTime, err)
}

func (gc *gaugeCleaner) updateMetric(key string, labels map[string]string) {
	if gc.timeout == 0 {
		//this is important due to channel receiver being disabled
		logger.Debug("gauge_cleaner", "Skipped update due to 0 timeout for %s, %+v", key, labels)
		return
	}
	MetricsUpdatesCounter.Inc()
	gc.input <- metricsDescriptor{key: key, labels: labels}
}

func (gc *gaugeCleaner) receive() {
	for update := range gc.input {
		logger.Debug("gauge_cleaner", "Received new update %+v", update)
		gc.lock.Lock()
		gc.updates[calculateHash(update.key, update.labels)] = metricsUpdateDescriptor{
			updateTime:        gc.clk.Now(),
			metricsDescriptor: update,
		}
		gc.lock.Unlock()
	}
}

func (gc *gaugeCleaner) cleaner() {
	for {
		select {
		case <-gc.clk.After(5 * time.Second):
			gc.clean()
		}
	}
}

func (gc *gaugeCleaner) clean() {
	logger.Debug("gauge_cleaner", "Running cleanup")
	if gc.timeout == 0 {
		logger.Debug("gauge_cleaner", "Cleanup skipped due to timeout set to 0")
		return
	}
	expired := make([]string, 0)
	for key, descriptor := range gc.updates {
		if descriptor.updateTime.Add(gc.timeout).Before(gc.clk.Now()) {
			expired = append(expired, key)
		}
	}
	if len(expired) == 0 {
		logger.Debug("gauge_cleaner", "Nothing to clean.. skipping")
		return
	}
	logger.Debug("gauge_cleaner", "Got %d entities to clean", len(expired))
	//first finding then removing expired metrics optimises lock time
	start := time.Now()
	gc.lock.Lock()
	for _, key := range expired {
		descriptor := gc.updates[key]
		if descriptor.updateTime.Add(gc.timeout).Before(gc.clk.Now()) {
			logger.Info("gauge_cleaner", "Cleaning %s metrics for %+v", descriptor.metricsDescriptor.key, descriptor.metricsDescriptor.labels)
			gc.metrics[descriptor.metricsDescriptor.key].Delete(descriptor.metricsDescriptor.labels)
			delete(gc.updates, key)
			MetricsCleanCounter.Inc()
		}
	}
	gc.lock.Unlock()
	end := time.Now()
	MetricsCleanGauge.Set(end.Sub(start).Seconds())
	logger.Debug("gauge_cleaner", "Cleaning completed after %s", end.Sub(start))
}

func calculateHash(key string, labels map[string]string) string {
	sortedKeys := make([]string, 0)
	for key, _ := range labels {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	hash := key
	for _, value := range sortedKeys {
		hash = hash + fmt.Sprintf(" %s", labels[value])
	}
	return string(crypto.SHA1.New().Sum([]byte(hash)))
}
