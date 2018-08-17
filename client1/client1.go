package main

import (
	"encoding/hex"
	"fmt"
	"gateway/appcodec"
	"github.com/cihub/seelog"
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
	msg := appcodec.NewHandShakeReqMsg()
	var key []byte
	conn.Write(msg.Encode())
	go func() {
		for {
			codec := appcodec.StringMessageCodec{}
			response, err := codec.Decode(conn)

			resp := response.(*appcodec.StringMessage)

			if err != nil {
				seelog.Error(err)
				return
			}
			switch resp.Header.MsgType {
			case appcodec.Msg_Type_Handshake:
				key, err = hex.DecodeString(resp.Body.Content)
				if err != nil {
					seelog.Errorf("Parse EncryptKey fail:%s,err:%s", string(key), err.Error())
				} else {
					seelog.Infof("Key:%s", resp.Body.Content)
				}

			case appcodec.Msg_type_Heartbeat:
				seelog.Infof(resp.ToString())
			}

		}
	}()
	for {
		select {
		case <-time.After(10 * time.Second):
			msg = appcodec.NewHeatbeatReqMsg()
			conn.Write(msg.Encode())
			seelog.Infof("Write Heatbeat Message:%s", msg.ToString())
		}
	}

}
