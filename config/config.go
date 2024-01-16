package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	TokenExpiration        int32
	RefreshTokenExpiration int32
}

func GetConfig() Config {
	var config Config
	_, err := toml.DecodeFile("../config.toml", &config)
	if err != nil {
		panic(err)
	}
	return config

}
