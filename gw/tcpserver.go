package gw

import (
	"context"
	"fmt"
	"github.com/cihub/seelog"
	"net"
	"sync"
	"time"
)

type TcpServer struct {
	conf          AppConf
	connectors    *sync.Map
	delConnectors *sync.Map
	msgCodec      MessageCodec
	MsgHandler    Handler
	filters       []Filter
	ctx           context.Context
	cancelFunc    context.CancelFunc
	count         Counter
}

type Server interface {
	Listen()
	Dial(network string, remoteAddr string)
	MessageCodec(codec MessageCodec)
	Handler(handler Handler)
	AddFilter(filter Filter)
}

func NewAppMgt(conf AppConf) *TcpServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &TcpServer{
		ctx:           ctx,
		cancelFunc:    cancel,
		conf:          conf,
		connectors:    &sync.Map{},
		delConnectors: &sync.Map{},
	}
}

func (s *TcpServer) MessageCodec(codec MessageCodec) {
	s.msgCodec = codec
}

func (s *TcpServer) Handler(handler Handler) {
	s.MsgHandler = handler
}

func (s *TcpServer) AddFilter(filter Filter) {

}

func (s *TcpServer) Listen() {
	network := s.conf.Network
	addr := s.conf.ServerAddr
	l, err := net.Listen(network, addr)
	if err != nil {
		fmt.Println("Listen fail ", network, addr)
		return
	}
	defer l.Close()
	//启动状态检查
	go s.checkConnectStatus()
	//启动删除检查
	go s.delConnector()
	for {
		conn, err := l.Accept()
		if err != nil {
			seelog.Error(err)
			continue
		}
		netId := s.count.GetAndIncrement()
		connector := NewConnector(conn, netId, s)
		//app.connectors.LoadOrStore(connector, true)
		s.connectors.LoadOrStore(connector, true)
		go connector.process()
		seelog.Infof("accepted client %s, id: %d, total: %d\n",
			conn.RemoteAddr(),
			netId,
			s.getConnectorSize())
	}

}

func (s *TcpServer) Dial(network string, remoteAddr string) {

}

func (s *TcpServer) getConnectorSize() int {
	count := 0
	s.connectors.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (s *TcpServer) getDelConnectorSize() int {
	count := 0
	s.delConnectors.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (s *TcpServer) checkConnectStatus() {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			s.connectors.Range(func(key, value interface{}) bool {
				c := key.(*Connector)
				seelog.Infof("Connector:%d, Closed:%T", c.NetId, c.Closed)
				if c.Closed {
					c.LatestActivity = time.Now()
					s.delConnectors.Store(c, true)
					seelog.Infof("remove connector from Connectors:%d", c.NetId)
				} else {
					if time.Now().Sub(c.LatestActivity).Seconds() > 30 {
						c.Cancel()
						c.LatestActivity = time.Now()
						s.delConnectors.Store(c, true)
						seelog.Infof("remove not acitvity connector from Connectors:%d", c.NetId)
					}
				}
				return true
			})
			if s.getConnectorSize() > 0 {
				seelog.Infof("Current Used Connector have %d", s.getConnectorSize())
			}
			time.Sleep(5 * time.Second)
		}
	}

}

func (s *TcpServer) delConnector() {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			s.delConnectors.Range(func(key, value interface{}) bool {
				c := key.(*Connector)
				s.connectors.Delete(c)
				if time.Now().Sub(c.LatestActivity).Seconds() > 10 {
					c.Reset()
					s.delConnectors.Delete(c)
					seelog.Infof("Clean Connector:%d", c.NetId)
				}
				return true
			})
			if s.getDelConnectorSize() > 0 {
				seelog.Infof("Current DelsConnector have %d", s.getDelConnectorSize())
			}
			time.Sleep(10 * time.Second)
		}
	}

}
