package util

import (
	"errors"
	"runtime"
)

// IsLinux returns true if the current OS is Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// CheckPlatformSupport returns an error if the platform is not supported
func CheckPlatformSupport() error {
	if !IsLinux() {
		return errors.New("this operation requires Linux. Please run this application inside Docker or on a Linux system")
	}
	return nil
}
