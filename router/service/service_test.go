package service

import (
	"testing"

	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Author:Boyn
// Date:2020/9/13

var r *gin.Engine

func init() {
	r = gin.Default()
	InitServiceRouter(r)
	pkg.InitDB()
}

func TestListService(t *testing.T) {
	message, _, err := pkg.GetTest(r, pkg.GetParam{
		Uri: "/service/list",
		Query: map[string]string{
			"page_size":  "20",
			"page_index": "1",
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, message.Code)
}

func TestGetServiceDetail(t *testing.T) {
	message, _, err := pkg.GetTest(r, pkg.GetParam{
		Uri: "/service/detail/1",
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, message.Code)
	if message.Code != 200 {
		t.Fatal("fail message:", message.Message)
	}
}
