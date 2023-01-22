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

type Redis struct {
	Net      string
	Addr     string
	Password string
}

type COS struct {
	Addr      string
	SecretId  string
	SecretKey string
}

type Config struct {
	DB    Mysql `toml:"mysql"`
	JWT   Jwt   `toml:"jwt"`
	Redis Redis `toml:"redis"`
	COS   COS   `toml:"cos"`
}

type Jwt struct {
	SecretKey string
}

var Conf Config

func init() {
	if _, err := toml.DecodeFile("./config/conf.toml", &Conf); err != nil {
		panic(fmt.Errorf("cannot decode config file: %s", err))
	}
}
