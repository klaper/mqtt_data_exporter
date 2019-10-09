package naming

import (
	"testing"
)

var inputFile = "./testdata/sampleConfiguration.yaml"

func Test_loadConfiguration_groups(t *testing.T) {
	//when
	result := loadConfiguration(inputFile)

	//then
	_, ok := result["dev1"]

	if !ok {
		t.Errorf("group search => Result: %q not found in: %+v", "dev1", result)
	}
}

func Test_loadConfiguration_devices(t *testing.T) {
	//given
	expectedCount := 4

	//when
	result := loadConfiguration(inputFile)

	//then
	if len(result) != expectedCount {
		t.Errorf("device count => Result: value %d != expected %d", len(result), expectedCount)
	}
}

func Test_loadConfiguration_correctResult(t *testing.T) {
	//when
	configuration := loadConfiguration(inputFile)

	//then
	result, _ := configuration["dev1"]

	if result.Group != "livingroom" {
		t.Errorf("group => Result: value %q != expected %q", result.Group, "livingroom")
	}

	if result.Device != "dev1" {
		t.Errorf("device => Result: value %q != expected %q", result.Device, "dev1")
	}
	if result.Name != "d1" {
		t.Errorf("name => Result: value %q != expected %q", result.Name, "d1")
	}
}

func Test_NewNamer_loadConfig(t *testing.T) {
	//when
	_ = NewNamer(inputFile)

	//then
	if configuration == nil {
		t.Errorf("configuration was not loaded")
	}
}

func Test_TranslateDevice(t *testing.T) {
	//given
	namer := NewNamer(inputFile)

	//when
	result, ok := namer.TranslateDevice("dev3")

	//then
	if !ok {
		t.Errorf("device should be found")
	}
	if result.Group != "someroom1" {
		t.Errorf("group => Result: value %q != expected %q", result.Group, "someroom1")
	}

	if result.Device != "dev3" {
		t.Errorf("device => Result: value %q != expected %q", result.Device, "dev3")
	}
	if result.Name != "d3" {
		t.Errorf("name => Result: value %q != expected %q", result.Name, "d3")
	}
}

func Test_TranslateDevice_deviceNotFound(t *testing.T) {
	//given
	namer := NewNamer(inputFile)

	//when
	_, ok := namer.TranslateDevice("not existing")

	//then
	if ok {
		t.Errorf("device should NOT be found")
	}
}
