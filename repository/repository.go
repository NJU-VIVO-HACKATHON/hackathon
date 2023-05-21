package repository

import (
	m_logger "github.com/NJU-VIVO-HACKATHON/hackathon/m-logger"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func GetDataBase() (db *gorm.DB, err error) {

	//set database  logger
	logger, err, closeFunc := m_logger.InitLogFile("hackathon_"+time.Now().Format("20060102")+".log", "[GetDataBase]")
	if err != nil {
		logger.Println("Failed to init logger", err)
	}
	defer closeFunc()

	// Capture connection properties.
	cfg := mysqlDriver.Config{}
	viper.SetConfigFile("./config.json")
	if err := viper.ReadInConfig(); err == nil {
		// 读取数据库连接信息
		cfg = mysqlDriver.Config{
			User:   viper.GetString("database.DBUser"),
			Passwd: viper.GetString("database.DBPassword"),
			DBName: viper.GetString("database.DBName"),
			Addr:   viper.GetString("database.DBAddr"),
			Net:    "tcp",
			Params: map[string]string{
				"loc":       "Local",
				"parseTime": "True",
			},
		}
	}

	//log cfg
	logCfg := mysqlDriver.Config{
		User:   "user",
		Passwd: "password",
		DBName: viper.GetString("database.DBName"),
		Addr:   "localhost:8089",
		Net:    "tcp",
		Params: map[string]string{
			"loc":       "Local",
			"parseTime": "True",
		},
	}

	db, err = gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{})

	if err != nil {
		logger.Println("Failed to open database", logCfg.FormatDSN())
		return nil, err
	}
	logger.Println("Get database success!", logCfg.FormatDSN())
	return db, err
}
