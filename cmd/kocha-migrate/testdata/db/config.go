package db

import "github.com/woremacx/kocha"

var DatabaseMap = kocha.DatabaseMap{
	"default": {
		Driver: "sqlite3",
		DSN:    ":memory:",
	},
}
