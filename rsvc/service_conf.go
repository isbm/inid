package rsvc

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
	Serial      []string
	Concurrent  []string
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
