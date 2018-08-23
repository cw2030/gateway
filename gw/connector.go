/*
表示每一个连接到服务端的连接工作者
*/
package gw

import (
	"context"
	"crypto/rand"
	"io"
	"net"
	"time"
)

type Connector struct {
	Tcpserver      *TcpServer
	Conn           net.Conn
	SinceNow       time.Time
	LatestActivity time.Time
	Codec          MessageCodec
	MsgHandler     Handler
	PreFilter      []Filter
	WriteChan      chan []byte
	Ctx            context.Context
	Cancel         context.CancelFunc
	Key            []byte
	NetId          int64
	Closed         bool
	Conf           ServerConf
}

func NewConnector(conn net.Conn, netId int64, server *TcpServer) *Connector {
	ctx, cancel := context.WithCancel(server.ctx)
	key := getEncryptKey(server.conf.Encrypt)
	return &Connector{Conn: conn,
		Tcpserver:      server,
		SinceNow:       time.Now(),
		LatestActivity: time.Now(),
		Codec:          server.msgCodec,
		MsgHandler:     server.MsgHandler,
		Ctx:            ctx,
		Cancel:         cancel,
		NetId:          netId,
		WriteChan:      make(chan []byte, 128),
		Conf:           server.conf,
		Key:            key,
	}
}

/*
process client request and write response
*/
func (connector *Connector) process() {
	defer func() {
		if err := recover(); err != nil {
			connector.closeConn()
			Logger.Error(err)
		}
	}()
	go connector.write()
	for {
		select {
		case <-connector.Ctx.Done():
			Logger.Info("Read:Close conn in Write when client error or exit.")
			return
		case <-connector.Tcpserver.ctx.Done():
			Logger.Info("Read:Close conn when server exit.")
			break
		default:
			msg, err := connector.Codec.Decode(connector.Conn, connector)
			if err != nil {
				if err.Error() == "EOF" {
					Logger.Error("NetID：%d，Connector close:%T", connector.NetId, connector)
					connector.Cancel()
				}
				switch err {
				case io.EOF:
				case io.ErrUnexpectedEOF:
				case io.ErrClosedPipe:
				case io.ErrShortBuffer:
				case io.ErrShortWrite:
					Logger.Error(err)
					Logger.Infof("NetID：%d，Connector close:%T", connector.NetId, connector)
					connector.Cancel()
				default:
					_, flag := err.(*net.OpError)
					if flag {
						connector.Cancel()
					}
					Logger.Error(err)
				}
				continue
			}
			Logger.Debugf("Receive Message:%s", msg.ToString())
			if connector.MsgHandler != nil {
				//设置最新活动时间
				connector.LatestActivity = time.Now()
				connector.MsgHandler.HandleFunc(connector, msg, err)
			} else {
				Logger.Errorf("Read:Can't find MsgHandler.Msg:%s", msg)
			}
		}
	}

}

/*
监听数据回写
*/
func (connector *Connector) write() {
	defer func() {
		if err := recover(); err != nil {
			connector.closeConn()
			Logger.Error(err)
		}
	}()
	for {
		select {
		case bs := <-connector.WriteChan:
			_, err := connector.Conn.Write(bs)
			if err != nil {
				Logger.Error(err)
				connector.closeConn()
				connector.Cancel()
			}
		case <-connector.Ctx.Done():
			Logger.Info("Write:Close conn in Write when client error or exit.")
			connector.closeConn()
			return
		case <-connector.Tcpserver.ctx.Done():
			Logger.Info("Write:Close conn when server exit.")
			connector.closeConn()
			return
		}
	}

}

/*
关闭当前客户端连接
*/
func (connector *Connector) closeConn() {
	//标记该客户端连接为关闭状态
	connector.Closed = true
	err := connector.Conn.Close()
	if err != nil {
		Logger.Error(err)
	}
}

/*
release resource
*/
func (connector *Connector) Reset() {
	connector.Tcpserver = nil
	connector.Conn = nil
	connector.Cancel = nil
	connector.MsgHandler = nil
	connector.Ctx = nil
	connector.Key = nil
	if connector.WriteChan != nil {
		close(connector.WriteChan)
	}
	connector.MsgHandler = nil
	connector.Codec = nil
	connector.PreFilter = nil
}

/*
按指定位数生成随机的密钥
*/
func RandAESEncryptKey(size int) []byte {
	var keySize int
	switch keySize {
	case 16:
	case 24:
	case 32:
		keySize = size
	default:
		keySize = 16
	}
	key := make([]byte, keySize)
	io.ReadFull(rand.Reader, key)
	return key
	/*kinds, result := []int{26, 65}, make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		scope, base := kinds[0], kinds[1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result*/
}

func RandRSA2048EncryptKey() []byte {
	result := make([]byte, 2048)
	return result
}

func getEncryptKey(encryptType uint8) []byte {
	switch encryptType {
	case Encrypt_AES:
		return RandAESEncryptKey(16)
	case Encrypt_RSA2048:
		return RandRSA2048EncryptKey()
	default:
		return nil

	}
}
