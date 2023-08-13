package dns

type (
	Config struct {
		DNS ConfigItem `yaml:"dns"`
	}
	ConfigItem struct {
		Addr string `yaml:"addr"`
	}
)

func (v *Config) Default() {
	v.DNS.Addr = "0.0.0.0:53"
}
