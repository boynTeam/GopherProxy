package user

import (
	"fmt"
	"testing"

	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Author:Boyn
// Date:2020/9/8

var r *gin.Engine

func init() {
	r = gin.Default()
	InitUserRouter(r)
	pkg.InitDB()
}

func TestFindUser(t *testing.T) {
	message, err := pkg.GetTest(fmt.Sprintf("/admin/user/id/%d", 3), r)
	assert.Nil(t, err)
	assert.Equal(t, message.Code, 200)
	assert.Equal(t, message.Code, 200)
	data := message.Data.(map[string]interface{})
	assert.Equal(t, data["user_name"], "test01")
}

func TestRegisterUser(t *testing.T) {

}
