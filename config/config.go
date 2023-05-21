package config

import (
	m_logger "github.com/NJU-VIVO-HACKATHON/hackathon/m-logger"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
	"time"
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
	//set logger
	logger, err, closeFunc := m_logger.InitLogFile("hackathon_"+time.Now().Format("20060102")+".log", "[GetConfig]")
	if err != nil {
		logger.Println("Failed to init logger", err)
	}
	defer closeFunc()

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
