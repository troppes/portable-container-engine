package main

import (
	"fmt"
	"os"

	dl "github.com/troppes/portable-container-engine/internal/image"
	pce "github.com/troppes/portable-container-engine/internal/runtime"
)

func main() {

	args := os.Args

	// Check if the number of arguments is at least 3 (program name + two arguments)
	if args[1] != "download" && len(args) < 4 {
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
		dlPath, err := dl.DownloadImage(path)
		if err != nil {
			fmt.Printf("Error downloading image %v\n", err)
		} else {
			fmt.Printf("Image downloaded to %v\n", dlPath)
		}
	default:
		fmt.Printf("Unknown command %v\n", mode)
	}
}
