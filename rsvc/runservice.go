package rsvc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/isbm/inid/rtutils"
	"github.com/isbm/inid/services/inidmounter"
	"github.com/isbm/processman"
	"gopkg.in/yaml.v2"
)

type RunService struct {
	serialCommands     []*RunServiceCommand
	concurrentCommands []*RunServiceCommand
	conf               *ServiceConfiguration
	procman            *processman.Processman
	BaseService
}

func NewRunService() *RunService {
	svc := new(RunService)
	svc.serialCommands = make([]*RunServiceCommand, 0)
	svc.concurrentCommands = make([]*RunServiceCommand, 0)
	svc.procman = processman.New(nil)

	return svc
}

func (svc *RunService) Init(descrPath string) error {
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

func (svc *RunService) loadSerialCommands() *RunService {
	for _, command := range svc.conf.GetSerialCommands() {
		svc.serialCommands = append(svc.serialCommands, NewRunServiceCommand(command).SetConcurrent(false))
	}
	return svc
}

func (svc *RunService) loadConcurrentCommands() *RunService {
	for _, command := range svc.conf.GetConcurrentCommands() {
		svc.concurrentCommands = append(svc.concurrentCommands, NewRunServiceCommand(command).SetConcurrent(true))
	}
	return svc
}

func (svc *RunService) GetServiceConfiguration() *ServiceConfiguration {
	return svc.conf
}

func (svc *RunService) GetProcesses() map[int]*processman.Process {
	return svc.procman.Processes()
}

func (svc *RunService) Start() error {
	switch svc.GetServiceConfiguration().GetKind() {
	case SVC_SERVICE:
		return svc.startService()
	case SVC_MOUNTER:
		return svc.startMounter()
	default:
		return fmt.Errorf("Unknown service type")
	}
}

func (svc *RunService) startMounter() error {
	log.Printf("Starting mounter service")
	if svc.conf == nil {
		return fmt.Errorf("Mounter service was not initialised!")
	}

	// Start concurrent mounts
	go func(conf *ServiceConfiguration) {
		if err := inidmounter.NewInidPremounter(conf.GetConcurrentConfig()).Start(); err != nil {
			log.Printf("Concurrent mounter %s failed: %s", svc.GetServiceConfiguration().GetName(), err.Error())
		}
	}(svc.GetServiceConfiguration())

	// Start serial mounts
	if err := inidmounter.NewInidPremounter(svc.conf.GetSerialConfig()).Start(); err != nil {
		log.Printf("Serial mounter %s failed: %s", svc.GetServiceConfiguration().GetName(), err.Error())
	}

	return nil
}

func (svc *RunService) formatSTD(p *processman.Process) string {
	stream, err := p.Stdout()
	var buff bytes.Buffer
	if err == nil {
		o := strings.TrimSpace(rtutils.RCloser2String(stream))
		if o != "" {
			buff.WriteString(o + "\n")
		}
	}

	stream, err = p.Stderr()
	if err == nil {
		o := strings.TrimSpace(rtutils.RCloser2String(stream))
		if o != "" {
			buff.WriteString(o + "\n")
		}
	}

	if len(buff.Bytes()) > 0 {
		return fmt.Sprintf("Output log for process %s:\n%s\n", svc.GetServiceConfiguration().GetName(), buff.String())
	}

	return ""
}

func (svc *RunService) startService() error {
	if svc.conf == nil {
		return fmt.Errorf("Service was not initialised!")
	}
	failed := false
	for _, c := range svc.concurrentCommands {
		go func(rsc *RunServiceCommand) {
			p, err := svc.procman.StartConcurrent(rsc.command, rsc.args, svc.env)
			if err != nil {
				log.Printf("Service %s failed background command '%s': %s\n", svc.GetServiceConfiguration().GetName(), rsc.command, err.Error())
				failed = true
			} else if out := svc.formatSTD(p); out != "" {
				log.Println(out)
			}
		}(c)
	}

	for _, c := range svc.serialCommands {
		p, err := svc.procman.StartSerial(c.command, c.args, svc.env)
		if err != nil {
			log.Printf("Service %s failed serial command '%s': %s\n", svc.GetServiceConfiguration().GetName(), c.command, err.Error())
			failed = true
		} else if out := svc.formatSTD(p); out != "" {
			log.Println(out)
		}
	}
	if !failed {
		log.Printf("Service %s started successfully\n", svc.GetServiceConfiguration().GetName())
	} else {
		log.Printf("Service %s failed\n", svc.GetServiceConfiguration().GetName())
	}

	return nil
}

func (svc *RunService) Kill() error {
	return svc.procman.KillAll()
}

func (svc *RunService) Stop() error {
	return svc.procman.StopAll()
}

// Restart service
func (svc *RunService) Restart() error {
	if err := svc.Stop(); err != nil {
		return err
	}

	return svc.Start()
}
