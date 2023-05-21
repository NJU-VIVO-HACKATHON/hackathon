package repository

import (
	"errors"
	"github.com/NJU-VIVO-HACKATHON/hackathon/global"
	"github.com/NJU-VIVO-HACKATHON/hackathon/model"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strconv"
)

func GetDataBase() (db *gorm.DB, err error) {

	dbConfig := global.GetConfig()

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
func CreateUser(email, sms *string, db *gorm.DB) (ID uint, RowsAffected int64, err error) {

	user := model.User{
		Email: email,
		Sms:   sms,
	}
	result := db.Create(&user)
	if result.Error != nil {
		log.Println("Fail to create user in database", result.Error)
		return 0, 0, result.Error
	}
	log.Println("Create user in database success!", user.ID, "row affected", result.RowsAffected)
	return user.ID, result.RowsAffected, nil
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

// GetUserInfoByAuth 通过Em获取用户信息
func GetUserInfoByAuth(mode, email, sms *string, db *gorm.DB) (*model.User, error) {

	var user model.User
	var result *gorm.DB

	if *mode == "email" {
		result = db.Where("email=?", email).First(&user)
	} else if *mode == "sms" {
		result = db.Where("sms=?", sms).First(&user)
	} else {
		return nil, errors.New("非法输入")
	}

	if result.Error != nil {
		log.Println("File to select user", result.Error)
		return &user, result.Error
	}

	log.Println("Select user success!", "user.UID", user.ID)
	return &user, nil
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(id int64, newUser model.User, db *gorm.DB) error {

	var user model.User
	result := db.First(&user, id)
	if result.Error != nil {
		log.Println("File to select user", result.Error)
		return result.Error
	}

	result = db.Model(&user).Updates(newUser)
	if result.Error != nil {
		log.Println("File to update user", result.Error)
		return result.Error
	}

	log.Println("Update user success!", "user.UID", user.ID)
	return result.Error

}

// CreatePost 创建帖子
func CreatePost(uidStr, title, content *string, db *gorm.DB) (ID uint, RowsAffected int64, err error) {
	uid, err := strconv.ParseInt(*uidStr, 10, 64)
	post := model.Post{
		Uid:     &uid,
		Title:   title,
		Content: content,
	}
	result := db.Create(&post)
	if result.Error != nil {
		log.Println("Fail to post user in database", result.Error)
		return 0, 0, result.Error
	}
	log.Println("Create post in database success!", post.ID, "row affected", result.RowsAffected)
	return post.ID, result.RowsAffected, nil
}

// EditPost 编辑帖子
func EditPost(pid *string, newPost *model.Post, db *gorm.DB) error {
	var post model.Post
	result := db.First(&post, pid)
	if result.Error != nil {
		log.Println("File to select post", result.Error)
		return result.Error
	}

	result = db.Model(&post).Updates(newPost)
	if result.Error != nil {
		log.Println("File to update post", result.Error)
		return result.Error
	}

	return result.Error
}

// DeletePost 删除帖子
func DeletePost(pid *string, db *gorm.DB) error {
	result := db.Delete(&model.Post{}, pid)
	if result.Error != nil {
		log.Println("File to delete post", result.Error)
		return result.Error
	}
	return result.Error
}

// GetPosts 列举帖子
func GetPosts(tagId *string, db *gorm.DB) ([]*model.Post, error) {
	var posts []*model.Post
	var postTags []*model.PostTag
	var result *gorm.DB
	if tagId != nil {
		pids := db.Where("tagId=?", tagId).Find(&postTags).Select("pid")
		//时间倒序 新帖在前
		result = db.Where("id IN (?)", pids).Order("updated_at desc").Find(&posts)
	}

	result = db.Order("updated_at desc").Find(&posts)

	if result.Error != nil {
		log.Println("File to select posts", result.Error)
		return posts, result.Error
	}
	return posts, nil
}

// GetPost 获取帖子
func GetPost(pid *string, db *gorm.DB) (*model.Post, error) {
	var post model.Post
	result := db.First(&post, pid)
	if result.Error != nil {
		log.Println("File to select post", result.Error)
		return &post, result.Error
	}
	return &post, nil
}

// SearchPost 搜索帖子
func SearchPost(pageNum, pageSize, keyword *string, db *gorm.DB) ([]*model.Post, error) {
	limit, _ := strconv.Atoi(*pageSize)
	offset, _ := strconv.Atoi(*pageNum)
	offset = offset * limit
	var posts []*model.Post
	result := db.Where("title LIKE ?", "%"+*keyword+"%").Limit(limit).Offset(offset).Find(&posts)
	if result.Error != nil {
		log.Println("File to select posts", result.Error)
		return posts, result.Error
	}
	return posts, nil
}

// GetComments todo  获取评论
func GetComments(pid, pageNum, pageSize *string, db *gorm.DB) {}
