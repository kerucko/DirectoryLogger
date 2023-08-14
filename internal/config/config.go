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
	Path          string   `yaml:"path"`
	IncludeRegexp []string `yaml:"include_regexp"`
	ExcludeRegexp []string `yaml:"exclude_regexp"`
}

var C Config

func ReadConfig() {
	log.Println("read config")
	if err := cleanenv.ReadConfig("config.yml", &C); err != nil {
		log.Printf("ERROR ReadConfig")
		panic(err)
	}
	log.Println("CONFIG:", C)
}
