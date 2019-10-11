package tasmota

import (
	broadcast "github.com/dustin/go-broadcast"
	"github.com/klaper_/mqtt_data_exporter/naming"
	"github.com/klaper_/mqtt_data_exporter/prom"
)

type TasmotaCollector struct {
	state  *prometheusTasmotaStateCollector
	sensor *prometheusTasmotaSensorCollector
}

func NewTasmotaCollector(prometheusTopicPrefix string, metricsStore *prom.Metrics, file string) *TasmotaCollector {
	return &TasmotaCollector{
		state:  newPrometheusTasmotaStateCollector(metricsStore),
		sensor: newPrometheusTasmotaSensorCollector(prometheusTopicPrefix, naming.NewNamer(file)),
	}
}

func (collector *TasmotaCollector) InitializeMessageReceiver(broadcaster broadcast.Broadcaster) {
	broadcaster.Register(collector.state.channel)
	broadcaster.Register(collector.sensor.channel)
	go collector.state.collector()
	go collector.sensor.collector()
}
