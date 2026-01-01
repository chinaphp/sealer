// Copyright © 2021 Alibaba Group Holding Ltd.
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

package hash

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

const aesKey = "ZU9WbzRMVXRQZ2pzTGowR2hNWUpIZjRkWld4aWVRWko="

func AesEncrypt(origData []byte) (string, error) {
	key, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode key base64: %v", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData = pkcs7Padding(origData, blockSize)

	// Create random IV
	ciphertext := make([]byte, blockSize+len(origData))
	iv := ciphertext[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCEncrypter(block, iv) // #nosec G407
	blockMode.CryptBlocks(ciphertext[blockSize:], origData)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func AesDecrypt(ciphertext []byte) (string, error) {
	key, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode key base64: %v", err)
	}
	ciphertext, err = base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return "", fmt.Errorf("failed to decode key base64: %v", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()

	// Check if ciphertext is long enough to contain IV
	if len(ciphertext) < blockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]

	// Check validity after extracting IV
	if len(ciphertext)%blockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	plaintext := pkcs7UnPadding(ciphertext)
	return string(plaintext), nil
}

func pkcs7Padding(origData []byte, blockSize int) []byte {
	padding := blockSize - len(origData)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padText...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return nil
	}
	unPadding := int(origData[length-1])
	if length < unPadding {
		return nil
	}
	return origData[:(length - unPadding)]
}
