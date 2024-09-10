package helper

import (
	"os"
	"path"
)

// IsExistsPath check if path exists
func IsExistsPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// MkdirPath create directory path if not exists
func MkdirP(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0777)
		if err != nil {
			return
		}
	}
}

// ExplandHome ~/foo -> /home/hant/foo
func ExplandHome(filePath string) string {
	if len(filePath) < 2 {
		return filePath
	}

	if filePath[:2] != "~/" {
		return filePath
	}

	return path.Join(os.Getenv("HOME"), filePath[2:])
}
