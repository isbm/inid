package rsvc

import (
	"fmt"

	"github.com/isbm/processman"
)

type InidService interface {
	Init(descrPath string) error
	SetEnviron(env map[string]string) InidService
	GetServiceConfiguration() *ServiceConfiguration
	GetProcesses() map[int]*processman.Process
	Start() error
	Kill() error
	Stop() error
	Restart() error
}

type BaseService struct {
	env []string
}

func (svc *BaseService) SetEnviron(env map[string]string) InidService {
	svc.env = make([]string, 0)
	for k := range env {
		svc.env = append(svc.env, fmt.Sprintf("%s=%s", k, env[k]))
	}
	return svc
}

func (svc *BaseService) Init(descrPath string) error                    { return nil }
func (svc *BaseService) GetServiceConfiguration() *ServiceConfiguration { return nil }
func (svc *BaseService) GetProcesses() map[int]*processman.Process      { return nil }
func (svc *BaseService) Start() error                                   { return nil }
func (svc *BaseService) Kill() error                                    { return nil }
func (svc *BaseService) Stop() error                                    { return nil }
func (svc *BaseService) Restart() error                                 { return nil }
