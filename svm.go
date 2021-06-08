package runit_svm

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/isbm/runit-svm/rsvc"
)

type SVM struct {
	confd      string
	services   []*rsvc.RunitService
	defaultEnv map[string]string
}

func NewSVM() *SVM {
	svm := new(SVM)
	svm.services = make([]*rsvc.RunitService, 0)
	svm.defaultEnv = map[string]string{"PATH": "/sbin:/bin:/usr/sbin:/usr/bin"}
	svm.confd = "/etc/runit.d"

	return svm
}

// Init the svm by reading /etc/runit.d directory
func (svm *SVM) Init() error {
	filenames, err := ioutil.ReadDir(svm.confd)
	if err != nil {
		return fmt.Errorf("Error reading %s directory: %s", svm.confd, err.Error())
	}

	for _, servConfPath := range filenames {
		if !strings.HasSuffix(servConfPath.Name(), ".service") {
			continue
		}
		s := rsvc.NewRunitService()
		spath := path.Join(svm.confd, servConfPath.Name())
		if err := s.Init(spath); err != nil {
			return err
		}
		svm.services = append(svm.services, s)
	}

	return nil
}

func (svm *SVM) Run() error {
	for _, service := range svm.services {
		service.Start()
	}

	return nil
}
