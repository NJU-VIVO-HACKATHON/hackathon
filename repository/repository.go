package repository

import (
	"github.com/NJU-VIVO-HACKATHON/hackathon/config"
	m_logger "github.com/NJU-VIVO-HACKATHON/hackathon/m-logger"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func GetDataBase() (db *gorm.DB, err error) {

	dbConfig := config.GetConfig()

	//set database  logger
	logger, err, closeFunc := m_logger.InitLogFile("hackathon_"+time.Now().Format("20060102")+".log", "[GetDataBase]")
	if err != nil {
		logger.Println("Failed to init logger", err)
	}
	defer closeFunc()

	// Capture connection properties.
	cfg := mysqlDriver.Config{
		User:   dbConfig.Database.Username,
		Passwd: dbConfig.Database.Password,
		DBName: dbConfig.Database.DbName,
		Addr:   dbConfig.Database.Hostname + ":" + strconv.Itoa(dbConfig.Database.Port),
		Net:    "tcp",
		Params: map[string]string{
			"loc":       "Local",
			"parseTime": "True",
		},
	}

	//log cfg
	logCfg := mysqlDriver.Config{
		User:   "user",
		Passwd: "password",
		DBName: dbConfig.Database.DbName,
		Addr:   "localhost:8080",
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
