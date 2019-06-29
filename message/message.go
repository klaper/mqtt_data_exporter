package message

import (
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
)

type MessageState string

const (
	MessageProcessed MessageState = "processed"
	MessageIgnored   MessageState = "ignored"
)

const prometheusTopicPrefix = "mqtt_exporter"

var processingGauge = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: prometheusTopicPrefix + "_message_count",
		Help: "Count of MQTT messages processed",
	}, []string{"source_instance", "processing_state", "exporter_module"})

func init() {
	prometheus.MustRegister(processingGauge)
}

type ExporterMessage struct {
	msg MQTT.Message
}

func NewExporterMessage(msg MQTT.Message) *ExporterMessage {
	return &ExporterMessage{msg: msg}
}

func (e *ExporterMessage) GetDeviceName() string {
	return strings.Split(e.msg.Topic(), "/")[1]
}

func (e *ExporterMessage) ProcessMessage(exporterModule string, state MessageState) {
	processingGauge.WithLabelValues(e.GetDeviceName(), string(state), exporterModule).Inc()
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
