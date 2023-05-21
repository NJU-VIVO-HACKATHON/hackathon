package repository

import (
	m_logger "github.com/NJU-VIVO-HACKATHON/hackathon/m-logger"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func GetDataBase() (db *gorm.DB, err error) {

	//set database zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server"+".log", "[GetDataBase]")
	if err != nil {
		logger.Error("Failed to init zap logger", zap.Error(err))
	}
	defer closeFunc()

	// Capture connection properties.
	cfg := mysqlDriver.Config{
		User:   os.Getenv("DBUser"),
		Passwd: os.Getenv("DBPassword"),
		DBName: os.Getenv("DBName"),
		Addr:   "localhost:3306",
		Net:    "tcp",
		Params: map[string]string{
			"loc":       "Local",
			"parseTime": "True",
		},
	}
	//log cfg
	logCfg := mysqlDriver.Config{
		User:   "DBUser",
		Passwd: "DBPassword",
		DBName: os.Getenv("DBName"),
		Addr:   "localhost:3306",
		Net:    "tcp",
		Params: map[string]string{
			"loc":       "Local",
			"parseTime": "True",
		},
	}
	db, err = gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to open database", zap.Error(err))
		return nil, err
	}

	logger.Info("Get database success!", zap.String("logCfg", logCfg.FormatDSN()))

	return db, err
}
