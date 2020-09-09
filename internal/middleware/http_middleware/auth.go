package http_middleware

import (
	"encoding/json"

	"github.com/BoynChan/GopherProxy/dto"
	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
)

// Author:Boyn
// Date:2020/9/9

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := pkg.CookieSession.Get(c.Request, pkg.UserCookieName)
		if err != nil {
			c.AbortWithStatusJSON(200, pkg.ErrorMessage(pkg.NotLoginErrorCode))
			return
		}
		var us dto.AdminUserSession
		err = json.Unmarshal([]byte(session.Values["info"].(string)), &us)
		if err != nil {
			c.AbortWithStatusJSON(200, pkg.ErrorMessage(pkg.NotLoginErrorCode))
			return
		}
		c.Set("info", us)
		c.Next()
	}
}
