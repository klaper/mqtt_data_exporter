package tasmota

import (
	"github.com/klaper_/mqtt_data_exporter/prom"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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
	metricsStore *prom.Metrics
	channel      chan interface{}
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

func newPrometheusTasmotaStateCollector(metricsStore *prom.Metrics) (collector *prometheusTasmotaStateCollector) {
	metricsStore.RegisterMetric(
		prom.GAUGE,
		"upTimeGauge",
		"tasmota_state_uptime",
		"Uptime of tasmota entity",
		[]string{},
	)
	metricsStore.RegisterMetric(
		prom.GAUGE,
		"rssiGauge",
		"tasmota_state_rssi",
		"Signal strength of tasmota entity",
		[]string{"ssid", "channel", "ap_index"},
	)
	metricsStore.RegisterMetric(
		prom.GAUGE,
		"powerGauge",
		"tasmota_power",
		"Power state of tasmota entity",
		[]string{},
	)
	return &prometheusTasmotaStateCollector{
		metricsStore: metricsStore,
		channel:      make(chan interface{}),
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
		err = yaml.Unmarshal((message).Payload(), &state)
		if err != nil {
			fatal("state", "error while unmarshaling", err)
			continue
		}
		collector.metricsStore.Set("upTimeGauge", message.GetDeviceName(), map[string]string{}, state.Uptime.Seconds())
		collector.metricsStore.Set("powerGauge", message.GetDeviceName(), map[string]string{}, state.Power)

		collector.metricsStore.Set(
			"rssiGauge",
			message.GetDeviceName(),
			map[string]string{
				"ssid":     state.Wifi.Ssid,
				"channel":  strconv.Itoa(state.Wifi.Channel),
				"ap_index": strconv.Itoa(state.Wifi.Ap),
			},
			float64(state.Wifi.Rssi),
		)
	}
}
