package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/BoynChan/GopherProxy/pkg"
	_ "github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Start a simple real http server.
// Author:Boyn
// Date:2020/8/31

func newGinServer(addr string) *gin.Engine {
	r := gin.Default()
	r.GET("/*path", func(c *gin.Context) {
		// in this handler, we just simply send some basic info back to proxy response.
		req := c.Request
		urlPath := fmt.Sprintf("%s://%s%s", req.URL.Scheme, addr, req.RequestURI)
		realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-Ip"))
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Data(gin.H{
			"path": urlPath,
			"ip":   realIP,
		}).Build())
	})
	return r
}

func main() {
	wg := sync.WaitGroup{}
	addrs := viper.GetStringSlice("Http.RealServer")
	certFile := viper.GetString("Http.Cert")
	keyFile := viper.GetString("Http.Key")
	wg.Add(len(addrs))
	for _, a := range addrs {
		go func(addr string) {
			port := strings.Split(addr, "/")[2]
			if err := newGinServer(port).RunTLS(port, certFile, keyFile); err != nil {
				wg.Done()
			}
		}(a)
	}
	wg.Wait()
}
