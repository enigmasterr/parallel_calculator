package main

import (
	"github.com/enigmasterr/calchttp/internal/application"
)

func main() {
	app := application.New()
	//app.Run()
	app.RunServer()
}
