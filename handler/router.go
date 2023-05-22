package handler

import (
	"fmt"
	"github.com/NJU-VIVO-HACKATHON/hackathon/global"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

func SetupRouter() *gin.Engine {

	r := gin.Default()
	sessionGroup := r.Group("/session") // 用户鉴权
	{
		sessionGroup.POST("", Session)           // 登陆/注册
		sessionGroup.POST("/authcode", Authcode) // 验证码
	}
	userGroup := r.Group("/users/:uid", JwtAuthMiddleware) // 个人信息
	{
		userGroup.GET("/info", GetUserInfo)         // 获取个人信息
		userGroup.POST("/info", UpdateUserInfo)     // 更新个人信息
		userGroup.POST("/my_tags", InitLoveTags)    // 更新用户初始感兴趣的标签
		userGroup.GET("/history/:type", GetHistory) // 获取足迹
	}

	postGroup := r.Group("/posts", JwtAuthMiddleware) // 首页&帖子模块
	{
		postGroup.GET("", GetPosts)           // 列举帖子
		postGroup.GET("/search", SearchPosts) // 搜索帖子
		postGroup.GET("/local", LocalPosts)   // 附近的帖子

		postGroup.POST("", CreatePosts) //创建帖子

		postInfoGroup := postGroup.Group("/:pid")
		{
			postInfoGroup.GET("", GetPostContext)       // 获取帖子内容
			postInfoGroup.GET("/comments", GetComments) // 获取评论
			postInfoGroup.PUT("", EditPosts)            // 编辑帖子
			postInfoGroup.DELETE("", DelPosts)          // 删除帖子
			postGroup.POST("/tags", PostPostTags)       // 修改文章标注标签
			postGroup.POST("/tags", GetPostTags)        // 获取文章标注标签
		}

		postGroup.POST("/:pid/bookmark/:type", PostBookmark) //点赞/收藏 / 取消点赞/收藏  帖子
	}

	tagGroup := r.Group("tags", JwtAuthMiddleware)
	{
		tagGroup.POST("", CreateTag) //创建标签
		tagGroup.GET("", SearchTags) //搜索标签
	}
	r.GET("/basic_tags", GetAllTags) //获取所有基础标签

	//上传附件
	r.POST("/attachment", UploadFiles)
	r.GET("/attachment/:uuid", DownloadFile)

	return r
}
