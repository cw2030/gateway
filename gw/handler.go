package gw

type Handler interface {
	HandleFunc(message Message) error
}
