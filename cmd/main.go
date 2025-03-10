package main

import (
	"github.com/enigmasterr/parallel_calculator/internal/application"
)

func main() {
	app := application.New()
	//app.Run()
	app.RunServer()
}
