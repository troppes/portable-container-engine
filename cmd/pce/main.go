package main

import (
	"fmt"
	"os"

	dl "github.com/troppes/portable-container-engine/internal/image"
	pce "github.com/troppes/portable-container-engine/internal/runtime"
)

// go run main.go run <cmd> <args>
func main() {

	args := os.Args

	// Check if the number of arguments is at least 3 (program name + two arguments)
	if len(args) < 4 {
		fmt.Println("Please provide at least three arguments. Usage: pce <download|run> <image> <command>")
		return
	}

	// Retrieve the first and second arguments
	mode := args[1]
	path := args[2]
	command := args[3:]

	if path == "" {
		fmt.Println("Please provide an image path")
		return
	}

	switch mode {
	case "run":
		pce.Run(path, command)
	case "internalrun":
		pce.CreateChildProcess(path, command)
	case "download":
		fmt.Printf("Downloading image %v\n", path)
		dl.DownloadImage(path)
	default:
		fmt.Printf("Unknown command %v\n", mode)
	}
}
