//go:build darwin

package runtime

import (
	"fmt"
	"runtime"
)

type platformRuntime struct{}

func (r *platformRuntime) Run(image string, command []string) error {
	return fmt.Errorf("container functionality is not supported on %s. Please use Linux", runtime.GOOS)
}

func (r *platformRuntime) CreateChildProcess(path string, command []string) error {
	return fmt.Errorf("container functionality is not supported on %s. Please use Linux", runtime.GOOS)
}
