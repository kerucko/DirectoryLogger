package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Storage     DBConfig    `yaml:"storage"`
	Directories []DirConfig `yaml:"directories"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type DirConfig struct {
	Path   string   `yaml:"path"`
	Regexp []string `yaml:"regexp"`
}

var C Config

func ReadConfig() {
	if err := cleanenv.ReadConfig("config.yml", &C); err != nil {
		help, _ := cleanenv.GetDescription(C, nil)
		log.Fatalf("ERROR ReadConfig: %v/n%v", err, help)
	}
}
