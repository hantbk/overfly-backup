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

package compressor

import (
	"os/exec"

	"github.com/hantbk/vtsbackup/helper"
)

type Tar struct {
	Base
}

func (tar *Tar) perform() (archivePath string, err error) {
	filePath := tar.archiveFilePath(tar.ext)

	opts := tar.options()
	opts = append(opts, filePath)
	opts = append(opts, tar.name)
	archivePath = filePath

	_, err = helper.Exec("tar", opts...)

	return
}

func (tar *Tar) options() (opts []string) {
	if helper.IsGnuTar {
		opts = append(opts, "--ignore-failed-read")
	}

	var useCompressProgram bool
	if len(tar.parallelProgram) > 0 {
		if path, err := exec.LookPath(tar.parallelProgram); err == nil {
			useCompressProgram = true
			opts = append(opts, "--use-compress-program", path)
		}
	}
	if !useCompressProgram {
		opts = append(opts, "-a")
	}
	opts = append(opts, "-cf")

	return
}
