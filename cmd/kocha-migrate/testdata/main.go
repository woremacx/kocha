package main

import (
	"testappname/config"

	"github.com/woremacx/kocha"
)

func main() {
	kocha.Run(config.AppConfig)
}
