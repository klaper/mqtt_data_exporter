package tasmota

import (
	"github.com/klaper_/mqtt_data_exporter/prom"
	"testing"

	exporterMessage "github.com/klaper_/mqtt_data_exporter/message"
)

type messageMock struct {
	topic      string
	payload    []byte
	deviceName string
}

func (e messageMock) Duplicate() bool {
	return false
}

func (e messageMock) Qos() byte {
	return byte(1)
}

func (e messageMock) Retained() bool {
	return false
}

func (e messageMock) Topic() string {
	return e.topic
}

func (e messageMock) MessageID() uint16 {
	return 0
}

func (e messageMock) Payload() []byte {
	return e.payload
}

func (e messageMock) Ack() {}

func Test_receiveMessageNonMessage(t *testing.T) {
	//given
	inputTmp := "some"
	inputModule := "module"

	//when
	result, err := receiveMessage(inputTmp, inputModule, func(string) bool { return true })

	//then
	if _, ok := err.(NotExporterMessage); !ok || result != nil {
		t.Errorf("NonExporterMessage => Error was: %q (expecting: NotExporterMessage), and Result was: %+v (expected: nil)", err, result)
	}
}

func Test_receiveMessageDifferentTopic(t *testing.T) {
	//given
	inputTmp := exporterMessage.NewExporterMessage(
		messageMock{topic: "some/topic/value"},
		prom.NewMetrics("", nil),
	)
	inputModule := "module"

	//when
	result, err := receiveMessage(inputTmp, inputModule, func(string) bool { return false })

	//then
	if _, ok := err.(TopicValidatedToFalse); !ok || result != nil {
		t.Errorf("TopicValidatedToFalse => Error was: %q (expecting: TopicValidatedToFalse), and Result was: %+v (expected: nil)", err, result)
	}
}

func Test_receiveMessageSameTopic(t *testing.T) {
	//given
	inputTmp := exporterMessage.NewExporterMessage(
		messageMock{topic: "some/topic/value"},
		prom.NewMetrics("", nil),
	)
	inputModule := "module"

	//when
	result, err := receiveMessage(inputTmp, inputModule, func(string) bool { return true })

	//then
	if result == nil || err != nil {
		t.Errorf("CorrectMessage => Error was: %+v (expecting: nil), and Result was: %+v (expected: not nil)", err, result)
	}
}
