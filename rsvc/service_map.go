package rsvc

import (
	"fmt"
	"sort"
)

type SvmServices struct {
	services map[uint8]*ServiceOrder
}

func NewSvmServices() *SvmServices {
	svs := new(SvmServices)
	svs.services = map[uint8]*ServiceOrder{}
	return svs
}

func (svs *SvmServices) AddService(service *RunitService) {
	if _, so := svs.services[service.GetServiceConfiguration().Stage]; !so {
		svs.services[service.GetServiceConfiguration().Stage] = NewServiceOrder()
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

func (svs *SvmServices) GetServiceByName(name string) (*RunitService, error) {
	for _, so := range svs.services {
		for _, s := range so.services {
			if s.GetServiceConfiguration().GetName() == name {
				return s, nil
			}
		}
	}
	return nil, fmt.Errorf("Service '%s' not found", name)
}

func (svs *SvmServices) GetRunlevels() []*ServiceOrder {
	slots := []*ServiceOrder{}
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
