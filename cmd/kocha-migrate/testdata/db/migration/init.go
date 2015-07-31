package migration

import (
	"fmt"

	"github.com/woremacx/kocha"
)

type Migration struct{}

func Up(config kocha.DatabaseConfig, n int) error {
	fmt.Printf("call Up: n => %v\n", n)
	return nil
}

func Down(config kocha.DatabaseConfig, n int) error {
	fmt.Printf("call Down: n => %v\n", n)
	return nil
}
