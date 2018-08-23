package appcodec

import (
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

func TestAppEncryptAES(t *testing.T) {
	key := make([]byte, 16)
	io.ReadFull(rand.Reader, key)
	body := Body{BType: "b", SessionId: "sid10234", SvrName: "loginHandler", Action: "POST", Resource: "/act/login",
		Content: "userName=test\u0026password=12341"}

	bs := EncryptAES(key, &body)
	fmt.Sprintf("%x", bs)

	bb := DecryptAES(key, bs)

	fmt.Println(bb.ToString())

}
