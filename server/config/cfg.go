package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Database struct {
		Addr       string `yaml:"host"`
		Port       string `yaml:"port"`
		Username   string `yaml:"user"`
		Password   string `yaml:"pass"`
		DBname     string `yaml:"dbname"`
		DriverName string `yaml:"driverName"`
	} `yaml:"database"`
}

func (cfg *Config) InitFile() {
	f, err := os.Open("config/config.yml")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
