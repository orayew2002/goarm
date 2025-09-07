package viper

import (
	"fmt"

	"github.com/spf13/viper"
)

func Parse(path string, cfg interface{}) error {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %s ", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("fatal error config file: %s ", err)
	}

	return nil
}
