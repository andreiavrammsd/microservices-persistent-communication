package main

type ServiceRequest struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
	Body   string `json:"body"`
}

type Service struct {
	Address string
}

type Services map[string]Service

var services = Services{
	"notification": Service{
		Address: "http://zonga.ro",
	},
}

func serviceIsDefined(name string) bool {
	_, isDefined := services[name]
	return isDefined
}
