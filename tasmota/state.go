package tasmota

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
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
	Wifi    tasmotaWifi
}

type prometheusTasmotaStateCollector struct {
	upTimeGauge *prometheus.GaugeVec
	rssiGauge   *prometheus.GaugeVec

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

func (state *tasmotaState) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		Uptime  string      `yaml:"Uptime"`
		Loadavg int         `yaml:"LoadAvg"`
		Vcc     float64     `yaml:"Vcc"`
		Wifi    tasmotaWifi `yaml:"Wifi"`
	}
	var tmp alias
	err := unmarshal(&tmp)
	if err != nil {
		return err
	}

	state.Uptime = parseDuration(tmp.Uptime)
	state.Loadavg = tmp.Loadavg
	state.Wifi = tmp.Wifi
	state.Vcc = tmp.Vcc

	return nil
}

func newPrometheusTasmotaStateCollector(metricsPrefix string) (collector *prometheusTasmotaStateCollector) {
	upTimeGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricsPrefix + "_tasmota_state_uptime",
			Help: "Uptime of tasmota entity",
		},
		[]string{"tasmota_instance"},
	)
	rssiGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricsPrefix + "_tasmota_state_rssi",
			Help: "Signal strength of tasmota entity",
		},
		[]string{"tasmota_instance", "ssid", "channel", "ap_index"},
	)

	prometheus.MustRegister(upTimeGauge)
	prometheus.MustRegister(rssiGauge)

	return &prometheusTasmotaStateCollector{
		upTimeGauge: upTimeGauge,
		rssiGauge:   rssiGauge,
		channel:     make(chan interface{}),
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
		device := message.GetDeviceName()

		collector.upTimeGauge.WithLabelValues(device).Set(state.Uptime.Seconds())

		collector.rssiGauge.WithLabelValues(
			device,
			state.Wifi.Ssid,
			strconv.Itoa(state.Wifi.Channel),
			strconv.Itoa(state.Wifi.Ap),
		).Set(float64(state.Wifi.Rssi))
	}
}
