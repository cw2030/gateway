package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
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

func main() {
	fmt.Printf("%s", time.Now().Format("2006-01-02 15:04:05"))
	println()
	header := Header{ProtocolFlag: 182,
		MainVersion: 49, SecondaryVersion: 48, ReqType: 49, MsgType: 49, EncryptType: 50, Length: 102, Priority: 49, Extend: 30}
	buf := make([]byte, 12)

	binary.BigEndian.PutUint16(buf, header.ProtocolFlag)
	buf[2] = header.MainVersion
	buf[3] = header.SecondaryVersion
	buf[4] = header.ReqType
	buf[5] = header.MsgType
	buf[6] = header.EncryptType
	buf[7] = byte(header.Length >> 8)
	buf[8] = byte(header.Length)
	buf[9] = header.Priority
	buf[10] = byte(header.Extend >> 8)
	buf[11] = byte(header.Extend)
	fmt.Println(hex.EncodeToString(buf))
	fmt.Println(buf)

	s := "1.0"
	fmt.Println([]byte(s))

	var r uint = 1

	fmt.Println(strconv.FormatUint(uint64(r), 10))

}
