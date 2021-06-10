package runit_svm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/isbm/runit-svm/rsvc"
	"github.com/isbm/runit-svm/rtutils"
)

type SvmServices struct {
	services map[uint8]*rsvc.ServiceOrder
}

func NewSvmServices() *SvmServices {
	svs := new(SvmServices)
	svs.services = map[uint8]*rsvc.ServiceOrder{}
	return svs
}

func (svs *SvmServices) AddService(service *rsvc.RunitService) {
	if _, so := svs.services[service.GetServiceConfiguration().Stage]; !so {
		svs.services[service.GetServiceConfiguration().Stage] = rsvc.NewServiceOrder()
	}
	svs.services[service.GetServiceConfiguration().Stage].AddSevice(service)
}

func (svs *SvmServices) GetStages() []uint8 {
	stages := []uint8{}
	for stage := range svs.services {
		stages = append(stages, stage)
	}
	return stages
}

func (svs *SvmServices) GetRunlevels() []*rsvc.ServiceOrder {
	slots := []*rsvc.ServiceOrder{}
	idx := []int{}

	for key := range svs.services {
		idx = append(idx, int(key))
	}

	sort.Ints(idx)

	for _, i := range idx {
		slots = append(slots, svs.services[uint8(i)])
	}

	return slots
}

type SVM struct {
	confd      string
	defaultEnv map[string]string
	services   *SvmServices
	stage      uint8
}

func NewSVM() *SVM {
	svm := new(SVM)
	svm.defaultEnv = map[string]string{"PATH": "/sbin:/bin:/usr/sbin:/usr/bin"}
	svm.services = NewSvmServices()
	svm.confd = "/etc/runit.d"

	return svm
}

// Set the runlevel
func (svm *SVM) setRunlevel() error {
	me, err := filepath.Abs(os.Args[0])
	if err != nil {
		return fmt.Errorf("Unable to obtain executable: %s", err.Error())
	}

	meBase := path.Base(me)
	if !rtutils.InAny(meBase, "1", "2", "3", "init") {
		return fmt.Errorf("Please symlink me at /etc/runit/ to '1', '2' or '3'. Or directly as /sbin/init")
	}

	// Compat to runit, run as level 2
	if meBase == "init" {
		svm.stage = 2
		return nil
	}

	// Put stage for runit
	s, err := strconv.Atoi(meBase)
	if err != nil {
		return err
	}
	svm.stage = uint8(s)

	if me != path.Join("/etc/runit", strconv.Itoa(int(svm.stage))) {
		return fmt.Errorf("I must be put as /etc/runit/%d, not as %s", svm.stage, me)
	}

	return nil
}

// Init the svm by reading /etc/runit.d directory
func (svm *SVM) Init() error {
	if err := svm.setRunlevel(); err != nil {
		return err
	}
	// Skip init during runlevel 1 and 3.
	if svm.stage == 1 || svm.stage == 3 {
		return nil
	}

	filenames, err := ioutil.ReadDir(svm.confd)
	if err != nil {
		return fmt.Errorf("Error reading %s directory: %s", svm.confd, err.Error())
	}

	for _, servConfPath := range filenames {
		if !strings.HasSuffix(servConfPath.Name(), ".service") {
			continue
		}
		s := rsvc.NewRunitService().SetEnviron(svm.defaultEnv)
		spath := path.Join(svm.confd, servConfPath.Name())
		if err := s.Init(spath); err != nil {
			return err
		}
		svm.services.AddService(s)
	}

	// Rearrange runlevel ordering
	for _, runlevel := range svm.services.GetRunlevels() {
		runlevel.Sort()
	}

	return nil
}

func (svm *SVM) Run() {
	// For runit integration, skip traditional runlevel 1 and 3 competely. Everything is happening in runlevel 2.
	if svm.stage == 1 || svm.stage == 3 {
		return
	}

	for idx, runlevel := range svm.services.GetRunlevels() {
		fmt.Printf("Processing stage %d\n", idx+1)
		for _, service := range runlevel.GetServices() {
			fmt.Print("Starting ", service.GetServiceConfiguration().Info, " ... ")
			if err := service.Start(); err != nil {
				fmt.Println("Failed")
			}
			fmt.Println("Done")
		}
	}

	// Forever loop
	for {
		select {}
	}
}
