package gw

type Handler interface {
	HandleFunc(connector *Connector, message Message, err error)
}
