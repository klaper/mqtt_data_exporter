package tasmota

import (
	broadcast "github.com/dustin/go-broadcast"
	naming "github.com/klaper_/mqtt_data_exporter/naming"
)

type TasmotaCollector struct {
	state  *prometheusTasmotaStateCollector
	sensor *prometheusTasmotaSensorCollector
}

func NewTasmotaCollector(prometheusTopicPrefix string, namingConfigurationFile string) *TasmotaCollector {
	var converter = naming.NewNamer(namingConfigurationFile)
	return &TasmotaCollector{
		state:  newPrometheusTasmotaStateCollector(prometheusTopicPrefix, converter),
		sensor: newPrometheusTasmotaSensorCollector(prometheusTopicPrefix, converter),
	}
}

func (collector *TasmotaCollector) InitializeMessageReceiver(broadcaster broadcast.Broadcaster) {
	broadcaster.Register(collector.state.channel)
	broadcaster.Register(collector.sensor.channel)
	go collector.state.collector()
	go collector.sensor.collector()
}
