package main

import (
	exporterMessage "github.com/klaper_/mqtt_data_exporter/message"
	"github.com/klaper_/mqtt_data_exporter/naming"
	"github.com/klaper_/mqtt_data_exporter/prom"
	tasmota "github.com/klaper_/mqtt_data_exporter/tasmota"
	"github.com/prometheus/client_golang/prometheus"

	"log"
	"net/http"

	broadcast "github.com/dustin/go-broadcast"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	opsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_exporter_total_message_count",
			Help: "Count of MQTT messages processed",
		}, []string{"source_instance"})
)

func init() {
	prometheus.MustRegister(opsProcessed)
}

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
		log.Printf("Connected to MQTT")
		if token := c.Subscribe("tele/#", byte(1), onMessageReceived); token.Wait() && token.Error() != nil {
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
	msg := exporterMessage.NewExporterMessage(message)
	log.Printf("Received message on topic: %s", message.Topic())
	opsProcessed.WithLabelValues(msg.GetDeviceName()).Inc()
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
			"prefix for metrics names",
		).Default("mqtt_exporter").String()
	)

	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	var metricsStore = prom.NewMetrics(*metricsPrefix, naming.NewNamer(*namingFile))
	var tasmotaCollector = tasmota.NewTasmotaCollector(metricsStore)
	tasmotaCollector.InitializeMessageReceiver(broadcaster)

	mqttInit(mqttHost, mqttClientId, mqttUsername, mqttPassword)
	prometheusListenAndServer(listenAddress, metricsPath)
}
