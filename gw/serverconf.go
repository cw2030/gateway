package gw

type ServerConf struct {
	Network    string `yaml:"network"	comment:"Network; tcp, tcp4, tcp6, unix or unixpacket"`
	ServerAddr string `yaml:"server_addr"	comment:"listen ip and port string,example:192.168.1.1:3100"`
	Checked    bool   `yaml:"checked"	comment:"check connection is ok or not"`
	IdleTime   uint16 `yaml:"idle_time" comment:"connector idle timeout for seconds"`
	LocalIp    string `yaml:"local_ip"	comment:"localhost ip addr"`
	LocalPort  uint16 `yaml:"local_port"	comment:"localhost port"`
	Encrypt    uint8  `yaml:"encrypt" comment:"set encrypt type"`
	Zuul       string `yaml:"zuul"`
}
