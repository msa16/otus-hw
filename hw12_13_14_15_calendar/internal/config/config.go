package config

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
	Kafka   KafkaConf  `yaml:"kafka"`
	Storage string
	DB      StorageConf
	Timer   TimerConf
}

type LoggerConf struct {
	Level string
	File  string
}

type ServerConf struct {
	HTTP AddrConf `yaml:"http"`
}

type KafkaConf struct {
	Host  string
	Port  int
	Topic string
}

type AddrConf struct {
	Host string
	Port int
}

type StorageConf struct {
	Driver string
	Dsn    string
}

type TimerConf struct {
	ReminderEvents int `yaml:"reminder_events"`
	OldEvents      int `yaml:"old_events"`
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
