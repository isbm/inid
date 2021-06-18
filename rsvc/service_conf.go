package rsvc

import (
	"log"
)

const (
	SVC_MOUNTER = iota + 1
	SVC_SERVICE
)

type ServiceConfiguration struct {
	Info        string
	After       string
	Before      string
	Stage       uint8
	Environment map[string]string
	Serial      interface{}
	Concurrent  interface{}
	serviceName string
	kind        int
}

// SetName of the service outside of YAML parsing
func (sc *ServiceConfiguration) SetName(name string) *ServiceConfiguration {
	sc.serviceName = name
	return sc
}

// GetName of the service
func (sc *ServiceConfiguration) GetName() string {
	return sc.serviceName
}

// SetKind of the service.
func (sc *ServiceConfiguration) SetKind(kind int) *ServiceConfiguration {
	sc.kind = kind
	return sc
}

// GetKind of the service
func (sc *ServiceConfiguration) GetKind() int {
	return sc.kind
}

// GetConcurrentCommands returns an array of strings, each is a command to execute in parallel
func (sc *ServiceConfiguration) GetConcurrentCommands() []string {
	return sc.getCommands(sc.Concurrent)
}

// GetSerialCommands returns an array of strings, each is a command to execute in serial sequence
func (sc *ServiceConfiguration) GetSerialCommands() []string {
	return sc.getCommands(sc.Serial)
}

// GetConcurrentConfig returns a general map->map->interface config to run a service concurrently
func (sc *ServiceConfiguration) GetConcurrentConfig() map[string]map[string]interface{} {
	return sc.getConfig(sc.Concurrent)
}

// GetSerialConfig returns a general map->map->interface config to run a service in a sequence
func (sc *ServiceConfiguration) GetSerialConfig() map[string]map[string]interface{} {
	return sc.getConfig(sc.Serial)
}

func (sc *ServiceConfiguration) getCommands(data interface{}) []string {
	if data != nil {
		switch data := data.(type) {
		case []interface{}:
			buff := []string{}
			for _, item := range data {
				if tItem, ok := item.(string); ok {
					buff = append(buff, tItem)
				}
			}
			return buff
		}
	}
	return []string{}
}

func (sc *ServiceConfiguration) getConfig(data interface{}) map[string]map[string]interface{} {
	buff := make(map[string]map[string]interface{})
	if data != nil {
		switch data := data.(type) {
		case map[interface{}]interface{}:
			for device, conf := range data {
				switch device := device.(type) {
				case string:
					switch conf := conf.(type) {
					case map[interface{}]interface{}:
						devConf := make(map[string]interface{})
						for kConf, vConf := range conf {
							if kdata, ok := kConf.(string); ok {
								devConf[kdata] = vConf
							} else {
								log.Printf("Unsupported type for key %v in configuration for device %s", kConf, device)
							}
						}
						buff[device] = devConf
					default:
						log.Printf("Unsupported configuration for device %s", device)
					}
				default:
					log.Printf("Unsupported type in mounter configuration for device %s", device)
				}
			}
		}
	}
	return buff
}
