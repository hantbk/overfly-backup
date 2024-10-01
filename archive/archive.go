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
package archive

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/hantbk/vtsbackup/config"
	"github.com/hantbk/vtsbackup/helper"
	"github.com/hantbk/vtsbackup/logger"
)

// Run archive
func Run(model config.ModelConfig) error {
	logger := logger.Tag("Archive")

	if model.Archive == nil {
		return nil
	}

	if err := helper.MkdirP(model.DumpPath); err != nil {
		logger.Errorf("Failed to mkdir dump path %s: %v", model.DumpPath, err)
		return err
	}

	includes := model.Archive.GetStringSlice("includes")
	includes = cleanPaths(includes)

	excludes := model.Archive.GetStringSlice("excludes")
	excludes = cleanPaths(excludes)

	if len(includes) == 0 {
		return fmt.Errorf("archive.includes have no config")
	}
	logger.Info("=> includes", len(includes), "rules")

	opts := options(model.DumpPath, excludes, includes)

	_, err := helper.Exec("tar", opts...)
	return err
}

func options(dumpPath string, excludes, includes []string) (opts []string) {
	tarPath := path.Join(dumpPath, "archive.tar")
	if helper.IsGnuTar {
		opts = append(opts, "--ignore-failed-read")
	}
	opts = append(opts, "-cPf", tarPath)

	for _, exclude := range excludes {
		opts = append(opts, "--exclude="+filepath.Clean(exclude))
	}

	opts = append(opts, includes...)

	return opts
}

func cleanPaths(paths []string) (results []string) {
	for _, p := range paths {
		results = append(results, filepath.Clean(p))
	}
	return
}
