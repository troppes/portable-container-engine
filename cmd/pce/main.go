package main

import (
	"fmt"
	"os"

	dl "github.com/troppes/portable-container-engine/internal/image"
	pce "github.com/troppes/portable-container-engine/internal/runtime"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("Usage: pce <download|run> <image> [<command>...]")
		return
	}

	mode := args[1]
	image := args[2]
	command := []string{}
	if len(args) > 3 {
		command = args[3:]
	}

	if image == "" {
		fmt.Println("Please provide a valid image name")
		return
	}

	// Get the platform-appropriate container runtime
	containerRuntime := pce.GetRuntime()

	switch mode {
	case "run":
		if len(command) == 0 {
			fmt.Println("Please provide a command to run")
			return
		}
		if err := containerRuntime.Run(image, command); err != nil {
			fmt.Printf("Error running container: %v\n", err)
			return
		}

	case "internalrun":
		if len(command) == 0 {
			fmt.Println("Please provide a command for internal run")
			return
		}
		if err := containerRuntime.CreateChildProcess(image, command); err != nil {
			fmt.Printf("Error creating child process: %v\n", err)
			return
		}

	case "download":
		extract := false
		if len(command) > 0 && (command[0] == "--extract" || command[0] == "-x") {
			extract = true
		}

		fmt.Printf("Downloading image %v (extract: %v)\n", image, extract)
		dlPath, err := dl.RetrieveImage(image, extract)
		if err != nil {
			fmt.Printf("Error downloading image: %v\n", err)
		} else {
			if extract {
				fmt.Printf("Image extracted to %v\n", dlPath)
			} else {
				fmt.Printf("Image downloaded to %v\n", dlPath)
			}
		}

	default:
		fmt.Printf("Unknown command: %v\n", mode)
		fmt.Println("Usage: pce <download|run|internalrun> <image> [<command|--extract>...]")
	}
}
