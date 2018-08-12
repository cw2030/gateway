/*
表示每一个连接到服务端的连接工作者
*/
package main

import (
	"context"
	"fmt"
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
		case msg, err := connector.messageCodec.Decode(connector.Conn):
			fmt.Println("ReadMsg:", msg, err)
			connector.writeChan <- []byte("server:" + msg.ToString())
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
