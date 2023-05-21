package router

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"

	"testing"
)

import "encoding/json"

func TestSessionLogin(t *testing.T) {
	r := SetupRouter()

	// 组成请求体
	requestBody := gin.H{
		"mode": "email",
		"auth": gin.H{
			"email": "i@zhangzqs.cn",
			"code":  "<code>",
		},
	}

	// 将请求体转换为JSON字符串
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		t.Error(err)
	}

	// 发送POST请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/session", bytes.NewBuffer(requestBodyJSON))
	r.ServeHTTP(w, req)

	// 检查返回结果是否符合预期
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestSessionRegister(t *testing.T) {
	r := SetupRouter()

	// 组成请求体
	requestBody := gin.H{
		"mode": "email",
		"auth": gin.H{
			"email": "i@zhangzqs.cn",
			"code":  "<code>",
		},
	}

	// 将请求体转换为JSON字符串
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		t.Error(err)
	}

	// 发送POST请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/session", bytes.NewBuffer(requestBodyJSON))
	r.ServeHTTP(w, req)

	// 检查返回结果是否符合预期
	assert.Equal(t, http.StatusOK, w.Code)

}
