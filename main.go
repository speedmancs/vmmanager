package main

import (
	"github.com/speedmancs/vmmanager/app"
)

func main() {
	app := app.App{}
	app.Initialize()
	app.Run(":8010")
}
