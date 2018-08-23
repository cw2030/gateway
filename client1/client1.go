package main

import (
	"encoding/hex"
	"fmt"
	"gateway/appcodec"
	"gateway/gw"
	"math/rand"
	"net"
	"os"
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
	connector := &gw.Connector{Conn: conn, Conf: gw.ServerConf{}}
	conn.Write(msg.Encode(connector))

	go func() {
		for {
			codec := appcodec.StringMessageCodec{}
			response, err := codec.Decode(conn, connector)
			if err != nil {
				gw.Logger.Error(err)
				os.Exit(0)
			}
			resp, ok := response.(*appcodec.StringMessage)
			if !ok {
				return
			}

			switch resp.Header.MsgType {
			case appcodec.Msg_Type_Handshake:
				key, err = hex.DecodeString(resp.Body.Content)
				connector.Conf.Encrypt = resp.Header.EncryptType
				if err != nil {
					gw.Logger.Errorf("Parse EncryptKey fail:%s,err:%s", string(key), err.Error())
				} else {
					gw.Logger.Infof("Key:%s", resp.Body.Content)
					connector.Key = key
				}

			case appcodec.Msg_type_Heartbeat:
				gw.Logger.Infof("Receive Heartbeat: %s", resp.ToString())
			case appcodec.Msg_Type_Login:
				gw.Logger.Infof("Receive LoginMsg: %s", resp.ToString())
			case appcodec.Msg_Type_Busi:
			case appcodec.Msg_Type_Async:
				gw.Logger.Infof("Receive Business Message Resposne: %s", resp.ToString())
			default:
				gw.Logger.Infof("Receive UnknowMsg: %s", resp.ToString())
			}

		}
	}()
	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case <-time.After(5 * time.Second):
			ts := rand.Intn(100)
			gw.Logger.Infof("ts:%d", ts)
			if ts%2 == 0 {
				m := appcodec.NewEmptyMsg()
				sm := m.(*appcodec.StringMessage)
				h := sm.Header
				b := sm.Body

				h.ReqType = appcodec.Request
				h.MsgType = appcodec.Msg_Type_Login
				h.EncryptType = connector.Conf.Encrypt

				b.BType = "b"
				b.SvrName = "login"
				b.Resource = "/act/login"
				b.Action = "POST"
				b.Content = "userName=test&password=1234"
				conn.Write(sm.Encode(connector))
				gw.Logger.Infof("Send Login Msg: %s", sm.ToString())

			} else {
				/*msg = appcodec.NewHeatbeatReqMsg()
				conn.Write(msg.Encode(connector))
				gw.Logger.Infof("Send HeartBeat:%s", msg.ToString())*/
				busiMsg := appcodec.NewEmptyMsg().(*appcodec.StringMessage)
				h := busiMsg.Header
				body := busiMsg.Body
				h.ReqType = appcodec.Request
				h.MsgType = appcodec.Msg_Type_Busi
				h.EncryptType = gw.Encrypt_AES

				body.BType = "b"
				body.SvrName = "bw-transaction-monitor"
				body.Resource = "/health"
				body.Action = "GET"
				conn.Write(busiMsg.Encode(connector))
				gw.Logger.Infof("Send Business Message:%s", busiMsg.ToString())
			}

		}
	}

}
