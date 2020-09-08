package router

import (
	"sync"

	"github.com/BoynChan/GopherProxy/router/user"
	"github.com/gin-gonic/gin"
)

// Author:Boyn
// Date:2020/9/8
var once sync.Once

func InitRouter(r *gin.Engine) {
	once.Do(func() {
		user.InitUserRouter(r)
	})
}
