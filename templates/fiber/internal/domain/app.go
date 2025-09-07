package domain

type AppConfigs struct {
	App AppConfig `mapstructure:"app" yaml:"app"`
}

type AppConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port string `mapstructure:"port" yaml:"port"`
}
