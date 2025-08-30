//go:build linux

package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	img "github.com/troppes/portable-container-engine/internal/image"
	util "github.com/troppes/portable-container-engine/internal/util"
)

type platformRuntime struct {
}

func (r *platformRuntime) Run(image string, command []string) error {
	// Create temporary directory for this container run
	tempDir, err := os.MkdirTemp("", "container-runtime-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Schedule cleanup - this will run after the process finishes
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to cleanup temp directory %s: %v\n", tempDir, err)
		}
	}()

	// Retrieve image to the temporary directory
	imagePath, config, err := img.RetrieveImage(image, true, tempDir)
	if err != nil {
		return err
	}

	if len(command) == 0 {
		command = getDefaultCommand(config)
		if len(command) == 0 {
			return fmt.Errorf("no command specified and no default command found in image")
		}
	}

	// restart myself with the child flag /proc/self/exe is a symbolic link to the current process
	args := append([]string{"internalrun", imagePath}, command...)

	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// NEWNS => used for mounting
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Credential:   &syscall.Credential{Uid: 0, Gid: 0},                                    // make root in container
		UidMappings:  []syscall.SysProcIDMap{{ContainerID: 0, HostID: os.Getuid(), Size: 1}}, // outside of container be the user
		GidMappings:  []syscall.SysProcIDMap{{ContainerID: 0, HostID: os.Getgid(), Size: 1}},
		Unshareflags: syscall.CLONE_NEWNS, // remove the other mounts
	}

	// Handle signals in parent process
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if err := cmd.Start(); err != nil {
		return err
	}

	// Signal handling goroutine
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, shutting down container...")

		// Send SIGTERM first
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			fmt.Printf("Warning: failed to send SIGTERM: %v\n", err)
		}

		// Wait 5 seconds for graceful shutdown
		time.Sleep(5 * time.Second)

		// Force kill if still running
		if err := cmd.Process.Kill(); err != nil {
			fmt.Printf("Warning: failed to kill process: %v\n", err)
		}
	}()

	return cmd.Wait()
}

func (r *platformRuntime) CreateChildProcess(path string, command []string) error {
	fmt.Println("Current command: " + strings.Join(command, " "))
	fmt.Println("Current path on host:" + path)

	// Simple environment setup
	os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")

	util.Must(syscall.Sethostname([]byte("container")))
	util.Must(syscall.Chroot(path))
	util.Must(os.Chdir("/"))

	// Create essential device files
	os.MkdirAll("/dev", 0755)
	os.MkdirAll("/proc", 0755)

	// Create /dev/null
	if f, err := os.Create("/dev/null"); err == nil {
		f.Close()
	}

	// Create /dev/urandom with some random data
	if f, err := os.Create("/dev/urandom"); err == nil {
		f.WriteString("random_data_placeholder_" + fmt.Sprintf("%d", time.Now().UnixNano()))
		f.Close()
	}

	// Create /dev/random with some random data
	if f, err := os.Create("/dev/random"); err == nil {
		f.WriteString("random_data_placeholder_" + fmt.Sprintf("%d", time.Now().UnixNano()))
		f.Close()
	}

	// Create /dev/zero
	if f, err := os.Create("/dev/zero"); err == nil {
		f.Close()
	}

	util.Must(syscall.Mount("proc", "proc", "proc", 0, ""))

	err := syscall.Exec(command[0], command, os.Environ())
	return fmt.Errorf("exec failed: %v", err)
}

func getDefaultCommand(config *v1.ConfigFile) []string {
	var fullCommand []string

	if len(config.Config.Entrypoint) > 0 {
		fullCommand = append(fullCommand, config.Config.Entrypoint...)
		if len(config.Config.Cmd) > 0 {
			fullCommand = append(fullCommand, config.Config.Cmd...)
		}
	} else if len(config.Config.Cmd) > 0 {
		fullCommand = config.Config.Cmd
	}

	return fullCommand
}
