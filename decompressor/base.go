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

package decompressor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hantbk/vtsbackup/logger"
)

func Run(filePath string, modelName string) error {
	logger := logger.Tag("Decompressor")
	logger.Infof("Decompressing %s...", filePath)

	extractDir := filepath.Dir(filePath)

	// Ensure the extract directory exists
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("failed to create extract directory: %v", err)
	}

	// Use tar command to extract the file
	cmd := exec.Command("tar", "-xzvf", filePath, "-C", extractDir)

	// Create a buffer to capture the command output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to decompress file: %v\nOutput: %s", err, out.String())
	}

	logger.Infof("File decompressed successfully to: %s", extractDir)

	extractDir = filepath.Join(extractDir, modelName)

	// Search for archive.tar in the extract directory
	archiveTarPath := filepath.Join(extractDir, "archive.tar")
	if _, err := os.Stat(archiveTarPath); err == nil {
		logger.Infof("Found archive.tar at: %s", archiveTarPath)
		logger.Infof("Extracting archive.tar...")
		cmd = exec.Command("tar", "-xvf", archiveTarPath, "-C", extractDir)
		cmd.Stdout = &out
		cmd.Stderr = &out
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to extract archive.tar: %v\nOutput: %s", err, out.String())
		}
		logger.Infof("archive.tar extracted successfully")
		// Remove archive.tar after extraction
		if err := os.Remove(archiveTarPath); err != nil {
			logger.Warnf("Failed to remove archive.tar: %v", err)
		}
	} else {
		logger.Infof("archive.tar not found at %s, skipping extraction", archiveTarPath)
	}

	cleanUp(filePath)

	return nil
}

func cleanUp(filePath string) {
	// Delete the original .tar.gz file
	if err := os.Remove(filePath); err != nil {
		logger.Warnf("Failed to remove original file: %v", err)
	} else {
		logger.Infof("Original file removed: %s", filePath)
	}
}
