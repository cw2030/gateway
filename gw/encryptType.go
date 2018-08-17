package gw

const (
	Encrypt_None uint8 = iota
	Encrypt_AES
	Encrypt_RSA1024
	Encrypt_RSA2048
	Encrypt_SM4 //国密对称加密算法
	Encrypt_SM2 //国密非对称加密算法
)
