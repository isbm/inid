package rsvc

import (
	"github.com/isbm/processman"
)

type InidService interface {
	Init(descrPath string) error
	SetEnviron(env map[string]string) *InidService
	GetConfiguration() *ServiceConfiguration
	GetProcesses() map[int]*processman.Process
	Start() error
	Kill() error
	Stop() error
	Restart() error
}
