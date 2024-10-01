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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExistsPath(t *testing.T) {
	exist := IsExistsPath("foo/bar")
	assert.False(t, exist)

	exist = IsExistsPath("./filepath_test.go")
	assert.True(t, exist)
}

func TestMkdirP(t *testing.T) {
	dest := path.Join(os.TempDir(), "test-mkdir-p")
	exist := IsExistsPath(dest)
	assert.False(t, exist)

	assert.Nil(t, MkdirP(dest))
	defer os.Remove(dest)
	exist = IsExistsPath(dest)
	assert.True(t, exist)
}

func TestExplandHome(t *testing.T) {
	newPath := ExplandHome("")
	assert.Equal(t, newPath, "")

	newPath = ExplandHome("/home/hant/111")
	assert.Equal(t, newPath, "/home/hant/111")

	newPath = ExplandHome("~")
	assert.Equal(t, newPath, "~")

	newPath = ExplandHome("~/")
	assert.NotEqual(t, newPath[:2], "~/")

	newPath = ExplandHome("~/foo/bar/dar")
	assert.Equal(t, newPath, path.Join(os.Getenv("HOME"), "/foo/bar/dar"))
}

func TestAbsolutePath(t *testing.T) {
	pwd, _ := os.Getwd()
	newPath := AbsolutePath("foo/bar")
	assert.Equal(t, newPath, path.Join(pwd, "foo/bar"))

	newPath = AbsolutePath("/home/hant/111")
	assert.Equal(t, newPath, "/home/hant/111")

	newPath = AbsolutePath("~")
	assert.NotEqual(t, newPath[:2], "~/")

	newPath = AbsolutePath("~/")
	assert.NotEqual(t, newPath[:2], "~/")

	newPath = AbsolutePath("~/foo/bar/dar")
	assert.Equal(t, newPath, path.Join(os.Getenv("HOME"), "/foo/bar/dar"))
}
