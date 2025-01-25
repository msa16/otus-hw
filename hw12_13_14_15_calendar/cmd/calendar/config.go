package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf `yaml:"log"`
	Server  ServerConf `yaml:"server"`
	Storage string
	Db      StorageConf
}

type LoggerConf struct {
	Level string
	File  string
}

type ServerConf struct {
	Http AddrConf `yaml:"http"`
}

type AddrConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Driver string
	Dsn    string
}

func NewConfig(configFile string) Config {
	config := Config{Logger: LoggerConf{}}
	confFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = yaml.Unmarshal(confFile, &config) //nolint:typecheck
	if err != nil {
		log.Fatal(err.Error())
	}
	return config
}
