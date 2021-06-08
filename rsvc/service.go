package rsvc

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type RunitService struct {
	serialCommands     []*RunitServiceCommand
	concurrentCommands []*RunitServiceCommand
	conf               *ServiceConfiguration
}

func NewRunitService() *RunitService {
	svc := new(RunitService)
	svc.serialCommands = make([]*RunitServiceCommand, 0)
	svc.concurrentCommands = make([]*RunitServiceCommand, 0)

	return svc
}

func (svc *RunitService) Init(descrPath string) error {
	buff, err := ioutil.ReadFile(descrPath)
	if err != nil {
		return fmt.Errorf("Error reading service description: %s", err.Error())
	}

	if err := yaml.Unmarshal(buff, &svc.conf); err != nil {
		return fmt.Errorf("Error parsing service configuration: %s", err.Error())
	}

	svc.loadConcurrentCommands()
	svc.loadSerialCommands()

	return nil
}

func (svc *RunitService) loadSerialCommands() *RunitService {
	for _, command := range svc.conf.Serial {
		svc.serialCommands = append(svc.serialCommands, NewRunitServiceCommand(command).SetConcurrent(false))
	}
	return svc
}

func (svc *RunitService) loadConcurrentCommands() *RunitService {
	for _, command := range svc.conf.Concurrent {
		svc.concurrentCommands = append(svc.concurrentCommands, NewRunitServiceCommand(command).SetConcurrent(true))
	}
	return svc
}

func (svc *RunitService) GetServiceConfiguration() *ServiceConfiguration {
	return svc.conf
}

func (svc *RunitService) Start() error {
	if svc.conf == nil {
		return fmt.Errorf("Service was not initialised!")
	}
	for _, c := range svc.concurrentCommands {
		c.Start()
	}

	for _, c := range svc.serialCommands {
		c.Start()
	}

	return nil
}

func (svc *RunitService) Stop() error {
	return nil
}

// Restart service
func (svc *RunitService) Restart() error {
	if err := svc.Stop(); err != nil {
		return err
	}

	return svc.Start()
}
