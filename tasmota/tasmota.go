package tasmota

import (
	broadcast "github.com/dustin/go-broadcast"
	"github.com/klaper_/mqtt_data_exporter/prom"
)

type TasmotaCollector struct {
	state  *prometheusTasmotaStateCollector
	sensor *prometheusTasmotaSensorCollector
}

func NewTasmotaCollector(metricsStore *prom.Metrics) *TasmotaCollector {
	return &TasmotaCollector{
		state:  newPrometheusTasmotaStateCollector(metricsStore),
		sensor: newPrometheusTasmotaSensorCollector(metricsStore),
	}
}

func (collector *TasmotaCollector) InitializeMessageReceiver(broadcaster broadcast.Broadcaster) {
	broadcaster.Register(collector.state.channel)
	broadcaster.Register(collector.sensor.channel)
	go collector.state.collector()
	go collector.sensor.collector()
}
