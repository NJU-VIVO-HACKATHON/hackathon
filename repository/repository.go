package repository

import (
	"github.com/NJU-VIVO-HACKATHON/hackathon/config"
	"github.com/NJU-VIVO-HACKATHON/hackathon/model"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strconv"
)

func GetDataBase() (db *gorm.DB, err error) {

	dbConfig := config.GetConfig()

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

	db, err = gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{

		Logger: logger.Default.LogMode(logger.Info), // 控制日志级别为 info
	})

	if err != nil {
		log.Println("Failed to open database", logCfg.FormatDSN())
		return nil, err
	}
	log.Println("Get database success!", logCfg.FormatDSN())
	return db, err

}

// CreateUser 创建用户
func CreateUser(email, sms string, db *gorm.DB) (ID uint, RowsAffected int64, err error) {

	user := model.User{
		Email: email,
		Sms:   sms,
	}
	result := db.Create(&user)
	if result.Error == nil {
		log.Println("Fail to create file log in repository", result.Error)
	}
	log.Println("Insert file log success!", user.ID, "row affected", result.RowsAffected)
	return user.ID, result.RowsAffected, result.Error
}

// GetUserInfo 获取用户信息
func GetUserInfo(id int64, db *gorm.DB) (model.User, error) {

	var user model.User
	result := db.First(&user, id)
	if result.Error != nil {
		log.Println("File to select user", result.Error)
		return user, result.Error
	}

	log.Println("Select user success!", "user.UID", user.ID)
	return user, result.Error

}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(id int64, newUser model.User, db *gorm.DB) (model.User, error) {

	logger.Default.LogMode(logger.Info) // 控制日志级别为 info
	var user model.User
	result := db.First(&user, id)

	if result.Error != nil {
		log.Println("File to select user", result.Error)
		return user, result.Error
	}
	log.Println(user)
	db.Model(&user).Updates(newUser)
	log.Println(user)
	log.Println("Update user success!", "user.UID", user.ID)
	return user, result.Error

}
