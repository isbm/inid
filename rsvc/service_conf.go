package rsvc

type ServiceConfiguration struct {
	Info        string
	After       string
	Before      string
	Environment map[string]string
	Serial      []string
	Concurrent  []string
}
