package router

import "cgin/router/system"

type RouterGroup struct {
	System system.RouteGroup
}

var RouterGroupApp = new(RouterGroup)
