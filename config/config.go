package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Db_host     string `required:"true"`
	Db_port     string `required:"true" default:"3306"`
	Db_user     string `required:"true"`
	Db_password string `required:"true"`
	Db_name     string `required:"true"`
}

var Cfg Config

func LoadConfig() error {
	log.Println("LoadConfig...")
	err := envconfig.Process("", &Cfg)

	return err
}
