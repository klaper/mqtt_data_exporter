package tasmota

import (
	broadcast "github.com/dustin/go-broadcast"
)

type TasmotaCollector struct {
	state  *prometheusTasmotaStateCollector
	sensor *prometheusTasmotaSensorCollector
}

func NewTasmotaCollector(prometheusTopicPrefix string) *TasmotaCollector {
	return &TasmotaCollector{
		state:  newPrometheusTasmotaStateCollector(prometheusTopicPrefix),
		sensor: newPrometheusTasmotaSensorCollector(prometheusTopicPrefix),
	}
}

func (collector *TasmotaCollector) InitializeMessageReceiver(broadcaster broadcast.Broadcaster) {
	broadcaster.Register(collector.state.channel)
	broadcaster.Register(collector.sensor.channel)
	go collector.state.collector()
	go collector.sensor.collector()
}
