package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logging struct {
		Path     string `yaml:"path"`
		Level    string `yaml:"level"`
		RavenDSN string `yaml:"raven-dsn"`
	} `yaml:"logging"`
	DB struct {
		DSN     string `yaml:"dsn"`
		MaxConn int    `yaml:"maxconn"`
		Debug   bool   `yaml:"debug"`
	} `yaml:"db"`
	HTTP struct {
		Prefix string `yaml:"prefix"`
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
	} `yaml:"http"`
}

var (
	config  *Config = &Config{}
	Version string
)

func Init(path string) {
	filename, _ := filepath.Abs(path)
	yamlFile, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, config)

	if err != nil {
		panic(err)
	}
}

func Get() *Config {
	return config
}
