package rsvc

import (
	"fmt"
	"log"

	"github.com/isbm/inid/services/inidmounter"
)

type MountService struct {
	BaseService
}

func NewMountService() *MountService {
	svc := new(MountService)
	return svc
}

func (svc *MountService) Start() error {
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
