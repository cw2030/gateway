package gw

import "unsafe"

func Str2byte(s string) []byte {
	p := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{p[0], p[1], p[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Byte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
