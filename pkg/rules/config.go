package rules

type (
	Config struct {
		DumpRules    string            `yaml:"dump-rules"`
		AdblockRules []string          `yaml:"adblock-rules"`
		DNSRules     map[string]string `yaml:"dns-rules"`
		HostRules    map[string]string `yaml:"host-rules"`
		RegExpRules  map[string]string `yaml:"regexp-rules"`
		QueryRules   map[string]string `yaml:"query-rules"`
	}
)
