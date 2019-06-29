package tasmota

import (
	exporterMessage "github.com/klaper_/mqtt_data_exporter/message"
)

type NotExporterMessage struct {
	message string
}
type NotExpectedMessage struct {
	message string
}

func (err NotExporterMessage) Error() string {
	return err.message
}

func (err NotExpectedMessage) Error() string {
	return err.message
}

func receiveMessage(tmp interface{}, module string, topicValidator func(string) bool) (*exporterMessage.ExporterMessage, error) {
	message, ok := tmp.(*exporterMessage.ExporterMessage)
	debug(module, "message: %+v, ok: %t", message, ok)
	if !ok {
		info(module, "Message was not an ExporterMessage")
		return nil, NotExporterMessage{message: "Message was not an ExporterMessage"}
	}
	debug(module, "Message(%d).Topic %q", (message).MessageID(), (message).Topic())
	if !isTasmotaStateMessage((message).Topic()) {
		debug("state", "DEBUG: Message(%d) was skipped due to wrong topic", (message).MessageID())
		message.ProcessMessage(module, exporterMessage.MessageIgnored)
		return nil, NotExpectedMessage{message: "Skipped due to wrong topic"}
	}
	message.ProcessMessage(module, exporterMessage.MessageProcessed)
	debug(module, "Message(%d) was processed", (message).MessageID())
	return message, nil
}
