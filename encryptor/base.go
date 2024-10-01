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
	"github.com/hantbk/vtsbackup/config"
	"github.com/hantbk/vtsbackup/logger"
	"github.com/spf13/viper"
)

// Base encryptor
type Base struct {
	model       config.ModelConfig
	viper       *viper.Viper
	archivePath string
}

// Encryptor interface
type Encryptor interface {
	perform() (encryptPath string, err error)
}

func newBase(archivePath string, model config.ModelConfig) (base *Base) {
	base = &Base{
		archivePath: archivePath,
		model:       model,
		viper:       model.EncryptWith.Viper,
	}
	return
}

// Run compressor
func Run(archivePath string, model config.ModelConfig) (encryptPath string, err error) {
	logger := logger.Tag("Encryptor")

	base := newBase(archivePath, model)
	var enc Encryptor
	switch model.EncryptWith.Type {
	case "openssl":
		enc = NewOpenSSL(base)
	default:
		encryptPath = archivePath
		return
	}

	logger.Info("encrypt: " + model.EncryptWith.Type)
	encryptPath, err = enc.perform()
	if err != nil {
		return
	}
	logger.Info("encrypted:", encryptPath)

	// save Extension
	model.Viper.Set("Ext", model.Viper.GetString("Ext")+".enc")

	return
}
