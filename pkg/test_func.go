package pkg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// Author:Boyn
// Date:2020/9/8

func GetTest(uri string, router *gin.Engine) (Message, error) {
	// 构造get请求
	req := httptest.NewRequest(http.MethodGet, uri, nil)
	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return Message{}, err
	}
	var msg Message
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return Message{}, err
	}
	return msg, nil
}
