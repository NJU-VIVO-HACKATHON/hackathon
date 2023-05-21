package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

type Database struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db-name"`
}

type Config struct {
	Database Database `yaml:"database"`
}

var (
	config     Config
	configOnce sync.Once
)

func GetConfig() *Config {

	configOnce.Do(func() {
		log.Println("加载配置文件 config.yaml")
		bs, err := os.ReadFile(os.Getenv("CONFIG_PATH"))
		if err != nil {
			log.Fatalln(err)
		}
		_ = yaml.Unmarshal(bs, &config)
	})

	return &config
}
