package controller

import (
	"github.com/woremacx/kocha"
)

type Root struct {
	*kocha.DefaultController
}

func (ro *Root) GET(c *kocha.Context) kocha.Result {
	return kocha.Render(c)
}
