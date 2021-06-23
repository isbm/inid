package rsvc

import "fmt"

type ServiceOrder struct {
	services []InidService
}

func NewServiceOrder() *ServiceOrder {
	so := new(ServiceOrder)
	so.services = make([]InidService, 0)
	return so
}

// AddService in order
func (so *ServiceOrder) AddSevice(service InidService) *ServiceOrder {
	so.services = append(so.services, service)
	return so
}

func (so *ServiceOrder) inArray(v string, a []string) bool {
	for _, n := range a {
		if n == v {
			return true
		}
	}
	return false
}

func (so *ServiceOrder) popServices(names []string) {
	svc := []InidService{}
	for _, s := range so.services {
		if !so.inArray(s.GetServiceConfiguration().GetName(), names) {
			svc = append(svc, s)
		}
	}
	so.services = svc
}

func (so *ServiceOrder) insertOrder(set []InidService, service InidService) ([]InidService, bool) {
	buff := []InidService{}
	added := false
	for _, s := range set {
		if s.GetServiceConfiguration().GetName() == service.GetServiceConfiguration().Before {
			buff = append(buff, service)
			buff = append(buff, s)
			added = true
		} else if s.GetServiceConfiguration().GetName() == service.GetServiceConfiguration().After {
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
	ordered := []InidService{}

	// Add services those have no deps
	removed := []string{}
	for _, service := range so.services {
		if service.GetServiceConfiguration().After == "" && service.GetServiceConfiguration().Before == "" {
			ordered = append(ordered, service)
			removed = append(removed, service.GetServiceConfiguration().GetName())
		}
	}

	// Delete moved
	so.popServices(removed)

	// Stop, if there are no root services to hook on
	if len(ordered) == 0 {
		fmt.Println("There are no independent services found.")
		return
	}

	// Add all after.
	cycler := 0
	for {
		dependency := false
		removed = []string{}
		for _, cdx := range []int{1, 2} {
			for _, service := range so.services {
				var crit string

				switch cdx {
				case 1:
					crit = service.GetServiceConfiguration().After
				case 2:
					crit = service.GetServiceConfiguration().Before
				}

				if crit != "" {
					dependency = true
					var added bool
					ordered, added = so.insertOrder(ordered, service)
					if added {
						removed = append(removed, service.GetServiceConfiguration().GetName())
					}
				}
			}
		}
		// Delete moved
		so.popServices(removed)

		cycler++

		if len(so.services) == 0 || !dependency {
			break
		}
		if cycler > 10000 {
			fmt.Println("Service chain seems forever cycled. Giving up...")
			break
		}
	}

	so.services = ordered
}

// GetServices
func (so *ServiceOrder) GetServices() []InidService {
	return so.services
}
