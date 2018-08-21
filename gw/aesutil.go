package gw

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func encrypt(key []byte, iv []byte, src string) []byte {
	if src == "" || key == nil {
		return nil
	}
	srcBytes := Str2byte(src)
	padding := aes.BlockSize - len(srcBytes)%aes.BlockSize
	if padding > 0 {
		pad := bytes.Repeat([]byte{byte(padding)}, padding)
		srcBytes = append(srcBytes, pad...)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		Logger.Error(err)
		return nil
	}
	cipherBytes := make([]byte, aes.BlockSize+len(srcBytes))
	io.ReadFull(rand.Reader, iv)

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherBytes[aes.BlockSize:], srcBytes)
	return cipherBytes[aes.BlockSize:]
}

func decrypt(key []byte, iv []byte, src []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		Logger.Error(err)
		return nil
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	result := make([]byte, len(src))
	mode.CryptBlocks(result, src)
	return result
}
