//go:build linux

package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	img "github.com/troppes/portable-container-engine/internal/image"
	util "github.com/troppes/portable-container-engine/internal/util"
)

type platformRuntime struct{}

func (r *platformRuntime) Run(image string, command []string) error {
	imagePath, err := img.RetrieveImage(image, true)
	if err != nil {
		return err
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

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (r *platformRuntime) CreateChildProcess(path string, command []string) error {
	fmt.Println("Current command: " + strings.Join(command, " "))
	fmt.Println("Current path on host:" + path)

	var cmd *exec.Cmd
	if len(command) > 1 {
		cmd = exec.Command(command[0], command[1:]...)
	} else {
		cmd = exec.Command(command[0])
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	util.Must(syscall.Sethostname([]byte("container")))
	util.Must(syscall.Chroot(path))
	util.Must(os.Chdir("/"))
	util.Must(syscall.Mount("proc", "proc", "proc", 0, ""))

	util.Must(cmd.Run())

	util.Must(syscall.Unmount("proc", 0))

	return nil
}
