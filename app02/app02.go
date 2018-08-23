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
	h := appcodec.NewSimpleHeader()
	h.EncryptType = 1

	b := appcodec.Body{}
	b.BType = "1"
	b.SessionId = "1234"
	b.ProtType = "login"
	b.SvrName = "UserLogin"
	b.Resource = "/user/login"
	b.Content = "userName=test&password=1234"
	b.Attachment = "NO"
	msg := appcodec.SimpleMessage{h, &b}
	bs := msg.Encode()
	fmt.Println("send:", bs)
	conn.Write(bs)

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("timer out for 5 secs")
	default:
		codec := appcodec.SimpleMessageCodec{}
		response, err := codec.Decode(conn)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(response.ToString())
		}

	}
}
