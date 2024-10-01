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
	"strings"
	"testing"

	"github.com/hantbk/vtsbackup/config"
	"github.com/hantbk/vtsbackup/helper"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	// with nil Archive
	model := config.ModelConfig{
		Archive: nil,
	}
	err := Run(model)
	assert.NoError(t, err)
}

func TestOptions(t *testing.T) {
	includes := []string{
		"/foo/bar/dar",
		"/bar/foo",
		"/ddd",
	}

	excludes := []string{
		"/hello/world",
		"/cc/111",
	}

	dumpPath := "~/work/dir"

	opts := options(dumpPath, excludes, includes)
	cmd := strings.Join(opts, " ")
	if helper.IsGnuTar {
		assert.Equal(t, cmd, "--ignore-failed-read -cPf ~/work/dir/archive.tar --exclude=/hello/world --exclude=/cc/111 /foo/bar/dar /bar/foo /ddd")
	} else {
		assert.Equal(t, cmd, "-cPf ~/work/dir/archive.tar --exclude=/hello/world --exclude=/cc/111 /foo/bar/dar /bar/foo /ddd")
	}
}
