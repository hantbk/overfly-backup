// Copyright © 2024 Ha Nguyen <captainnemot1k60@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package helper

import (
	"os"
	"path"
	"path/filepath"

	"github.com/hantbk/vtsbackup/logger"
)

// IsExistsPath check path exist
func IsExistsPath(p string) bool {
	_, err := os.Stat(p)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// MkdirP like mkdir -p
func MkdirP(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0750)
	}
	return nil
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

// Convert a file path into an absolute path
func AbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	path = ExplandHome(path)

	path, err := filepath.Abs(path)
	if err != nil {
		logger.Error("Convert config file path to absolute path failed: ", err)
		return path
	}

	return path
}
