package main

import (
	app2 "tranc/server/app"
	"tranc/server/config"
)

func main() {
	cfg := new(config.Config)
	cfg.InitFile()
	app := app2.InitApp(*cfg)
	app.Run()
}
