package handler

import (
	m_logger "github.com/NJU-VIVO-HACKATHON/hackathon/m-logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
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

	logger, err, closeFunc := m_logger.InitZapLogger("hackathon_"+time.Now().Format("20060102")+".log", "[Session]")
	if err != nil {
		logger.Error("Failed to init zap logger", zap.Error(err))
	}
	defer closeFunc()

	var session SessionInfo
	if err := c.BindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//todo 连接数据库判断用户是否存在，不存在就注册，存在就校验登陆

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
