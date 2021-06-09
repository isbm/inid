package runit_svm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/isbm/runit-svm/rsvc"
	"github.com/isbm/runit-svm/rtutils"
)

type SVM struct {
	confd      string
	services   map[uint8]*rsvc.ServiceOrder
	defaultEnv map[string]string
	stage      uint8
}

func NewSVM() *SVM {
	svm := new(SVM)
	svm.services = map[uint8]*rsvc.ServiceOrder{
		1: rsvc.NewServiceOrder(),
		2: rsvc.NewServiceOrder(),
		3: rsvc.NewServiceOrder(),
	}
	svm.defaultEnv = map[string]string{"PATH": "/sbin:/bin:/usr/sbin:/usr/bin"}
	svm.confd = "/etc/runit.d"

	return svm
}

// Set the runlevel
func (svm *SVM) setRunlevel() error {
	me, err := os.Executable()
	if err != nil {
		return fmt.Errorf("Unable to obtain executable: %s", err.Error())
	}

	if !rtutils.InAny(me, "1", "2", "3") {
		return fmt.Errorf("Please symlink me at /etc/runit/ to '1', '2' or '3'.")
	}

	s, err := strconv.Atoi(me)
	if err != nil {
		return err
	}
	svm.stage = uint8(s)

	return nil
}

// Init the svm by reading /etc/runit.d directory
func (svm *SVM) Init() error {
	if err := svm.setRunlevel(); err != nil {
		return err
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
		fmt.Println("Initialised ", s.GetServiceConfiguration().Info, " service")
		svm.services[s.GetServiceConfiguration().Stage].AddSevice(s)
	}

	return nil
}

func (svm *SVM) Run() error {
	// Skip traditional runlevel 1 competely. Everything is happening in runlevel 2.
	if svm.stage == 1 {
		return nil
	}

	// Process runlevels
	for _, runlevel := range []uint8{1, 2, 3} {
		fmt.Printf("Processing stage %d\n", runlevel)
		for _, service := range svm.services[runlevel].GetServices() {
			fmt.Print("Starting ", service.GetServiceConfiguration().Info, " ... ")
			if err := service.Start(); err != nil {
				fmt.Println("Failed")
				return err
			}
			fmt.Println("Done")
		}
	}

	// Forever loop
	for {
		select {}
	}
}
