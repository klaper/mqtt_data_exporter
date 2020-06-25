package tasmota

import (
	"github.com/klaper_/mqtt_data_exporter/logger"
	exporterMessage "github.com/klaper_/mqtt_data_exporter/message"
)

type NotExporterMessage struct {
	message string
}
type TopicValidatedToFalse struct {
	message string
}

func (err NotExporterMessage) Error() string {
	return err.message
}

func (err TopicValidatedToFalse) Error() string {
	return err.message
}

func receiveMessage(tmp interface{}, module string, topicValidator func(string) bool) (*exporterMessage.ExporterMessage, error) {
	message, ok := tmp.(*exporterMessage.ExporterMessage)
	logger.Debug(module, "message: %+v, ok: %t", message, ok)
	if !ok {
		logger.Info(module, "Message was not an ExporterMessage")
		return nil, NotExporterMessage{message: "Message was not an ExporterMessage"}
	}
	logger.Debug(module, "Message(%d).Topic %q", (message).MessageID(), (message).Topic())
	if !topicValidator((message).Topic()) {
		logger.Debug(module, "DEBUG: Message(%d) was skipped due to wrong topic %s", (message).MessageID(), (message).Topic())
		message.ProcessMessage(module, exporterMessage.Ignored)
		return nil, TopicValidatedToFalse{message: "Skipped due to wrong topic"}
	}
	message.ProcessMessage(module, exporterMessage.Processed)
	logger.Debug(module, "Message(%d) was processed", (message).MessageID())
	return message, nil
}
