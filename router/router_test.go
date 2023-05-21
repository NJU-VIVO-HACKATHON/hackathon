package router

import (
	"bytes"
	"fmt"
	"github.com/NJU-VIVO-HACKATHON/hackathon/global"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

import "encoding/json"

const (
	testEmail = "xxx@email.com"
	testSms   = "1234567890"
)

func getUserToken(t *testing.T, r http.Handler, mode string) (token string, uid int) {
	// 组成请求体
	var requestBody gin.H
	if mode == "email" {
		requestBody = gin.H{
			"mode": mode,
			"auth": gin.H{
				"email": testSms,
				"code":  "23456",
			},
		}
	} else {
		requestBody = gin.H{
			"mode": "email",
			"auth": gin.H{
				"email": testEmail,
				"code":  "12345",
			},
		}
	}

	// 将请求体转换为JSON字符串
	requestBodyJSON, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	// 发送POST请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/session", bytes.NewBuffer(requestBodyJSON))
	r.ServeHTTP(w, req)

	var responseBody struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	assert.NotEmpty(t, responseBody.Token)

	token = responseBody.Token
	payload, err := global.GetJwt().ParseToken(token)
	assert.NoError(t, err)
	uid = payload.Uid
	return
}
func TestSessionLoginByEmail(t *testing.T) {
	r := SetupRouter()
	getUserToken(t, r, "email")
}

func TestSessionLoginBySms(t *testing.T) {
	r := SetupRouter()
	getUserToken(t, r, "sms")
}

func TestSessionAuthcode(t *testing.T) {
	r := SetupRouter()

	// 组成请求体
	requestBody := gin.H{
		"mode": "sms",
		"auth": gin.H{
			"sms": testSms,
		},
	}

	// 将请求体转换为JSON字符串
	requestBodyJSON, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	// 发送POST请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/session/authcode", bytes.NewBuffer(requestBodyJSON))
	r.ServeHTTP(w, req)

	// 检查返回结果是否符合预期
	assert.Equal(t, http.StatusOK, w.Code)
}

func updateUserInfo(t *testing.T, r http.Handler, info gin.H) {
	token, uid := getUserToken(t, r, "email")
	assert.NotEmpty(t, token)

	requestBodyJson, err := json.Marshal(info)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/users/%d/info", uid), bytes.NewBuffer(requestBodyJson))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	// 检查返回结果是否符合预期
	assert.Equal(t, http.StatusOK, w.Code)
}

func getUserInfo(t *testing.T, r http.Handler) gin.H {
	token, uid := getUserToken(t, r, "email")
	assert.NotEmpty(t, token)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d/info", uid), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	// 检查返回结果是否符合预期
	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody gin.H
	bs := w.Body.Bytes()
	err := json.Unmarshal(bs, &responseBody)
	assert.NoError(t, err)

	return responseBody
}

func getRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func TestSetAndGetUserInfo(t *testing.T) {
	r := SetupRouter()

	expectInfo := gin.H{
		"nickname":     "testNickname" + getRandomString(),
		"avatar":       "testAvatar" + getRandomString(),
		"introduction": "testIntroduction" + getRandomString(),
	}

	updateUserInfo(t, r, expectInfo)
	actualUserInfo := getUserInfo(t, r)
	assert.Equal(t, expectInfo["nickname"], actualUserInfo["nickname"])
	assert.Equal(t, expectInfo["avatar"], actualUserInfo["avatar"])
	assert.Equal(t, expectInfo["introduction"], actualUserInfo["introduction"])
}

func addTag(t *testing.T, r http.Handler, tag string) {
	token, uid := getUserToken(t, r, "email")
	assert.NotEmpty(t, token)

	requestBodyJson, err := json.Marshal(gin.H{
		"tag": tag,
	})
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/users/%d/tags", uid), bytes.NewBuffer(requestBodyJson))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	// 检查返回结果是否符合预期
	assert.Equal(t, http.StatusOK, w.Code)
}

func addPost(t *testing.T, r http.Handler, title, content, cover string) (pid int) {
	token, _ := getUserToken(t, r, "email")
	assert.NotEmpty(t, token)

	requestBodyJson, err := json.Marshal(gin.H{
		"title":   title,
		"content": content,
		"cover":   cover,
	})
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(requestBodyJson))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	// 检查返回结果是否符合预期
	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody gin.H
	bs := w.Body.Bytes()
	err = json.Unmarshal(bs, &responseBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseBody["pid"])

	return int(responseBody["pid"].(float64))
}

func getPost(t *testing.T, r http.Handler, pid int) gin.H {
	token, _ := getUserToken(t, r, "email")
	assert.NotEmpty(t, token)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/posts/%d", pid), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	// 检查返回结果是否符合预期
	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody gin.H
	bs := w.Body.Bytes()
	err := json.Unmarshal(bs, &responseBody)
	assert.NoError(t, err)
	return responseBody
}

func TestCreatePostAndGet(t *testing.T) {
	r := SetupRouter()
	expectTitle := "testTitle" + getRandomString()
	expectContent := "testContent" + getRandomString()
	expectCover := "testCover" + getRandomString()

	pid := addPost(t, r, expectTitle, expectContent, expectCover)
	log.Println("pid:", pid)
	actualPost := getPost(t, r, pid)
	assert.Equal(t, expectTitle, actualPost["title"])
	assert.Equal(t, expectContent, actualPost["content"])
	assert.Equal(t, expectCover, actualPost["cover"])
}
