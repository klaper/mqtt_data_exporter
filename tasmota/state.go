package tasmota

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v3"

	naming "github.com/klaper_/mqtt_data_exporter/naming"
)

const stateClientId = "tasmota_state"

var durationRegex = regexp.MustCompile(`(?P<days>\d+)?T?(?P<hours>\d+)?:(?P<minutes>\d+)?:(?P<seconds>\d+)?`)

type tasmotaWifi struct {
	Ap      int    `yaml:"AP"`
	Ssid    string `yaml:"SSId"`
	Channel int    `yaml:"Channel"`
	Rssi    int    `yaml:"RSSI"`
}

type tasmotaState struct {
	Uptime  time.Duration
	Vcc     float64
	Loadavg int
	Power   float64
	Wifi    tasmotaWifi
}

type prometheusTasmotaStateCollector struct {
	upTimeGauge *prometheus.GaugeVec
	rssiGauge   *prometheus.GaugeVec
	powerGauge  *prometheus.GaugeVec

	converter *naming.Namer

	channel chan interface{}
}

func parseDuration(str string) time.Duration {
	matches := durationRegex.FindStringSubmatch(str)

	days, _ := strconv.Atoi(matches[1])
	hours, _ := strconv.Atoi(matches[2])
	minutes, _ := strconv.Atoi(matches[3])
	seconds, _ := strconv.Atoi(matches[4])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(int64(days)*24*hour + int64(hours)*hour + int64(minutes)*minute + int64(seconds)*second)
}

func parsePower(str string) float64 {
	if str == "ON" {
		return 1
	}
	return 0
}

func (state *tasmotaState) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		Uptime  string      `yaml:"Uptime"`
		Loadavg int         `yaml:"LoadAvg"`
		Vcc     float64     `yaml:"Vcc"`
		Power   string      `yaml:"POWER"`
		Wifi    tasmotaWifi `yaml:"Wifi"`
	}
	var tmp alias
	err := unmarshal(&tmp)
	if err != nil {
		return err
	}
	debug("state", "Got %+v as state input", tmp)
	state.Uptime = parseDuration(tmp.Uptime)
	state.Loadavg = tmp.Loadavg
	state.Wifi = tmp.Wifi
	state.Vcc = tmp.Vcc
	state.Power = parsePower(tmp.Power)
	debug("state", "Got %+v as state output", *state)

	return nil
}

func newPrometheusTasmotaStateCollector(metricsPrefix string, namer *naming.Namer) (collector *prometheusTasmotaStateCollector) {
	upTimeGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricsPrefix + "_tasmota_state_uptime",
			Help: "Uptime of tasmota entity",
		},
		[]string{"device", "group", "friendly_name"},
	)
	rssiGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricsPrefix + "_tasmota_state_rssi",
			Help: "Signal strength of tasmota entity",
		},
		[]string{"device", "group", "friendly_name", "ssid", "channel", "ap_index"},
	)
	powerGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricsPrefix + "_tasmota_power",
			Help: "Power state of tasmota entity",
		},
		[]string{"device", "group", "friendly_name"},
	)

	prometheus.MustRegister(upTimeGauge)
	prometheus.MustRegister(rssiGauge)
	prometheus.MustRegister(powerGauge)

	return &prometheusTasmotaStateCollector{
		upTimeGauge: upTimeGauge,
		rssiGauge:   rssiGauge,
		powerGauge:  powerGauge,
		channel:     make(chan interface{}),
		converter:   namer,
	}
}

func isTasmotaStateMessage(topic string) bool {
	split := strings.Split(topic, "/")
	return len(split) == 3 && split[2] == "STATE"
}

func (collector *prometheusTasmotaStateCollector) collector() {
	for tmp := range collector.channel {
		message, err := receiveMessage(tmp, "state", isTasmotaStateMessage)
		if err != nil {
			continue
		}

		state := tasmotaState{}
		err = yaml.Unmarshal([]byte((message).Payload()), &state)
		if err != nil {
			fatal("state", "error while unmarshaling", err)
			continue
		}
		device, ok := collector.converter.TranslateDevice(message.GetDeviceName())
		if !ok {
			name := message.GetDeviceName()
			warn("state", "Device configuration %s was not found", name)
			device = &naming.NamerDevice{name, name, name}
		}

		collector.upTimeGauge.WithLabelValues(device.Device, device.Group, device.Name).Set(state.Uptime.Seconds())
		collector.powerGauge.WithLabelValues(device.Device, device.Group, device.Name).Set(state.Power)

		collector.rssiGauge.WithLabelValues(
			device.Device,
			device.Group,
			device.Name,
			state.Wifi.Ssid,
			strconv.Itoa(state.Wifi.Channel),
			strconv.Itoa(state.Wifi.Ap),
		).Set(float64(state.Wifi.Rssi))
	}
}
