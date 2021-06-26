package main

import "github.com/speedmancs/vmmanager/app"

func main() {
	app := app.App{}
	app.Initialize("logs.txt")
	app.Run(":8010")
}
