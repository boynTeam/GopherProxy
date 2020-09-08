package domain

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Author:Boyn
// Date:2020/9/8
var duplicateRegisterError = errors.New("不能重复注册")

type Admin struct {
	gorm.Model
	UserName string `json:"user_name" gorm:"column:user_name" description:"管理员用户名"`
	Salt     string `json:"salt" gorm:"column:salt" description:"盐"`
	Password string `json:"password" gorm:"column:password" description:"密码"`
}

func (t *Admin) LoginCheck(c *gin.Context, db *gorm.DB) (*Admin, error) {
	loginUser := Admin{
		UserName: t.UserName,
	}
	adminInfo, err := loginUser.Find(c, db)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if adminInfo.Password != generateSaltPassword(adminInfo.Salt, t.Password) {
		return nil, errors.New("wrong password")
	}
	return adminInfo, nil
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}

func (t *Admin) Save(c *gin.Context, db *gorm.DB) error {
	if t.Salt == "" {
		t.Salt = randomSalt()
	}
	search := &Admin{UserName: t.UserName}
	_, err := search.Find(c, db)
	if err == nil {
		return duplicateRegisterError
	}
	t.Password = generateSaltPassword(t.Salt, t.Password)
	return db.Save(t).Error
}

func (t *Admin) Find(c *gin.Context, db *gorm.DB) (*Admin, error) {
	out := &Admin{}
	err := db.Where(t).First(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *Admin) Update(c *gin.Context, db *gorm.DB) error {
	out := &Admin{}
	err := db.Where(t).First(out).Error
	if err != nil {
		return err
	}
	out.UserName = t.UserName
	out.Password = t.Password
	out.Salt = t.Salt
	return db.Save(out).Error
}

func (t *Admin) Delete(c *gin.Context, db *gorm.DB) error {
	return db.Delete(t).Error
}

func randomSalt() string {
	return uuid.NewV4().String()
}

func hashPassword(password string) string {
	digest := md5.New()
	sum := digest.Sum([]byte(password))
	return base64.StdEncoding.EncodeToString(sum)
}

func generateSaltPassword(salt, password string) string {
	digest := md5.New()
	sum := digest.Sum([]byte(fmt.Sprintf("%s%s", salt, password)))
	return base64.StdEncoding.EncodeToString(sum)
}
