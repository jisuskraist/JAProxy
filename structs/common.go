package structs

//Service defines the functionality that all services must follow.
type Service interface {
	Start()
	Stop()
}

//Services represents an array of service as a hub.
type Services struct {
	s []Service
}

//AddService add a service to the list of services,
func (s *Services) AddService(service Service) {
	s.s = append(s.s, service)
}

//Run starts all the registered services
func (s Services) Run() {
	for _, srv := range s.s {
		srv.Start()
	}
}

//Stop stops all the registered services
func (s Services) Stop() {
	for _, srv := range s.s {
		srv.Stop()
	}
}
