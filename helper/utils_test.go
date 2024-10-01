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
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestUtils_init(t *testing.T) {
	if runtime.GOOS == "linux" {
		assert.Equal(t, IsGnuTar, true)
	} else {
		assert.Equal(t, IsGnuTar, false)
	}
}

func TestCleanHost(t *testing.T) {
	assert.Equal(t, "foo.bar.com", CleanHost("foo.bar.com"))
	assert.Equal(t, "foo.bar.com", CleanHost("ftp://foo.bar.com"))
	assert.Equal(t, "foo.bar.com", CleanHost("http://foo.bar.com"))
	assert.Equal(t, "", CleanHost("http://"))
}

func TestFormatEndpoint(t *testing.T) {
	assert.Equal(t, "http://foo.bar.com", FormatEndpoint("http://foo.bar.com"))
	assert.Equal(t, "https://foo.bar.com", FormatEndpoint("https://foo.bar.com"))
	assert.Equal(t, "https://foo.bar.com", FormatEndpoint("https://foo.bar.com"))
}
