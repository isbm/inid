package rsvc

type ServiceConfiguration struct {
	Info        string
	After       string
	Before      string
	Stage       uint8
	Environment map[string]string
	Serial      []string
	Concurrent  []string
	serviceName string
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
