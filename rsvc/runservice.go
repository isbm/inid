package rsvc

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/isbm/inid/rtutils"
	"github.com/isbm/processman"
)

type RunService struct {
	serialCommands     []*RunServiceCommand
	concurrentCommands []*RunServiceCommand
	procman            *processman.Processman
	BaseService
}

func NewRunService() *RunService {
	svc := new(RunService)
	svc.serialCommands = make([]*RunServiceCommand, 0)
	svc.concurrentCommands = make([]*RunServiceCommand, 0)
	svc.procman = processman.New(nil)
	svc.ref = svc

	return svc
}

// Called at the end of the main Init()
func (svc *RunService) postInit() error {
	for _, command := range svc.conf.GetSerialCommands() {
		svc.serialCommands = append(svc.serialCommands, NewRunServiceCommand(command).SetConcurrent(false))
	}

	for _, command := range svc.conf.GetConcurrentCommands() {
		svc.concurrentCommands = append(svc.concurrentCommands, NewRunServiceCommand(command).SetConcurrent(true))
	}

	return nil
}

func (svc *RunService) GetProcesses() map[int]*processman.Process {
	return svc.procman.Processes()
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

func (svc *RunService) Start() error {
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
