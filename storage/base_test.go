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

package storage

import (
	"github.com/hantbk/vtsbackup/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBase_newBase(t *testing.T) {
	model := config.ModelConfig{}
	archivePath := "/tmp/vtsbackup/test-storage/foo.zip"
	s, _ := newBase(model, archivePath, config.SubConfig{})

	assert.Equal(t, s.archivePath, archivePath)
	assert.Equal(t, s.model, model)
	assert.Equal(t, s.viper, model.Viper)
	assert.Equal(t, s.keep, 0)
}
