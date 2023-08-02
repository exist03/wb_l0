package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
	"wb_l0/pkg/logger"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type"`
		BindIP string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"listen"`
	PsqlStorage `yaml:"psqlStorage"`
}

type PsqlStorage struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

var instance *Config
var once sync.Once

func GetConfigYml() *Config {
	once.Do(func() {
		log := logger.GetLogger()
		log.Info().Msg("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Info().Msg(help)
			log.Err(err)
		}
	})
	return instance
}
