package rsvc

import "fmt"

type ServiceOrder struct {
	services []*RunitService
}

func NewServiceOrder() *ServiceOrder {
	so := new(ServiceOrder)
	so.services = make([]*RunitService, 0)
	return so
}

// AddService in order
func (so *ServiceOrder) AddSevice(service *RunitService) *ServiceOrder {
	so.services = append(so.services, service)
	return so
}

func (so *ServiceOrder) popService(idx int) {
	so.services[len(so.services)-1], so.services[idx] = so.services[idx], so.services[len(so.services)-1]
	so.services = so.services[:len(so.services)-1]
}

func (so *ServiceOrder) addBefore(set []*RunitService, service *RunitService) ([]*RunitService, bool) {
	buff := []*RunitService{}
	added := false
	for _, s := range set {
		if s.conf.GetName() == service.conf.Before {
			buff = append(buff, service)
			buff = append(buff, s)
			added = true
		} else {
			buff = append(buff, s)
		}
	}
	return buff, added
}

func (so *ServiceOrder) addAfter(set []*RunitService, service *RunitService) ([]*RunitService, bool) {
	buff := []*RunitService{}
	added := false
	for _, s := range set {
		if s.conf.GetName() == service.conf.After {
			buff = append(buff, s)
			buff = append(buff, service)
			added = true
		} else {
			buff = append(buff, s)
		}
	}
	return buff, added
}

func (so *ServiceOrder) Sort() {
	ordered := []*RunitService{}

	// Add services those have no deps
	removed := []int{}
	for idx, service := range so.services {
		if service.conf.After == "" && service.conf.Before == "" {
			ordered = append(ordered, service)
			removed = append(removed, idx)
		}
	}

	// Delete moved
	for _, idx := range removed {
		so.popService(idx)
	}

	// Stop, if there are no root services to hook on
	if len(ordered) == 0 {
		fmt.Println("There are no independent services found.")
		return
	}

	// Add all after.
	cycler := 0
	for {
		dependency := false
		removed = []int{}
		for idx, service := range so.services {
			if service.conf.After != "" {
				dependency = true
				var added bool
				ordered, added = so.addAfter(ordered, service)
				if added {
					removed = append(removed, idx)
				}
			}
		}
		// Delete moved
		for _, idx := range removed {
			so.popService(idx)
		}

		removed = []int{}
		for idx, service := range so.services {
			if service.conf.Before != "" {
				fmt.Println("Found candidate for before:", service.conf.GetName())
				dependency = true
				var added bool
				ordered, added = so.addBefore(ordered, service)
				if added {
					removed = append(removed, idx)
				}
			}
		}
		// Delete moved
		for _, idx := range removed {
			so.popService(idx)
		}

		cycler++

		if len(so.services) == 0 || !dependency {
			break
		}
		if cycler > 10000 {
			fmt.Println("Service chain is wrong, cannot quit ordering. Giving up...")
			break
		}
	}

	so.services = ordered

	fmt.Println("Services in this scope:")
	for idx, s := range so.services {
		fmt.Printf("%d. %s\n", idx+1, s.conf.serviceName)
	}
}

// GetServices
func (so *ServiceOrder) GetServices() []*RunitService {
	return so.services
}
