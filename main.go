package main

import (
	"github.com/klaper_/mqtt_data_exporter/devices"
	"github.com/klaper_/mqtt_data_exporter/logger"
	exporterMessage "github.com/klaper_/mqtt_data_exporter/message"
	"github.com/klaper_/mqtt_data_exporter/prom"
	"github.com/klaper_/mqtt_data_exporter/tasmota"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dustin/go-broadcast"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

const moduleId = "main"

var (
	metricsStore *prom.Metrics
)

func prometheusListenAndServer(listenAddress *string, metricsPath *string) {
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", *metricsPath)
		w.WriteHeader(http.StatusMovedPermanently)
	})
	log.Panic(http.ListenAndServe(*listenAddress, nil))
}

func mqttInit(mqttHost *string, mqttClientId *string, mqttUser *string, mqttPassword *string) {
	connOpts := MQTT.
		NewClientOptions().
		AddBroker(*mqttHost).
		SetClientID(*mqttClientId).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetUsername(*mqttUser).
		SetPassword(*mqttPassword)

	connOpts.OnConnect = func(c MQTT.Client) {
		logger.Debug(moduleId, "Connected to MQTT")
		if token := c.Subscribe("#", byte(1), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

var broadcaster = broadcast.NewBroadcaster(100)

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	msg := exporterMessage.NewExporterMessage(message, metricsStore)
	logger.Debug(moduleId, "Received message on topic: %s", message.Topic())
	metricsStore.CounterInc(
		"total_message_count",
		msg.GetDeviceName(),
		map[string]string{},
	)
	broadcaster.Submit(msg)
}

func main() {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(":2112").String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		mqttHost = kingpin.Flag(
			"mqtt.host",
			"Mqtt host address and port.",
		).Default("127.0.0.1:1883").String()
		mqttClientId = kingpin.Flag(
			"mqtt.clientId",
			"Mqtt clientId",
		).Default("mqtt_exporter").String()
		mqttUsername = kingpin.Flag(
			"mqtt.username",
			"Mqtt username",
		).Default().String()
		mqttPassword = kingpin.Flag(
			"mqtt.password",
			"Mqtt password",
		).Default().String()
		namingFile = kingpin.Flag(
			"naming.config",
			"File containg naming convertions",
		).Default("/etc/mqtt_exporter/naming.yaml").String()
		metricsPrefix = kingpin.Flag(
			"metrics.prefix",
			"Prefix for metrics names",
		).Default("mqtt_exporter").String()
		logLevel = kingpin.Flag(
			"log.level",
			"DEBUG = 1; INFO = 2; WARN = 3; ERROR = 4; OFF = 5",
		).Default("2").String()
		metricsCleanerTimeout = kingpin.Flag(
			"cleaner.gauge.timeout",
			"Timeout for gauge value cleaner (0 = disabled)",
		).Default("0s").Duration()
	)

	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	toSet, err := strconv.Atoi(*logLevel)
	if err != nil {
		toSet = 2
	}
	logger.SetLogLevel(logger.Loglevel(toSet))

	prepareMetricsStore(metricsPrefix, namingFile, metricsCleanerTimeout)

	var tasmotaCollector = tasmota.NewTasmotaCollector(metricsStore)
	tasmotaCollector.InitializeMessageReceiver(broadcaster)

	mqttInit(mqttHost, mqttClientId, mqttUsername, mqttPassword)
	prometheusListenAndServer(listenAddress, metricsPath)
}

func prepareMetricsStore(metricsPrefix *string, namingConfiguration *string, metricsCleanerTimeout *time.Duration) {
	metricsStore = prom.NewMetrics(*metricsPrefix, devices.NewProperties(*namingConfiguration), *metricsCleanerTimeout)
	metricsStore.RegisterCounter(
		"total_message_count",
		"total_message_count",
		"Count of MQTT messages processed",
		[]string{},
	)
	metricsStore.RegisterCounter(
		"message_count",
		"message_count",
		"Count of MQTT messages processed",
		[]string{"processing_state", "exporter_module"},
	)
}
