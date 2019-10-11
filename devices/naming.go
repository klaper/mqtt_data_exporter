package devices

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Properties struct {
	Device string
	Name   string
	Group  string
}
type Configuration map[string]Properties

var configuration = Configuration(nil)

func loadConfiguration(file string) Configuration {
	type YamlDevice struct {
		Device string `yaml:"device"`
		Name   string `yaml:"name"`
	}

	type YamlConfiguration map[string][]YamlDevice

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	configuration := YamlConfiguration{}
	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		panic(err)
	}

	result := Configuration{}

	for group := range configuration {
		devices := configuration[group]
		for i := range devices {
			device := devices[i]

			result[device.Device] = Properties{device.Device, device.Name, group}
		}
	}

	return result
}

type PropertiesProvider struct{}

func NewProperties(configFile string) *PropertiesProvider {
	if configuration == nil {
		configuration = loadConfiguration(configFile)
	}
	return &PropertiesProvider{}
}

func (*PropertiesProvider) GetProperties(deviceName string) (*Properties, bool) {
	device, ok := configuration[deviceName]
	if !ok {
		return nil, false
	}
	return &Properties{device.Device, device.Name, device.Group}, true
}
