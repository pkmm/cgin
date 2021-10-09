package core

import (
	"cgin/global"
	"cgin/initialize"
	"fmt"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.DB != nil {
		// TODO
	}
	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.Config.System.Addr)

	s := initServer(address, Router)
	fmt.Printf("Cgin is running: http://127.0.0.1%s\n", address)
	global.GLog.Error(s.ListenAndServe().Error())
}
