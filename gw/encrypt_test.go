package gw

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"testing"
)

func TestEncrypt(t *testing.T) {
	src := "wo shi shui"
	key := make([]byte, 16)
	iv := make([]byte, 16)
	io.ReadFull(rand.Reader, iv)
	io.ReadFull(rand.Reader, key)
	result := encrypt(key, iv, src)
	fmt.Println(hex.EncodeToString(result))

	strResult := decrypt(key, iv, result)
	fmt.Println(Byte2str(strResult))

}
