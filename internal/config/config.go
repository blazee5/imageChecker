package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type Config struct {
	Env        string        `env:"ENV"      env-default:"local"  yaml:"env"`
	Timeout    time.Duration `env:"TIMEOUT"  env-default:"4s"     yaml:"timeout"`
	HTTPServer `yaml:"http_server"`
	Redis      `yaml:"redis"`
}

type HTTPServer struct {
	Port string `env:"PORT" env-default:"3000"      yaml:"port"`
}

type Redis struct {
	Host     string `env:"REDIS_HOST" env-default:"localhost" yaml:"host"`
	Port     string `env:"REDIS_PORT" env-default:"6379"      yaml:"port"`
	Password string `env:"REDIS_PASSWORD" env-default:""      yaml:"password"`
}

func LoadConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		log.Fatalf("error while reading config file: %s", err)
	}

	return &cfg
}
