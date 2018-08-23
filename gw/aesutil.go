package gw

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func EncryptByCBC(key []byte, iv []byte, src string) []byte {
	if src == "" || key == nil {
		return nil
	}
	Logger.Debugf("Source Content：%s", src)
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

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherBytes[aes.BlockSize:], srcBytes)
	Logger.Debugf("Hex Bytes for Encrypt：%x", cipherBytes[aes.BlockSize:])
	return cipherBytes[aes.BlockSize:]
}

func DecryptByCBC(key []byte, iv []byte, src []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		Logger.Errorf("NewCipher failure: %s", err.Error())
		return nil
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	result := make([]byte, len(src))
	mode.CryptBlocks(result, src)

	Logger.Debugf("Content of Decryption：%s", Byte2str(result))
	return result
}
