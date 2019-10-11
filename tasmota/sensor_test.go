package tasmota

import (
	"testing"
)

/*
"{\"Time\":\"2019-06-25T21:29:37\",\"Switch1\":\"OFF\",\"ANALOG\":{\"A0\":3},\"BMP280\":{\"temperature\":30.3,\"pressure\":1010.7},\"BH1750\":{\"illuminance\":0},\"PressureUnit\":\"hPa\",\"TempUnit\":\"C\"}"
*/

func Test_getKeys(t *testing.T) {
	//given
	inputData := map[string]interface{}{
		"TempUnit":         "1",
		"SomeExampleFiled": "2",
		"HumidityUnit":     "3",
		"PressureUnit":     "4",
		"OtheStuff":        "5",
	}

	expectedKeys := map[string]bool{"SomeExampleFiled": true, "OtheStuff": true}
	expectedUnits := map[string]bool{"TempUnit": true, "HumidityUnit": true, "PressureUnit": true}

	//when
	resultKeys, resultUnits := getKeys(inputData)

	//then
	if len(resultKeys) != len(expectedKeys) {
		t.Errorf("Size(keys) => Expected: %d, got: %d", len(expectedKeys), len(resultKeys))
	}
	for i := range resultKeys {
		_, ok := expectedKeys[resultKeys[i]]
		if !ok {
			t.Errorf("Key => Result: %q not found in: %+v", resultKeys[i], expectedKeys)
		}
	}
	if len(resultUnits) != len(expectedUnits) {
		t.Errorf("Size(units) => Expected: %d, got: %d", len(expectedUnits), len(resultUnits))
	}
	for i := range resultKeys {
		_, ok := expectedUnits[resultUnits[i]]
		if !ok {
			t.Errorf("Units => Result: %q not found in: %+v", resultUnits[i], expectedUnits)
		}
	}
}

func Test_getSensorData(t *testing.T) {
	//given
	var inputData = map[string]interface{}{
		"SI7021": map[interface{}]interface{}{
			"Temperature": 123.23,
		},
	}

	//when
	result := getSensorData(inputData)

	//then
	if len(result) != 1 {
		t.Errorf("Size => Expected: %d, got: %d", 1, len(result))
	}

	if result[0].Type != temperature {
		t.Errorf("Type => Expected: %+v, got: %+v", temperature, result[0].Type)
	}

	if result[0].Value != 123.23 {
		t.Errorf("Value => Expected: %f, got: %+v", 123.23, result[0].Value)
	}

	if result[0].SensorName != "SI7021" {
		t.Errorf("Value => Expected: %q, got: %q", "SI7021", result[0].SensorName)
	}
}

func Benchmark_getSensorData(b *testing.B) {
	//given
	var inputData = map[string]interface{}{
		"SI7021": map[interface{}]interface{}{
			"Temperature": 123.23,
		},
	}

	//when
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getSensorData(inputData)
	}

}

func Test_getSensorData_withTwoSensors(t *testing.T) {
	//given
	var inputData = map[string]interface{}{
		"BMP280": map[interface{}]interface{}{
			"Temperature": 124.23,
			"Pressure":    1024,
		},
	}

	//when
	result := getSensorData(inputData)

	//then
	if len(result) != 2 {
		t.Errorf("Size => Expected: %d, got: %d", 2, len(result))
		return
	}

	if result[0].Type != temperature {
		t.Errorf("Type => Expected: %+v, got: %+v", temperature, result[0].Type)
	}

	if result[0].Value != 124.23 {
		t.Errorf("Value => Expected: %f, got: %+v", 124.23, result[0].Value)
	}

	if result[1].Type != pressure {
		t.Errorf("Type => Expected: %+v, got: %+v", pressure, result[1].Type)
	}

	if result[1].Value != 1024 {
		t.Errorf("Value => Expected: %d, got: %+v", 1024, result[1].Value)
	}
}
