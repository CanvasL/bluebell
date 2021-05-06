package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `{
		"community_id": 1,
		"title": "test",
		"content": "just a test"
	}`		//三个字段都是required

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// 判断是否按预期反悔了需要登录的错误
	// 方法1：判断响应内容中是否是包含相应的字符串
	assert.Contains(t, w.Body.String(), "需要登录")
	// 方法2：将响应到的内容反序列化到ResponseData，然后判断字段与预期是否一致
	//res := new(ResponseData)
	//if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
	//	t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
	//}
	//assert.Equal(t, res.Code, CodeNeedLogin)
}