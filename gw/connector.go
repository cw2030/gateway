/*
表示每一个连接到服务端的连接工作者
*/
package gw

import (
	"context"
	"math/rand"
	"net"
	"time"
)

type Connector struct {
	mgt            *AppMgt
	Conn           net.Conn
	SinceNow       time.Time
	LatestActivity time.Time
	messageCodec   MessageCodec
	handlerFunc    Handler
	filters        []Filter
	writeChan      chan []byte
	ctx            context.Context
	cancelFunc     context.CancelFunc
	key            []byte
}

func NewConnector(conn net.Conn, appMgt *AppMgt) *Connector {
	ctx, cancel := context.WithCancel(appMgt.ctx)
	return &Connector{Conn: conn,
		mgt:            appMgt,
		SinceNow:       time.Now(),
		LatestActivity: time.Now(),
		messageCodec:   appMgt.msgCodec,
		handlerFunc:    appMgt.handlerFunc,
		ctx:            ctx,
		cancelFunc:     cancel,
		key:            RandEncryptKey(16),
	}
}

/*
process client request and write response
*/
func (connector *Connector) process() {

	go connector.write()
	for {
		select {
		case <-connector.ctx.Done():
			break
		case <-connector.mgt.ctx.Done():
			break
		default:
			msg, err := connector.messageCodec.Decode(connector.Conn)
			connector.handlerFunc.HandleFunc(connector, msg, err)
		}
	}

}

func (connector *Connector) write() {
	select {
	case bs := <-connector.writeChan:
		connector.Conn.Write(bs)
	case <-connector.ctx.Done():
		break
	case <-connector.mgt.ctx.Done():
		break
	}
}

func RandEncryptKey(size int) []byte {
	kinds, result := []int{26, 65}, make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		scope, base := kinds[0], kinds[1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
