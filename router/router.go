package router

import (
	"fmt"
	"github.com/NJU-VIVO-HACKATHON/hackathon/global"
	"github.com/NJU-VIVO-HACKATHON/hackathon/handler"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// JwtAuthMiddleware 添加解析JWT的中间件
func JwtAuthMiddleware(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}
	// 去除Bearer
	if strings.HasPrefix(authorizationHeader, "Bearer ") {
		authorizationHeader = authorizationHeader[7:]
	}
	log.Println(authorizationHeader)
	claims, err := global.GetJwt().ParseToken(authorizationHeader)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		c.Abort()
		return
	}
	c.Set("uid", claims.Uid)
}

// GetPageInfo 分页
func GetPageInfo(c *gin.Context) (pageId int, pageSize int) {
	var err error
	pid := c.Query("pageId")
	ps := c.Query("pageSize")
	pageId, err = strconv.Atoi(pid)
	if err != nil {
		pageId = 0
	}
	pageSize, err = strconv.Atoi(ps)
	if err != nil {
		pageSize = 10
	}
	return
}

func SetupRouter() *gin.Engine {

	r := gin.Default()
	sessionGroup := r.Group("/session") // 用户鉴权
	{
		sessionGroup.POST("", handler.Session)           // 登陆/注册
		sessionGroup.POST("/authcode", handler.Authcode) // 验证码
	}
	userGroup := r.Group("/users/:uid") // 个人信息
	{
		userGroup.GET("/info", handler.GetUserInfo)         // 获取个人信息
		userGroup.POST("/info", handler.UpdateUserInfo)     // 更新个人信息
		userGroup.GET("/my_tags", handler.GetMyTags)        // 更新用户初始感兴趣的标签
		userGroup.GET("/history/:type", handler.GetHistory) // 获取足迹
	}
	r.GET("/tags", handler.GetAllTags) // 获取所有基本标签

	postGroup := r.Group("/posts", JwtAuthMiddleware) // 首页&帖子模块
	{
		postGroup.GET("", handler.GetPosts)           // 列举帖子
		postGroup.GET("/search", handler.SearchPosts) // 搜索帖子
		postGroup.GET("/local", handler.LocalPosts)   // 附近的帖子

		postGroup.POST("", handler.CreatePosts) //创建帖子

		postInfoGroup := r.Group("/:pid")
		{
			postInfoGroup.GET("", handler.GetPostContext)       // 获取帖子内容
			postInfoGroup.GET("/comments", handler.GetComments) // 获取评论
			postGroup.PUT("", handler.EditPosts)                // 编辑帖子
			postGroup.DELETE("", handler.DelPosts)              // 删除帖子
		}

		postGroup.POST(" /posts/:pid/bookmark/:type", handler.PostBookmark) //点赞/收藏 / 取消点赞/收藏  帖子
	}

	//上传附件
	r.POST("/attachment", handler.Attachment)

	return r
}
