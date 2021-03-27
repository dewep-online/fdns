package webserver

//go:generate easyjson

//easyjson:json
type (
	//MiddlewareConfig model
	MiddlewareConfig struct {
		Middleware ConfigItem `yaml:"middleware" json:"middleware"`
	}
	//ConfigItem model
	ConfigItem struct {
		Throttling int64 `yaml:"throttling" json:"throttling"`
	}
)
