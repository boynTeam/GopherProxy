package main

import (
	"github.com/BoynChan/GopherProxy/internal/loadbalance"
	"github.com/BoynChan/GopherProxy/internal/proxy"
	_ "github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Author:Boyn
// Date:2020/9/1

func main() {
	addrs := viper.GetString("Http.WebsockerServer")
	proxyHandler, err := proxy.NewHttpProxyHandler([]string{addrs}, loadbalance.Random)
	if err != nil {
		panic(err)
	}
	r := gin.New()
	r.GET("/*path", gin.Logger(), proxyHandler)

	if err := r.Run(":2000"); err != nil {
		panic(err)
	}
}
