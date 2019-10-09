package naming

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type NamerDevice struct {
	Device string
	Name   string
	Group  string
}
type NamerConfiguration map[string]NamerDevice

var configuration = NamerConfiguration(nil)

func loadConfiguration(file string) NamerConfiguration {
	type YamlNamerDevice struct {
		Device string `yaml:"device"`
		Name   string `yaml:"name"`
	}

	type YamlNamerConfiguration map[string][]YamlNamerDevice

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	configuration := YamlNamerConfiguration{}
	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		panic(err)
	}

	result := NamerConfiguration{}

	for group := range configuration {
		devices := configuration[group]
		for i := range devices {
			device := devices[i]

			result[device.Device] = NamerDevice{device.Device, device.Name, group}
		}
	}

	return result
}

type Namer struct{}

func NewNamer(configFile string) *Namer {
	if configuration == nil {
		configuration = loadConfiguration(configFile)
	}
	return &Namer{}
}

func (*Namer) TranslateDevice(deviceName string) (*NamerDevice, bool) {
	device, ok := configuration[deviceName]
	if !ok {
		return nil, false
	}
	return &NamerDevice{device.Device, device.Name, device.Group}, true
}
