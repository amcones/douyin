package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Mysql struct {
	Username  string
	Password  string
	Net       string
	Addr      string
	DbName    string
	Charset   string
	ParseTime bool `toml:"parse_time"`
	Loc       string
}

type Config struct {
	DB Mysql `toml:"mysql"`
}

type Jwt struct {
	SecretKey string
}

type CommonConfig struct {
	JWT Jwt `toml:"JWT"`
}

var Conf Config

var CommonConf CommonConfig

func init() {
	if _, err := toml.DecodeFile("./config/conf.toml", &Conf); err != nil {
		panic(fmt.Errorf("cannot decode config file: %s", err))
	}
	if _, err := toml.DecodeFile("./config/common_conf.toml", &CommonConf); err != nil {
		panic(fmt.Errorf("cannot decode config file: %s", err))
	}
}
