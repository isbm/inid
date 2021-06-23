package rsvc

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/isbm/processman"
	"gopkg.in/yaml.v2"
)

type InidService interface {
	Init(descrPath string) error
	postInit() error
	SetEnviron(env map[string]string) InidService
	GetServiceConfiguration() *ServiceConfiguration
	GetProcesses() map[int]*processman.Process
	Start() error
	Kill() error
	Stop() error
	Restart() error
}

type BaseService struct {
	env  []string
	conf *ServiceConfiguration
	ref  InidService
}

func (svc *BaseService) Init(descrPath string) error {
	buff, err := ioutil.ReadFile(descrPath)
	if err != nil {
		return fmt.Errorf("Error reading service description: %s", err.Error())
	}

	if err := yaml.Unmarshal(buff, &svc.conf); err != nil {
		return fmt.Errorf("Error parsing service configuration: %s", err.Error())
	}

	// Set name of the service, taken from the filename, always lowercase
	svc.conf.SetName(strings.ToLower(strings.Split(path.Base(descrPath), ".")[0]))

	if svc.ref != nil {
		return svc.ref.postInit()
	}

	return nil
}

func (svc *BaseService) SetEnviron(env map[string]string) InidService {
	svc.env = make([]string, 0)
	for k := range env {
		svc.env = append(svc.env, fmt.Sprintf("%s=%s", k, env[k]))
	}
	return svc
}

func (svc *BaseService) GetServiceConfiguration() *ServiceConfiguration { return svc.conf }
func (svc *BaseService) GetProcesses() map[int]*processman.Process      { return nil }
func (svc *BaseService) Start() error                                   { return nil }
func (svc *BaseService) Kill() error                                    { return nil }
func (svc *BaseService) Stop() error                                    { return nil }
func (svc *BaseService) Restart() error                                 { return nil }
func (svc *BaseService) postInit() error                                { return nil }
