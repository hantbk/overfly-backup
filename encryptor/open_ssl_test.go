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

package encryptor

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestOpenSSL_options(t *testing.T) {
	base := &Base{
		viper:       viper.New(),
		archivePath: "/foo/bar",
	}

	enc := NewOpenSSL(base)
	assert.Equal(t, false, enc.base64)
	assert.Equal(t, true, enc.salt)
	assert.Equal(t, "", enc.password)
	assert.Equal(t, "/foo/bar.enc", enc.encryptPath)
	assert.Equal(t, "", enc.args)
	assert.Equal(t, "aes-256-cbc", enc.chiper)
	assert.Equal(t, "aes-256-cbc -salt -k ", strings.Join(enc.options(), " "))

	base.viper.Set("base64", true)
	base.viper.Set("salt", false)
	base.viper.Set("args", "-pbkdf2 -iter 1000")
	base.viper.Set("password", "backup-123")
	base.viper.Set("chiper", "rc4")

	enc = NewOpenSSL(base)
	assert.Equal(t, true, enc.base64)
	assert.Equal(t, false, enc.salt)
	assert.Equal(t, "rc4", enc.chiper)
	assert.Equal(t, "backup-123", enc.password)
	assert.Equal(t, "-pbkdf2 -iter 1000", enc.args)

	assert.Equal(t, "rc4 -base64 -pbkdf2 -iter 1000 -k backup-123", strings.Join(enc.options(), " "))
}
