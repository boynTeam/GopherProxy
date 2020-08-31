package main

import (
	"github.com/BoynChan/GopherProxy/internal/loadbalance"
	"github.com/BoynChan/GopherProxy/internal/proxy"
	_ "github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Start Proxy Server.
// Author:Boyn
// Date:2020/8/31

func main() {
	addrs := viper.GetStringSlice("Http.RealServer")
	proxyHandler, err := proxy.NewHttpProxyHandler(addrs, loadbalance.Random)
	if err != nil {
		panic(err)
	}
	r := gin.New()
	r.GET("/*path", proxyHandler)

	if err := r.Run(":2000"); err != nil {
		panic(err)
	}
}
