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
	conf          ServerConf
	connectors    *sync.Map
	delConnectors *sync.Map
	msgCodec      MessageCodec
	MsgHandler    Handler
	filters       []Filter
	ctx           context.Context
	cancelFunc    context.CancelFunc
	count         Counter
}

var (
	Logger    seelog.LoggerInterface
	logConfig = `
		<seelog type="sync">
			<outputs formatid="main">
				<console/>
			</outputs>
			<formats>
				<format id="main" format="%Date/%Time [%LEV] [%File:%Line] [%Func] %Msg%n"/>
			</formats>
		</seelog>
	`
)

func init() {
	var err error
	Logger, err = seelog.LoggerFromConfigAsBytes([]byte(logConfig))
	if err != nil {
		fmt.Println("seelog init error:", err)
	}
}

type Server interface {
	Listen()
	Dial(network string, remoteAddr string)
	MessageCodec(codec MessageCodec)
	Handler(handler Handler)
	AddFilter(filter Filter)
}

func NewAppMgt(conf ServerConf) *TcpServer {
	ctx, cancel := context.WithCancel(context.Background())
	if conf.IdleTime <= 0 {
		conf.IdleTime = 300
	}
	Logger.Info("Server side configuration: ", conf)
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
		Logger.Info("Listen fail ", network, addr)
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
			Logger.Error(err)
			continue
		}
		netId := s.count.GetAndIncrement()
		connector := NewConnector(conn, netId, s)
		//app.connectors.LoadOrStore(connector, true)
		s.connectors.LoadOrStore(connector, true)
		go connector.process()
		Logger.Infof("accepted client %s, id: %d, total: %d\n",
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
	idleTimeout := float64(s.conf.IdleTime)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			s.connectors.Range(func(key, value interface{}) bool {
				c := key.(*Connector)
				Logger.Debugf("Connector:%d, Closed:%T", c.NetId, c.Closed)
				if c.Closed {
					c.LatestActivity = time.Now()
					s.delConnectors.Store(c, true)
					Logger.Debugf("remove connector from Connectors:%d", c.NetId)
				} else {
					if time.Now().Sub(c.LatestActivity).Seconds() > idleTimeout {
						c.Cancel()
						c.LatestActivity = time.Now()
						s.delConnectors.Store(c, true)
						Logger.Debugf("remove not acitvity connector from Connectors:%d", c.NetId)
					}
				}
				return true
			})
			if s.getConnectorSize() > 0 {
				Logger.Debugf("Current Used Connector have %d", s.getConnectorSize())
			}
			time.Sleep(5 * time.Second)
		}
	}

}

func (s *TcpServer) delConnector() {
	defer func() {
		if err := recover(); err != nil {
			Logger.Error(err)
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
					Logger.Infof("Clean Connector:%d", c.NetId)
				}
				return true
			})
			if s.getDelConnectorSize() > 0 {
				Logger.Infof("Current DelsConnector have %d", s.getDelConnectorSize())
			}
			time.Sleep(10 * time.Second)
		}
	}

}
