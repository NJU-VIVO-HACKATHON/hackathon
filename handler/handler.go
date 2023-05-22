package handler

import (
	"fmt"
	"github.com/NJU-VIVO-HACKATHON/hackathon/global"
	"github.com/NJU-VIVO-HACKATHON/hackathon/model"
	"github.com/NJU-VIVO-HACKATHON/hackathon/repository"
	"github.com/gin-gonic/gin"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
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
type TagInfo struct {
	Tid   int64  `json:"tid"`
	Name  string `json:"name"`
	Cover string `json:"cover"`
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
		_, _, err = repository.CreateUser(session.Auth.Email, session.Auth.Sms, db)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
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

	pid, _, err := repository.CreatePost(model.Post{
		Uid: func() *int64 {
			a := int64(uid.(int))
			return &a
		}(),
		Content: postInfo.Content,
		Title:   postInfo.Title,
		Cover:   postInfo.Cover,
	}, db)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"pid": int(pid)})

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
	pageSize, pageNum := GetPageInfo(c)
	posts, err := repository.SearchPosts(pageNum, pageSize, &keyword, db)
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
		Pid:           int64(post.ID),
		Title:         post.Title,
		Cover:         post.Cover,
		Nickname:      user.Nickname,
		Avatar:        user.Avatar,
		Content:       post.Content,
		IsLike:        false,
		LikeCount:     post.LikeCount,
		FavoriteCount: post.FavoriteCount,
	})

}

// PostBookmark 点赞收藏
func PostBookmark(c *gin.Context) {
	db, _ := repository.GetDataBase()
	uid, isExit := c.Get("uid")
	if !isExit {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", "uid is not exit"))
		return
	}

	_typeStr := c.Param("type")
	pidStr := c.Param("pid")
	_type, err := strconv.ParseInt(_typeStr, 10, 64)
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	err = repository.PostBookMark(int64(uid.(int)), pid, _type, db)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

// InitLoveTags 初始化喜欢的标签
func InitLoveTags(c *gin.Context) {
	var tids []int64
	if err := c.BindJSON(&tids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, _ := repository.GetDataBase()
	uid, isExit := c.Get("uid")
	if !isExit {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", "uid is not exit"))
		return
	}

	var loveTags []*model.LoveTag
	for _, tid := range tids {
		loveTags = append(loveTags, &model.LoveTag{Uid: uid.(int64), Tid: tid, Level: 10})
	}
	repository.InitLoveTags(loveTags, db)
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	db, _ := repository.GetDataBase()
	var tag model.Tag
	if err := c.BindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, _, err := repository.CreateTag(tag, db)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

// SearchTags 搜索标签
func SearchTags(c *gin.Context) {
	db, _ := repository.GetDataBase()
	keyword := c.Query("q")
	pageSize, pageNum := GetPageInfo(c)
	tags, err := repository.SearchTags(&keyword, pageNum, pageSize, db)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	var result []TagInfo
	for _, tag := range tags {
		result = append(result, TagInfo{
			Tid:   int64(tag.ID),
			Name:  tag.Name,
			Cover: tag.Cover,
		})
	}

	c.IndentedJSON(http.StatusOK, result)
}

// GetAllTags 获取所有基础标签
func GetAllTags(c *gin.Context) {
	db, _ := repository.GetDataBase()
	tags, err := repository.GetTags(db)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	var result []TagInfo
	for _, tag := range tags {
		result = append(result, TagInfo{
			Tid:   int64(tag.ID),
			Name:  tag.Name,
			Cover: tag.Cover,
		})
	}

	c.IndentedJSON(http.StatusOK, result)
}

// GetHistory 获取历史记录
func GetHistory(c *gin.Context) {
	_type := c.Param("type")
	uid, isExist := c.Get("uid")

	if !isExist {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", "uid is not exit"))
		return
	}
	pageSize, pageNum := GetPageInfo(c)
	db, _ := repository.GetDataBase()
	var postInfos []PostInfo

	if _type == "my" {
		posts, error := repository.GetPostsByUid(pageNum, pageSize, uid.(int64), db)
		if error != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", error.Error()))
			return
		}
		user, error := repository.GetUserInfo(uid.(int64), db)
		if error != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", error.Error()))
			return
		}
		for _, post := range posts {
			postInfos = append(postInfos, PostInfo{
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

		c.IndentedJSON(http.StatusOK, postInfos)
		return
	}

	users, err := repository.GetAllUsers(db)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	var result []PostInfo
	var posts []model.Post

	posts, err = repository.GetBookMarkUserLike(pageNum, pageSize, uid.(int64), _type, db)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err))
		return
	}

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

// LocalPosts todo 未实现
func LocalPosts(c *gin.Context) {}

func GetComments(c *gin.Context) {}

// UploadFiles 上传文件
func UploadFiles(c *gin.Context) {
	fieldName := c.DefaultPostForm("fieldName", "files")
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("Fail to get file from FormFile by fieldName", err)
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File[fieldName]
	db, _ := repository.GetDataBase()
	for _, file := range files {
		dst := "../target/upload/multiple/" + file.Filename
		// Save the uploaded file to the specified directory
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))

			return
		}

		userAgent := c.GetHeader("User-Agent")
		fileType := path.Ext(file.Filename)
		_, _, err = repository.InsertFileLog("../target/upload/", file.Filename, userAgent, fileType, file.Size, db)

	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))

}

// DownloadFile 下载文件
func DownloadFile(c *gin.Context) {

	fileName := c.Param("uuid")
	baseUrl := path.Join("..", "target", "upload")

	filePath := path.Join(baseUrl, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "file not found"})
		return
	}
	// Get ext
	ext := path.Ext(filePath)
	// Set response Header
	c.Header("Content-Type", mime.TypeByExtension(ext))
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Status(http.StatusOK)
	c.File(filePath)

}

// GetPostTags 获取文章标签
func GetPostTags(c *gin.Context) {
	db, _ := repository.GetDataBase()
	var tags []*model.Tag

	pidStr := c.Param("pid")
	pid, _ := strconv.Atoi(pidStr)

	tags, err := repository.GetTagByPid(int64(pid), db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, tags)
}

// PostPostTags 添加文章标签
func PostPostTags(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, _ := strconv.Atoi(pidStr)
	var tids []*int64

	if err := c.BindJSON(&tids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, _ := repository.GetDataBase()

	err := repository.PostPostTags(int64(pid), tids, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)

}
