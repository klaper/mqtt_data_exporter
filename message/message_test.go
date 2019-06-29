package message

import (
	"testing"
)

type mqttMessage struct {
	duplicate bool
	qos       byte
	retained  bool
	topic     string
	messageID uint16
	payload   []byte
	ack       func()
}

func (m mqttMessage) Duplicate() bool {
	return m.duplicate
}

func (m mqttMessage) Qos() byte {
	return m.qos
}

func (m mqttMessage) Retained() bool {
	return m.retained
}

func (m mqttMessage) Topic() string {
	return m.topic
}

func (m mqttMessage) MessageID() uint16 {
	return m.messageID
}

func (m mqttMessage) Payload() []byte {
	return m.payload
}

func (m mqttMessage) Ack() {}

func Test_GetDeviceName(t *testing.T) {
	//given
	input := []string{"cmd/device/SENSOR", "cmd/device_2/SENSOR"}
	expected := []string{"device", "device_2"}

	//when
	for i := range input {
		message := NewExporterMessage(&mqttMessage{topic: input[i]})
		result := message.GetDeviceName()
		if result != expected[i] {
			t.Errorf("DeviceName => For: %q expected: %q, but got %q", input[i], expected[i], result)
		}
	}
}
