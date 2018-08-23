package appcodec

import (
	"gateway/gw"
)

var iv = []byte{0x10, 0xA0, 0x37, 0xF2, 0xC5, 0x8A, 0x28, 0xB1, 0x8C, 0xEA, 0xF3, 0xA6, 0x22, 0x11, 0xAE, 0x52}

func EncryptAES(key []byte, content string) []byte {
	/*if body != nil {
		return body.ToBytes()
	}*/
	return gw.EncryptByCBC(key, iv, content)
}

func DecryptAES(key []byte, body []byte) []byte {
	/*if body != nil {
		b := &Body{}
		b.BytesTo(body)
		return b
	}*/
	bodyBytes := gw.DecryptByCBC(key, iv, body)
	length := len(bodyBytes)
	pad := int(bodyBytes[length-1])
	/*b := &Body{}
	b.BytesTo(bodyBytes[0 : length-pad])*/
	return bodyBytes[0 : length-pad]
}

func EncryptRSA2048(body []byte) []byte {

	return nil
}

func DecryptRSA2048(body []byte) []byte {

	return nil
}
