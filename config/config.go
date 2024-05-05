package config

import (
	"log"
	"reflect"

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

	cfgKeyList := reflect.ValueOf(Cfg)
	for i := 0; i < cfgKeyList.NumField(); i++ {
		log.Printf("%s: %v\n", cfgKeyList.Type().Field(i).Name, cfgKeyList.Field(i).String())
	}

	return err
}
