package tasmota

import (
	broadcast "github.com/dustin/go-broadcast"
)

type TasmotaCollector struct {
	state *prometheusTasmotaStateCollector
}

func NewTasmotaCollector(prometheusTopicPrefix string) *TasmotaCollector {
	return &TasmotaCollector{
		state: newPrometheusTasmotaStateCollector(prometheusTopicPrefix),
	}
}

func (collector *TasmotaCollector) InitializeMessageReceiver(broadcaster broadcast.Broadcaster) {
	broadcaster.Register(collector.state.channel)
	go collector.state.collector()
}
