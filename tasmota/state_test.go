package tasmota

import (
	"testing"

	"gopkg.in/yaml.v3"
)

var durationData = map[string]float64{
	"0T00:00:00": 0,
	"0T00:00:01": 1,
	"0T00:00:10": 10,
	"0T00:01:00": 60,
	"0T00:10:00": 600,
	"0T01:00:00": 3600,
	"0T10:00:00": 36000,
}

func Test_parseDuration(t *testing.T) {
	for input := range durationData {
		//when
		result := parseDuration(input)

		//then
		if result.Seconds() != durationData[input] {
			t.Errorf("\"%s\" != \"%f\" but %f", input, durationData[input], result.Seconds())
		}
	}
}

func Benchmark_parseDuration(b *testing.B) {
	// Benchmark_parseDuration-8   	 2000000	       954 ns/op
	for i := 0; i < b.N; i++ {
		parseDuration("4T12:34:56")
	}
}

var powerData = map[string]float64{
	"ON":    1,
	"OFF":   0,
	"OTHER": 0,
}

func Test_parsePower(t *testing.T) {
	for input := range powerData {
		//when
		result := parsePower(input)

		//then
		if result != powerData[input] {
			t.Errorf("\"%s\" != \"%f\" but %f", input, powerData[input], result)
		}
	}
}

var fullState = []byte("{\"Time\":\"2019-06-25T11:04:34\",\"Uptime\":\"41T12:28:52\",\"Vcc\":3.480,\"SleepMode\":\"Dynamic\",\"Sleep\":250,\"LoadAvg\":3,\"POWER\":\"ON\",\"Wifi\":{\"AP\":1,\"SSId\":\"example_ssid\",\"BSSId\":\"01:02:03:04:05:06\",\"Channel\":6,\"RSSI\":80}}")
var wifiState = []byte("{\"AP\":2,\"SSId\":\"example_ssid2\",\"BSSId\":\"06:05:04:03:02:01\",\"Channel\":2,\"RSSI\":52}")

func Test_unmarshal_loadavg(t *testing.T) {
	//given
	expected := 3
	result := state{}

	//when
	yaml.Unmarshal(fullState, &result)

	//then
	if result.Loadavg != expected {
		t.Errorf("expected: %q, got: %q", expected, result.Loadavg)
	}
}

func Test_unmarshal_power(t *testing.T) {
	//given
	expected := 1.00
	result := state{}

	//when
	yaml.Unmarshal(fullState, &result)

	//then
	if result.Power != expected {
		t.Errorf("expected: %f, got: %f", expected, result.Power)
	}
}

func Test_unmarshal_vcc(t *testing.T) {
	//given
	expected := 3.48
	result := state{}

	//when
	yaml.Unmarshal(fullState, &result)

	//then
	if result.Vcc != expected {
		t.Errorf("expected: %f, result: %f", expected, result.Vcc)
	}
}

func Test_unmarshal_wifi(t *testing.T) {
	//given
	expected := Wifi{Ap: 1, Ssid: "example_ssid", Bssid: "01:02:03:04:05:06",Channel: 6, Rssi: 80}
	result := state{}

	//when
	yaml.Unmarshal(fullState, &result)

	//then
	if result.Wifi != expected {
		t.Errorf("expected: %q, got: %q", expected, result.Wifi)
	}
}

func Test_unmarshal_wifi_ap(t *testing.T) {
	//given
	expected := 2
	result := Wifi{}

	//when
	yaml.Unmarshal(wifiState, &result)

	//then
	if result.Ap != expected {
		t.Errorf("expected: %q, got: %q", expected, result.Ap)
	}
}

func Test_unmarshal_wifi_ssid(t *testing.T) {
	//given
	expected := "example_ssid2"
	result := Wifi{}

	//when
	yaml.Unmarshal(wifiState, &result)

	//then
	if result.Ssid != expected {
		t.Errorf("expected: %q, got: %q", expected, result.Ssid)
	}
}

func Test_unmarshal_wifi_channel(t *testing.T) {
	//given
	expected := 2
	result := Wifi{}

	//when
	yaml.Unmarshal(wifiState, &result)

	//then
	if result.Channel != expected {
		t.Errorf("expected: %q, got: %q", expected, result.Channel)
	}
}

func Test_unmarshal_wifi_rssi(t *testing.T) {
	//given
	expected := 52
	result := Wifi{}

	//when
	yaml.Unmarshal(wifiState, &result)

	//then
	if result.Rssi != expected {
		t.Errorf("expected: %q, got: %q", expected, result.Rssi)
	}
}
func Test_isTasmotaStateMessage(t *testing.T) {
	//given
	input := []string{"cmd/device/STATE", "cmd/device_2/SENSOR"}
	expected := []bool{true, false}

	//when
	for i := range input {
		result := isStateMessage(input[i])
		if result != expected[i] {
			t.Errorf("istasmotaStateMessage => For: %q expected: %t, but got %t", input[i], expected[i], result)
		}
	}
}
