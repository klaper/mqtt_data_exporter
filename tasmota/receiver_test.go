package tasmota

import (
	"testing"
)

func Test_receiveMessageNonMessage(t *testing.T) {
	//given
	inputTmp := "some"
	inputModule := "module"

	//when
	result, err := receiveMessage(inputTmp, inputModule, func(string) bool { return true })

	//then
	if _, ok := err.(NotExporterMessage); !ok || result != nil {
		t.Errorf("NonExporterMessage => error was: %q, and result was: %q", err, result)
	}
}
