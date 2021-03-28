package dnsserver

type ConfigTCP struct {
	Server Srv `yaml:"server"`
}

type Srv struct {
	Addr string `yaml:"addr"`
}
