package config

type Deli struct {
	Season string `json:"season" mapstructure:"season" yaml:"season"` // winter|summer 季节
	Stop   bool   `json:"stop" yaml:"stop" mapstructure:"stop"`       // true|false 关闭此功能
}
