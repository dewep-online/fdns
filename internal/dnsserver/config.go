package dnsserver

type ConfigTCP struct {
	Srv Srv `yaml:"srv"`
	DoT DoT `yaml:"dot"`
}

type Srv struct {
	Enable bool   `yaml:"enable"`
	Addr   string `yaml:"addr"`
}

type DoT struct {
	Enable bool   `yaml:"enable"`
	Addr   string `yaml:"addr"`
	Cert   Cert   `yaml:"cert"`
}

type Cert struct {
	Public  string `yaml:"public"`
	Private string `yaml:"private"`
}
