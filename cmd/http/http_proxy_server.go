package main

import (
	"github.com/BoynChan/GopherProxy/internal/loadbalance"
	"github.com/BoynChan/GopherProxy/internal/middleware"
	"github.com/BoynChan/GopherProxy/internal/proxy"
	_ "github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Start Proxy Server.
// Author:Boyn
// Date:2020/8/31

func main() {
	serviceName := "TEST_HTTP_PROXY"
	proxyHandler, err := proxy.NewHttpProxyHandler(serviceName, loadbalance.Random)
	if err != nil {
		panic(err)
	}
	certFile := viper.GetString("Http.Cert")
	keyFile := viper.GetString("Http.Key")
	r := gin.New()

	// 速度限流器
	limiter := middleware.NewRateLimiter(2, 4)

	// 断路器
	circuitBreaker := middleware.NewCircuitBreaker("http_proxy", false)

	r.GET("/*path", gin.Logger(), circuitBreaker.GinMiddleWare(), limiter.GinMiddleWare(), proxyHandler)

	if err := r.RunTLS(":2000", certFile, keyFile); err != nil {
		panic(err)
	}
}
