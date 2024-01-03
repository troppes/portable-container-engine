package main

import (
	"os"

	engine "github.com/troppes/portable-container-engine/internal/engine"
)

// go run main.go run <cmd> <args>
func main() {
	switch os.Args[1] {
	case "run":
		engine.Run()
	case "container":
		engine.CreateChildProcess()
	default:
		panic("Please enter a valid command")
	}
}
