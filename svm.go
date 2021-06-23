package inid

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/isbm/inid/rsvc"
	"github.com/isbm/inid/rtutils"
)

type SVM struct {
	confd      string
	defaultEnv map[string]string
	services   *rsvc.SvmServices
	stage      uint8
}

func NewSVM() *SVM {
	svm := new(SVM)
	svm.defaultEnv = map[string]string{"PATH": "/sbin:/bin:/usr/sbin:/usr/bin"}
	svm.services = rsvc.NewSvmServices()
	svm.confd = "/etc/rc.d"

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
		if _, err := os.Stat(svm.confd); os.IsNotExist(err) {
			return fmt.Errorf("Error accessing %s: %s", svm.confd, err.Error())
		}
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

// Init the svm by reading /etc/rc.d directory
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
		scpath := strings.Split(servConfPath.Name(), ".")
		kind := scpath[len(scpath)-1]

		var s rsvc.InidService
		switch kind {
		case "service":
			s = rsvc.NewRunService()
		case "mount":
			s = rsvc.NewMountService()
		default:
			continue
		}

		s.SetEnviron(svm.defaultEnv)
		if err := s.Init(path.Join(svm.confd, servConfPath.Name())); err != nil {
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
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Oops, panic occurred: %v\n", err)
		}
	}()
	// For runit integration, skip traditional runlevel 1 and 3 competely. Everything is happening in runlevel 2.
	if svm.stage == 1 || svm.stage == 3 {
		return
	}

	for idx, runlevel := range svm.services.GetRunlevels() {
		log.Printf("Entering service stage %d\n", idx+1)
		for _, service := range runlevel.GetServices() {
			if err := service.Start(); err != nil {
				log.Printf("Error occurred: %s", err.Error())
			}
		}
	}
	NewIPCServer(svm.services).ServeForever()
}
