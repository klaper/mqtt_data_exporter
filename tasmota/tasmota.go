package tasmota

import (
	"github.com/dustin/go-broadcast"
	"github.com/klaper_/mqtt_data_exporter/prom"
)

type Collector struct {
	state  *stateCollector
	sensor *sensorCollector
}

func NewTasmotaCollector(metricsStore *prom.Metrics) *Collector {
	return &Collector{
		state:  newStateCollector(metricsStore),
		sensor: newSensorCollector(metricsStore),
	}
}

func (collector *Collector) InitializeMessageReceiver(broadcaster broadcast.Broadcaster) {
	broadcaster.Register(collector.state.channel)
	broadcaster.Register(collector.sensor.channel)
	go collector.state.collector()
	go collector.sensor.collector()
}
