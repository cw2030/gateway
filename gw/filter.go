package gw

type Filter interface {
	filter(message Message) error
}
