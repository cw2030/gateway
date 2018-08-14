/*
simple app msg header
include body length and encrypt type, total 4 bytes
*/
package appcodec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

const simpleMessageHeaderLength = 4

type SimpleHeader struct {
	EncryptType uint16
	BodyLength  uint16
}

func NewSimpleHeader() *SimpleHeader {
	return &SimpleHeader{}
}

func (h *SimpleHeader) ToString() string {
	return strconv.FormatInt(int64(h.EncryptType), 10) +
		strconv.FormatInt(int64(h.BodyLength), 10)
}

func (h *SimpleHeader) bytesTo(bs []byte) {
	fmt.Println(bs)
	h.EncryptType = binary.BigEndian.Uint16(bs[:2])
	h.BodyLength = binary.BigEndian.Uint16(bs[2:])
}

func (h *SimpleHeader) toBytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.BigEndian, h.EncryptType)
	binary.Write(buf, binary.BigEndian, h.BodyLength)
	fmt.Println(buf.Bytes())
	return buf.Bytes()
}
