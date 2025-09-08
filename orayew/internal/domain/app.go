package domain

import "orayew/pkg/pgxpool"

type AppConfigs struct {
	App	AppConfig	`mapstructure:"app" yaml:"app"`
	DB	pgxpool.Config	`mapstructure:"pgxpool" yaml:"pgxpool"`
}

type AppConfig struct {
	Host	string	`mapstructure:"host" yaml:"host"`
	Port	string	`mapstructure:"port" yaml:"port"`
}
