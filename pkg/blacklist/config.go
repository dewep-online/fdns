package blacklist

type Config struct {
	BlackHoleIP string   `yaml:"blackholeip"`
	BlackListIP []string `yaml:"blacklistip"`
}
