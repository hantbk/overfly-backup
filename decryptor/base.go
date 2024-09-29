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
