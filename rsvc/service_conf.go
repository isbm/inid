package rsvc

type ServiceConfiguration struct {
	Info        string
	After       string
	Before      string
	Stage       uint8
	Environment map[string]string
	Serial      []string
	Concurrent  []string
}
