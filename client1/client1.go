package main

import (
	"fmt"
	"gateway/appcodec"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp4", ":7722")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	h := appcodec.NewHeader()
	h.ReqType = 1
	h.MsgType = 1
	h.EncryptType = 2
	h.Priority = 5
	h.Extend = 11

	b := appcodec.Body{}
	b.BType = "1"
	b.SessionId = "1234"
	b.SvrType = "login"
	b.SvrName = "UserLogin"
	b.Resource = "/user/login"
	b.Content = "userName=test&password=1234"
	b.Attachment = "NO"
	msg := appcodec.StringMessage{h, &b}
	conn.Write(msg.Encode())
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("timer out for 5 secs")
	default:
		codec := appcodec.StringMessageCodec{}
		response, err := codec.Decode(conn)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(response.ToString())
		}

	}
}
