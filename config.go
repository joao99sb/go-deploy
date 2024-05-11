package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type resource struct {
	Port string
}
type configuration struct {
	Server struct {
		Listen_port  string
		Timeout      string
		Default_port string
	}
	Resources []resource
}

var Config *configuration

func GetConfig() (*configuration, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config file: %s", err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}
	return Config, nil
}
