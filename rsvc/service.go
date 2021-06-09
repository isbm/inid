package rsvc

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/isbm/processman"
	"gopkg.in/yaml.v2"
)

type RunitService struct {
	serialCommands     []*RunitServiceCommand
	concurrentCommands []*RunitServiceCommand
	conf               *ServiceConfiguration
	procman            *processman.Processman
	env                []string
}

func NewRunitService() *RunitService {
	svc := new(RunitService)
	svc.serialCommands = make([]*RunitServiceCommand, 0)
	svc.concurrentCommands = make([]*RunitServiceCommand, 0)
	svc.procman = processman.New(nil)

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

	// Set name of the service, taken from the filename, always lowercase
	svc.conf.SetName(strings.ToLower(strings.Split(path.Base(descrPath), ".")[0]))

	svc.loadConcurrentCommands()
	svc.loadSerialCommands()

	return nil
}

func (svc *RunitService) SetEnviron(env map[string]string) *RunitService {
	svc.env = make([]string, 0)
	for k := range env {
		svc.env = append(svc.env, fmt.Sprintf("%s=%s", k, env[k]))
	}
	return svc
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
		go func(rsc *RunitServiceCommand) {
			_, err := svc.procman.Command(rsc.command, rsc.args, svc.env)
			if err != nil {
				fmt.Printf("Error running concurrent command '%s': %s\n", rsc.command, err.Error())
			}
		}(c)
	}

	for _, c := range svc.serialCommands {
		_, err := svc.procman.Command(c.command, c.args, svc.env)
		if err != nil {
			fmt.Printf("Error running serial command '%s': %s\n", c.command, err.Error())
		}
	}

	return nil
}

func (svc *RunitService) Kill() error {
	return svc.procman.KillAll()
}

func (svc *RunitService) Stop() error {
	return svc.procman.StopAll()
}

// Restart service
func (svc *RunitService) Restart() error {
	if err := svc.Stop(); err != nil {
		return err
	}

	return svc.Start()
}
