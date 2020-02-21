package devices

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
	if result.Sensors == nil || len(result.Sensors) == 0 || result.Sensors["s1"] != "sensor1" {
		t.Errorf("sensor1 => Result: value %+v != expected %q", result.Sensors, "sensor1")
	}

}

func Test_NewNamer_loadConfig(t *testing.T) {
	//when
	_ = NewProperties(inputFile)

	//then
	if configuration == nil {
		t.Errorf("configuration was not loaded")
	}
}

func Test_TranslateDevice(t *testing.T) {
	//given
	namer := NewProperties(inputFile)

	//when
	result, ok := namer.GetProperties("dev3")

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
	if result.Sensors == nil || len(result.Sensors) != 0 {
		t.Errorf("sensors => Result: value %d != expected %d", len(result.Sensors), 0)
	}
}

func Test_TranslateDevice_deviceNotFound(t *testing.T) {
	//given
	namer := NewProperties(inputFile)

	//when
	_, ok := namer.GetProperties("not existing")

	//then
	if ok {
		t.Errorf("device should NOT be found")
	}
}
