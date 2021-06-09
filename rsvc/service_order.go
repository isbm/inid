package rsvc

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
	if len(so.services) == 0 {
		so.services = append(so.services, service)
	} else {
		slist := make([]*RunitService, 0)
		added := false

		if service.conf.After != "" { // Before works only if no after defined
			for _, current := range so.services {
				slist = append(slist, current)
				if current.conf.GetName() == service.conf.After {
					slist = append(slist, service)
					added = true
				}
			}
		} else if service.conf.Before != "" {
			for _, current := range so.services {
				if current.conf.GetName() == service.conf.Before {
					slist = append(slist, service)
					added = true
				}
				slist = append(slist, current)
			}
		} else {
			slist = append(slist, so.services...)
		}

		if !added {
			slist = append(slist, service)
		}
		so.services = slist
	}
	return so
}

// GetServices
func (so *ServiceOrder) GetServices() []*RunitService {
	return so.services
}
