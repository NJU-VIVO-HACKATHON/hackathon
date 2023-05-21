package global

import (
	"github.com/NJU-VIVO-HACKATHON/hackathon/config"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

var (
	configIns  config.Config
	configOnce sync.Once
)

func GetConfig() *config.Config {

	configOnce.Do(func() {
		log.Println("加载配置文件 config.yaml")
		bs, err := os.ReadFile(os.Getenv("CONFIG_PATH"))
		if err != nil {
			log.Fatalln(err)
		}
		_ = yaml.Unmarshal(bs, &configIns)
	})

	return &configIns
}
