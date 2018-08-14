package gw

import (
	"context"
	"fmt"
	"net"
	"time"
)

type AppMgt struct {
	conf        AppConf
	connectors  map[interface{}]bool
	msgCodec    MessageCodec
	handlerFunc Handler
	filters     []Filter
	ctx         context.Context
	cancelFunc  context.CancelFunc
	count       Counter
}

type App interface {
	Listen()
	Dial(network string, remoteAddr string)
	MessageCodec(codec MessageCodec)
	Handler(handler Handler)
	AddFilter(filter Filter)
}

func NewAppMgt(conf AppConf) *AppMgt {
	ctx, cancel := context.WithCancel(context.Background())
	return &AppMgt{
		ctx:        ctx,
		cancelFunc: cancel,
		conf:       conf,
		connectors: make(map[interface{}]bool, 20),
	}
}

func (app *AppMgt) MessageCodec(codec MessageCodec) {
	app.msgCodec = codec
}

func (app *AppMgt) Handler(handler Handler) {
	app.handlerFunc = handler
}

func (app *AppMgt) AddFilter(filter Filter) {

}

func (app *AppMgt) Listen() {
	network := app.conf.Network
	addr := app.conf.ServerAddr
	l, err := net.Listen(network, addr)
	if err != nil {
		fmt.Println("Listen fail ", network, addr)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		connector := NewConnector(conn, app)
		//app.connectors.LoadOrStore(connector, true)
		app.connectors[connector] = true
		go connector.process()
		netId := app.count.GetAndIncrement()
		fmt.Printf("accepted client %s, time:%s id: %d, total: %d\n",
			conn.RemoteAddr(),
			time.Now().Format("2006-01-02 15:04:05"),
			netId,
			app.getConnectorSize())
		fmt.Println()
	}

}

func (app *AppMgt) Dial(network string, remoteAddr string) {

}

func (app *AppMgt) getConnectorSize() int {
	return len(app.connectors)
}
