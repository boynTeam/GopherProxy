package main

import (
	"github.com/BoynChan/GopherProxy/internal/loadbalance"
	"github.com/BoynChan/GopherProxy/internal/proxy"
	_ "github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
)

// Author:Boyn
// Date:2020/9/1

func main() {
	proxyHandler, err := proxy.NewHttpProxyHandler("TEST_WEBSOCKET_PROXY", loadbalance.Random)
	if err != nil {
		panic(err)
	}
	r := gin.New()
	r.GET("/*path", gin.Logger(), proxyHandler)

	if err := r.Run(":2000"); err != nil {
		panic(err)
	}
}
