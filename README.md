# MQTT Data Exporter ![travis ci](https://travis-ci.org/klaper/mqtt_data_exporter.svg?branch=master)

### What does it do?

Simply it subscribes to MQTT topics and exposes data published there in way that is easy to read with Prometheus

### Links

- [Github](https://github.com/klaper/mqtt_data_exporter/)
- [Docker Hub](https://hub.docker.com/r/klaper/mqtt_data_exporter)

### Accepted parameters

```
web.listen-address: [Default: ":2112"] 				Address on which to expose metrics and web interface. 
web.telemetry-path: [Default: "/metrics"] 			Path under which to expose metrics.
mqtt.host:          [Default: "127.0.0.1:1883"] 	Mqtt host address and port.
```

### Running in docker

First build project:
```
cd $GOPATH/src/github.com/klaper_/mqtt_data_exporter
go build -o bin/mqtt_data_exporter
```
Build docker image:
```
docker build -t mqtt_data_exporter .
```
Run container:
```
docker run -p 2112:2112 mqtt_data_exporter --mqtt.host "mqtthost.local:1883"
```