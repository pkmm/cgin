package config

type Server struct {
	JWT      JWT      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Redis    Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	Mysql    Mysql    `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Zap      Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
	Local    Local    `mapstructure:"local" json:"local" yaml:"local"`
	System   System   `mapstructure:"system" json:"system" yaml:"system"`
	Deli     Deli     `mapstructure:"deli" yaml:"deli" json:"deli"`
	Wxpusher Wxpusher `napstructure:"wxpusher" json:"wxpusher" yaml:"wxpusher"`
}
