package user

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/BoynChan/GopherProxy/dto"
	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	data := message.Data.(map[string]interface{})
	assert.Equal(t, data["user_name"], "test01")
}

func TestRegisterUser(t *testing.T) {
	input := dto.AdminInput{
		UserName: "test02",
		Password: hashPassword("123456"),
	}
	message, err := pkg.PostTest("/admin/user", r, input)
	require.Nil(t, err)
	require.Equal(t, 2003, message.Code)
	data, ok := message.Data.(map[string]interface{})
	if ok {
		assert.Equal(t, data["user_name"], "test02")
	}
}

func TestLoginUser(t *testing.T) {
	input := dto.AdminInput{
		UserName: "test02",
		Password: hashPassword("123456"),
	}
	message, err := pkg.PostTest("/admin/login", r, input)
	require.Nil(t, err)
	require.Equal(t, 200, message.Code)
	data, ok := message.Data.(map[string]interface{})
	if ok {
		assert.Equal(t, data["user_name"], "test02")
	}
}

func hashPassword(password string) string {
	digest := md5.New()
	sum := digest.Sum([]byte(password))
	return base64.StdEncoding.EncodeToString(sum)
}
