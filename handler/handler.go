package handler

import (
	"fmt"
	"github.com/NJU-VIVO-HACKATHON/hackathon/global"
	"github.com/NJU-VIVO-HACKATHON/hackathon/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type SessionInfo struct {
	Mode string `json:"mode"`
	Auth struct {
		Email *string `json:"email,omitempty"`
		Code  *string `json:"code"`
		Sms   *string `json:"sms,omitempty"`
	} `json:"auth"`
}

func Session(c *gin.Context) {

	db, _ := repository.GetDataBase()

	var session SessionInfo
	if err := c.BindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 连接数据库判断用户是否存在，不存在就注册，存在就校验登陆
	user, err := repository.GetUserInfoByAuth(session.Auth.Email, session.Auth.Sms, db)
	//登陆
	if err == nil {
		c.IndentedJSON(http.StatusOK, gin.H{"token": global.GetJwt().GenerateToken(int(user.ID))})
	} else if err.Error() == "非法输入" {
		log.Print(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
	} else {
		repository.CreateUser(session.Auth.Email, session.Auth.Sms, db)
		c.IndentedJSON(http.StatusOK, gin.H{"token": global.GetJwt().GenerateToken(int(user.ID))})
	}

}

func Authcode(c *gin.Context)       {}
func GetUserInfo(c *gin.Context)    {}
func UpdateUserInfo(c *gin.Context) {}
func GetMyTags(c *gin.Context)      {}
func GetHistory(c *gin.Context)     {}
func GetAllTags(c *gin.Context)     {}
func GetPosts(c *gin.Context)       {}
func GetPostContext(c *gin.Context) {}
func LocalPosts(c *gin.Context)     {}
func GetComments(c *gin.Context)    {}
func DelPosts(c *gin.Context)       {}
func EditPosts(c *gin.Context)      {}
func Attachment(c *gin.Context)     {}
func SearchPosts(c *gin.Context)    {}
