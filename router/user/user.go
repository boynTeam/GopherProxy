package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/BoynChan/GopherProxy/domain"
	"github.com/BoynChan/GopherProxy/dto"
	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Author:Boyn
// Date:2020/9/8

const (
	userCookieName = "GATEWAY_USER_INFO"
)

func InitUserRouter(r *gin.Engine) {
	userControlloer := r.Group("/admin")
	userControlloer.POST("/login", userLogin)
	userControlloer.POST("/user", registerUser)
	userControlloer.GET("/user/id/:id", findUserById)
	userControlloer.GET("/user/username/:username", findUserByName)
	userControlloer.PUT("/user/:id", updateUser)
	userControlloer.DELETE("/user/:id", updateUser)
}

func userLogin(c *gin.Context) {
	var userInput dto.AdminInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	admin := domain.Admin{
		UserName: userInput.UserName,
		Password: userInput.Password,
	}
	check, err := admin.LoginCheck(c, pkg.DefaultDB)
	if err != nil {
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.PasswordErrorCode).Message(err).Build())
		return
	}
	sessionInfo := dto.AdminUserSession{
		Id:        check.ID,
		UserName:  check.UserName,
		LoginTime: time.Now(),
	}
	marshal, _ := json.Marshal(sessionInfo)
	session, _ := pkg.CookieSession.New(c.Request, userCookieName)
	session.Values["info"] = string(marshal)
	session.Save(c.Request, c.Writer)
	c.JSON(http.StatusOK, pkg.NewMessageBuilder().Data(check).Build())

}

func registerUser(c *gin.Context) {
	var userInput dto.AdminInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	admin := domain.Admin{
		UserName: userInput.UserName,
		Password: userInput.Password,
	}
	err := admin.Save(c, pkg.DefaultDB)
	if err == pkg.DuplicateRegisterError {
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DuplicateRegisterErrorCode).Message(err).Build())
		return
	} else if err != nil {
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	admin = domain.Admin{
		UserName: userInput.UserName,
	}
	find, err := admin.Find(c, pkg.DefaultDB)
	if err != nil {
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	c.JSON(http.StatusOK, pkg.NewMessageBuilder().Data(find).Build())
}

func findUserByName(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	admin := domain.Admin{
		UserName: username,
	}
	user, err := admin.Find(c, pkg.DefaultDB)
	if err != nil {
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	c.JSON(http.StatusOK, pkg.NewMessageBuilder().Data(user).Build())
}

func findUserById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	atoi, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	admin := domain.Admin{
		Model: gorm.Model{
			ID: uint(atoi),
		},
	}
	user, err := admin.Find(c, pkg.DefaultDB)
	if err != nil {
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	c.JSON(http.StatusOK, pkg.NewMessageBuilder().Data(user).Build())
}

func updateUser(c *gin.Context) {
	var userInput dto.AdminInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	atoi, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	admin := domain.Admin{
		Model: gorm.Model{
			ID: uint(atoi),
		},
		UserName: userInput.UserName,
		Password: userInput.Password,
	}
	err = admin.Update(c, pkg.DefaultDB)
	if err != nil {
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	admin.Password = ""
	user, err := admin.Find(c, pkg.DefaultDB)
	if err != nil {
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	c.JSON(http.StatusOK, pkg.NewMessageBuilder().Data(user).Build())
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	atoi, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	admin := domain.Admin{
		Model: gorm.Model{
			ID: uint(atoi),
		},
	}
	err = admin.Delete(c, pkg.DefaultDB)
	if err != nil {
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	c.JSON(http.StatusOK, pkg.NewMessageBuilder().Build())
}
