package message

import (
	"github.com/klaper_/mqtt_data_exporter/prom"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MessageState string

const (
	MessageProcessed MessageState = "processed"
	MessageIgnored   MessageState = "ignored"
)

type ExporterMessage struct {
	msg          MQTT.Message
	metricsStore *prom.Metrics
}

func NewExporterMessage(msg MQTT.Message, metricsStore *prom.Metrics) *ExporterMessage {
	return &ExporterMessage{msg: msg, metricsStore: metricsStore}
}

func (e *ExporterMessage) GetDeviceName() string {
	return strings.Split(e.msg.Topic(), "/")[1]
}

func (e *ExporterMessage) ProcessMessage(exporterModule string, state MessageState) {
	e.metricsStore.CounterInc(
		"message_count",
		e.GetDeviceName(),
		map[string]string{
			"processing_state": string(state),
			"exporter_module":  exporterModule,
		},
	)
}

func (e *ExporterMessage) Duplicate() bool {
	return e.msg.Duplicate()
}

func (e *ExporterMessage) Qos() byte {
	return e.msg.Qos()
}

func (e *ExporterMessage) Retained() bool {
	return e.msg.Retained()
}

func (e *ExporterMessage) Topic() string {
	return e.msg.Topic()
}

func (e *ExporterMessage) MessageID() uint16 {
	return e.msg.MessageID()
}

func (e *ExporterMessage) Payload() []byte {
	return e.msg.Payload()
}

func (e *ExporterMessage) Ack() {
	e.msg.Ack()
}
