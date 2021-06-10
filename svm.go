package runit_svm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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
	me, err := filepath.Abs(os.Args[0])
	if err != nil {
		return fmt.Errorf("Unable to obtain executable: %s", err.Error())
	}

	meRl := path.Base(me)
	fmt.Println(meRl)
	if !rtutils.InAny(meRl, "1", "2", "3", "init") {
		return fmt.Errorf("Please symlink me at /etc/runit/ to '1', '2' or '3'. Or directly as /sbin/init")
	}

	if meRl == "init" {
		svm.stage = 2
		return nil
	}

	// Put stage for runit
	s, err := strconv.Atoi(meRl)
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
		fmt.Println("Initialised ", s.GetServiceConfiguration().Info, " service")
		svm.services[s.GetServiceConfiguration().Stage].AddSevice(s)
	}

	// Rearrange orders
	fmt.Println("Rearranging")
	for _, order := range svm.services {
		order.Sort()
	}

	return nil
}

func (svm *SVM) Run() {
	// Skip traditional runlevel 1 and 3 competely. Everything is happening in runlevel 2.
	if svm.stage == 1 || svm.stage == 3 {
		return
	}

	// Process runlevels
	for _, runlevel := range []uint8{1, 2, 3} {
		fmt.Printf("Processing stage %d\n", runlevel)
		for _, service := range svm.services[runlevel].GetServices() {
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
