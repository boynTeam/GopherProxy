package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

// Author:Boyn
// Date:2020/9/8

type GetParam struct {
	Uri    string
	Header map[string]string
	Query  map[string]string
}

type PostParam struct {
	Uri    string
	Header map[string]string
	Query  map[string]string
	Body   interface{}
}

func GetTest(router *gin.Engine, param GetParam) (Message, http.Header, error) {
	// 构造get请求
	req := httptest.NewRequest(http.MethodGet, param.Uri, nil)
	if param.Header != nil {
		for k, v := range param.Header {
			req.Header.Set(k, v)
		}
	}
	if param.Query != nil {
		queryList := make([]string, len(param.Query))
		for k, v := range param.Query {
			queryList = append(queryList, fmt.Sprintf("%s=%s", k, v))
		}
		req.RequestURI = fmt.Sprintf("%s?%s", req.RequestURI, strings.Join(queryList, "&"))
	}
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
		return Message{}, nil, err
	}
	var msg Message
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return Message{}, nil, err
	}
	return msg, result.Header, nil
}

func PostTest(router *gin.Engine, param PostParam) (Message, http.Header, error) {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(param.Body)

	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest(http.MethodPost, param.Uri, bytes.NewReader(jsonByte))
	if param.Header != nil {
		for k, v := range param.Header {
			req.Header.Set(k, v)
		}
	}
	if param.Query != nil {
		queryList := make([]string, len(param.Query))
		for k, v := range param.Query {
			queryList = append(queryList, fmt.Sprintf("%s=%s", k, v))
		}
		req.RequestURI = fmt.Sprintf("%s?%s", req.RequestURI, strings.Join(queryList, "&"))
	}
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
		return Message{}, nil, err
	}
	var msg Message
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return Message{}, nil, err
	}
	return msg, result.Header, nil
}
