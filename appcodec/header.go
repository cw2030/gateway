package appcodec

import (
	"encoding/binary"
	"encoding/hex"
)

var (
	protocolFlag     = []byte{0xAA, 0xA1}
	mainVersion      = byte(49)
	secondaryVersion = byte(48)
	eof              = string(0x1)
	fieldSplit       = string(0x0)
)

const (
	headerLength = 12
)

type Header struct {
	ProtocolFlag     uint16
	MainVersion      uint8
	SecondaryVersion uint8
	ReqType          uint8
	MsgType          uint8
	EncryptType      uint8
	Length           uint16
	Priority         uint8
	Extend           uint16
}

func NewHeader() *Header {
	return &Header{ProtocolFlag: binary.BigEndian.Uint16(protocolFlag),
		MainVersion:      mainVersion,
		SecondaryVersion: secondaryVersion,
	}
}

func (h *Header) ToString() string {
	return hex.EncodeToString(h.toBytes())
}

func (h *Header) bytesTo(bs []byte) {
	h.ProtocolFlag = binary.BigEndian.Uint16(bs[:2])
	h.MainVersion = bs[2]
	h.SecondaryVersion = bs[3]
	h.ReqType = bs[4]
	h.MsgType = bs[5]
	h.EncryptType = bs[6]
	h.Length = binary.BigEndian.Uint16(bs[7:9])
	h.Priority = bs[9]
	h.Extend = binary.BigEndian.Uint16(bs[10:])
}

func (h *Header) toBytes() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf, h.ProtocolFlag)
	buf[2] = h.MainVersion
	buf[3] = h.SecondaryVersion
	buf[4] = h.ReqType
	buf[5] = h.MsgType
	buf[6] = h.EncryptType
	buf[7] = byte(h.Length >> 8)
	buf[8] = byte(h.Length)
	buf[9] = h.Priority
	buf[10] = byte(h.Extend >> 8)
	buf[11] = byte(h.Extend)

	return buf
}
