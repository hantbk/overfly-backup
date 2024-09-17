package helper

import (
	"fmt"
	"os"
	"runtime"
)

func CheckOS() (string, string, error) {
	goos := runtime.GOOS
	var info string

	switch goos {
	case "darwin":
		info = "darwin"
	case "linux":
		info = "linux"
	default:
		return "", "", fmt.Errorf("unsupported operating system: %s", goos)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", "", err
	}

	return info, hostname, nil
}
