package config

import (
	"testappname/app/controller"

	"github.com/woremacx/kocha"
)

type RouteTable kocha.RouteTable

var routes = RouteTable{
	{
		Name:       "root",
		Path:       "/",
		Controller: &controller.Root{},
	},
}

func init() {
	AppConfig.RouteTable = kocha.RouteTable(append(routes, RouteTable{
		{
			Name:       "static",
			Path:       "/*path",
			Controller: &kocha.StaticServe{},
		},
	}...))

}
