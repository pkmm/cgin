package system

import "cgin/service"

type ApiGroup struct {
	SystemApi
}

var apiService = service.ServiceGroupApp.SystemServiceGroup.ApiService
var deliAutoSignService = service.ServiceGroupApp.SystemServiceGroup.DeliAutoSign