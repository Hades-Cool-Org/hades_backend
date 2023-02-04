package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"hades_backend/app/logger"
	"os"
)

var Cfg *Config

func init() {
	Cfg = newConfig()
}

type Config struct {
	Server struct {
		Port string `yaml:"port", envconfig:"SERVER_PORT"`
		Host string `yaml:"host", envconfig:"SERVER_HOST"`
	} `yaml:"server"`
	Database struct {
		Username string `yaml:"user", envconfig:"DB_USERNAME"`
		Host     string `yaml:"host", envconfig:"DB_HOST"`
		Port     string `yaml:"port", envconfig:"DB_PORT"`
		Password string `yaml:"pass", envconfig:"DB_PASSWORD"`
		DbName   string `yaml:"db_name", envconfig:"DB_NAME"`
	} `yaml:"database"`
}

func newConfig() *Config {
	logger.Logger.Info("Loading config")
	var cfg Config
	readFile(&cfg)
	readEnv(&cfg)
	return &cfg
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *Config) {
	f, err := os.Open("config.yaml") // WINDOWS
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
