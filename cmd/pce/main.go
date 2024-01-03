package main

import (
	"flag"

	engine "github.com/troppes/portable-container-engine/internal/runtime"
)

// go run main.go run <cmd> <args>
func main() {
	var image string
	var command string

	flag.StringVar(&image, "image", "", "the image name")
	flag.StringVar(&command, "command", "", "the command to run")
	flag.Parse()

	if image == "container" {
		engine.CreateChildProcess(command)
	} else {
		engine.Run(image, command)
	}

}
