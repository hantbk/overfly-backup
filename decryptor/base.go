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

package decryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type Decryptor struct {
	key []byte
}

func NewDecryptor(key string) (*Decryptor, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes long")
	}

	return &Decryptor{key: []byte(key)}, nil
}

func (d *Decryptor) Decrypt(encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(d.key)
	if err != nil {
		return nil, err
	}

	nonce := encryptedData[:aes.BlockSize]
	ciphertext := encryptedData[aes.BlockSize:]

	stream := cipher.NewCBCDecrypter(block, nonce)
	stream.CryptBlocks(ciphertext, ciphertext)

	return ciphertext, nil
}
