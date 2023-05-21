package handler

import (
	"fmt"
	"github.com/NJU-VIVO-HACKATHON/hackathon/global"
	"github.com/NJU-VIVO-HACKATHON/hackathon/model"
	"github.com/NJU-VIVO-HACKATHON/hackathon/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type SessionInfo struct {
	Mode string `json:"mode"`
	Auth struct {
		Email *string `json:"email,omitempty"`
		Code  *string `json:"code"`
		Sms   *string `json:"sms,omitempty"`
	} `json:"auth"`
}
type UserInfo struct {
	Uid          int64   `json:"uid"`
	Email        *string `json:"email"`
	Sms          *string `json:"sms"`
	Nickname     *string `json:"nickname"`
	Avatar       *string `json:"avatar"`
	Introduction *string `json:"introduction"`
}

type PostInfo struct {
	Pid           int64   `json:"pid"`
	Uid           int64   `json:"uid"`
	Content       *string `json:"content"`
	Cover         *string `json:"cover"`
	Title         *string `json:"title"`
	IsLike        bool    `json:"isLike"`
	LikeCount     *int64  `json:"likeCount"`
	FavoriteCount *int64  `json:"isFavorite"`
	Nickname      *string `json:"nickname"`
	Avatar        *string `json:"avatar"`
}

func Session(c *gin.Context) {

	db, _ := repository.GetDataBase()

	var session SessionInfo
	if err := c.BindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 连接数据库判断用户是否存在，不存在就注册，存在就校验登陆
	user, err := repository.GetUserInfoByAuth(&session.Mode, session.Auth.Email, session.Auth.Sms, db)
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

// Authcode todo
func Authcode(c *gin.Context) {}

// GetUserInfo 获取个人信息
func GetUserInfo(c *gin.Context) {
	db, _ := repository.GetDataBase()
	uidStr := c.Param("uid")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	user, err := repository.GetUserInfo(uid, db)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	c.IndentedJSON(http.StatusOK, UserInfo{
		Uid:          int64(user.ID),
		Email:        user.Email,
		Sms:          user.Sms,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Introduction: user.Introduction,
	})

}

// UpdateUserInfo 更新个人信息
func UpdateUserInfo(c *gin.Context) {

	db, _ := repository.GetDataBase()
	var userInfo UserInfo

	if err := c.BindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uidStr := c.Param("uid")
	uid, err := strconv.ParseInt(uidStr, 10, 64)

	err = repository.UpdateUserInfo(uid, model.User{
		Email:        userInfo.Email,
		Sms:          userInfo.Sms,
		Nickname:     userInfo.Nickname,
		Avatar:       userInfo.Avatar,
		Introduction: userInfo.Introduction,
	}, db)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
	} else {
		c.Status(http.StatusOK)
	}

}

// CreatePosts 创建帖子
func CreatePosts(c *gin.Context) {
	db, _ := repository.GetDataBase()
	var postInfo PostInfo

	if err := c.BindJSON(&postInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, isExit := c.Get("uid")

	if !isExit {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", "uid is not exist"))
		return
	}

	_, _, err := repository.CreatePost(model.Post{
		Uid:     uid.(*int64),
		Content: postInfo.Content,
		Title:   postInfo.Title,
		Cover:   postInfo.Cover,
	}, db)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	c.Status(http.StatusOK)

}

// EditPosts 编辑帖子
func EditPosts(c *gin.Context) {
	db, _ := repository.GetDataBase()
	var postInfo PostInfo

	pidStr := c.Param("pid")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	if err := c.BindJSON(&postInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = repository.EditPost(pid, &model.Post{
		Content: postInfo.Content,
		Title:   postInfo.Title,
		Cover:   postInfo.Cover,
	}, db)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	c.Status(http.StatusOK)

}

// DelPosts 删除帖子
func DelPosts(c *gin.Context) {
	db, _ := repository.GetDataBase()
	pidStr := c.Param("pid")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	err = repository.DeletePost(pid, db)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

// GetPosts 列举帖子
func GetPosts(c *gin.Context) {
	db, _ := repository.GetDataBase()
	tagIdStr := c.Query("tag")
	tagId, err := strconv.ParseInt(tagIdStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	posts, err := repository.GetPosts(&tagId, db)
	users, err := repository.GetAllUsers(db)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	// 定义一个保存交集数据的切片
	var result []PostInfo

	// 遍历用户表
	for _, user := range users {
		// 遍历文章表
		for _, post := range posts {
			// 如果文章表中存在对应的用户ID，将该行数据放入新切片
			if int64(user.ID) == *post.Uid {
				result = append(result, PostInfo{
					Pid:      int64(post.ID),
					Title:    post.Title,
					Cover:    post.Cover,
					Nickname: user.Nickname,
					Avatar:   user.Avatar,
					//todo
					IsLike:        false,
					LikeCount:     post.LikeCount,
					FavoriteCount: post.FavoriteCount,
				})
			}
		}
	}
	c.IndentedJSON(http.StatusOK, result)

}

// SearchPosts 搜索帖子
func SearchPosts(c *gin.Context) {
	db, _ := repository.GetDataBase()
	keyword := c.Query("q")
	pageSize := c.Query("pageSize")
	pageNum := c.Query("pageNum")
	posts, err := repository.SearchPosts(&pageNum, &pageSize, &keyword, db)
	users, err := repository.GetAllUsers(db)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	// 定义一个保存交集数据的切片
	var result []PostInfo

	// 遍历用户表
	for _, user := range users {
		// 遍历文章表
		for _, post := range posts {
			// 如果文章表中存在对应的用户ID，将该行数据放入新切片
			if int64(user.ID) == *post.Uid {
				result = append(result, PostInfo{
					Pid:      int64(post.ID),
					Title:    post.Title,
					Cover:    post.Cover,
					Nickname: user.Nickname,
					Avatar:   user.Avatar,
					//todo
					IsLike:        false,
					LikeCount:     post.LikeCount,
					FavoriteCount: post.FavoriteCount,
				})
			}
		}
	}
	c.IndentedJSON(http.StatusOK, result)

}

// GetPostContext 获取帖子内容
func GetPostContext(c *gin.Context) {
	db, _ := repository.GetDataBase()
	pidStr := c.Param("pid")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	post, err := repository.GetPost(pid, db)
	user, err := repository.GetUserInfo(*post.Uid, db)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	c.IndentedJSON(http.StatusOK, PostInfo{
		Pid:      int64(post.ID),
		Title:    post.Title,
		Cover:    post.Cover,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		//todo
		IsLike:        false,
		LikeCount:     post.LikeCount,
		FavoriteCount: post.FavoriteCount,
	})

}

func GetMyTags(c *gin.Context)  {}
func GetHistory(c *gin.Context) {}
func GetAllTags(c *gin.Context) {}

func LocalPosts(c *gin.Context)  {}
func GetComments(c *gin.Context) {}

func Attachment(c *gin.Context) {}
