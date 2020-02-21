# MQTT Data Exporter ![travis ci](https://travis-ci.org/klaper/mqtt_data_exporter.svg?branch=master)

### What does it do?

Simply it subscribes to MQTT topics and exposes data published there in way that is easy to read with Prometheus

### Links

- [Github](https://github.com/klaper/mqtt_data_exporter/)
- [Docker Hub](https://hub.docker.com/r/klaper/mqtt_data_exporter)

### Accepted parameters

```
web.listen-address:     [Default: ":2112"]                          Address on which to expose metrics and web interface. 
web.telemetry-path:     [Default: "/metrics"]                       Path under which to expose metrics.
mqtt.host:              [Default: "127.0.0.1:1883"]                 Mqtt host address and port.
mqtt.clientId:          [Default: "mqtt_exporter"]                  Mqtt clientId.
mqtt.username:          [Default: ""]                               Mqtt username.
mqtt.password:          [Default: ""]                               Mqtt password
naming.config:          [Default: "/etc/mqtt_exporter/naming.yaml"] File containg naming convertions
metrics.prefix:         [Default: "mqtt_exporter"]                  Prefix for metrics names
log.level:              [Default: 2]                                Log level
```

#### naming conversion file format:
```
Group1:                             # "group" attribute of metric
  - device: device_name             # device name from mqtt message
                                    # used for matching
                                    # "device" attribute of metrics
    name: friendly_name             # "friendly_name" attribute value
    sensors:				        # list of "sensor_name"
        sensor_name: sensor_alias   # for this device for every readout from "sensor_name" there will be "sensor_alias" label added
```

#### log levels parameter values
| value | meaning |
|-------|---------|
|   1   |  DEBUG  | 
|   2   |  INFO   |
|   3   |  WARN   |
|   4   |  ERROR  | 
|   5   |  OFF    |

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