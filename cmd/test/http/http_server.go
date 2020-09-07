package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/BoynChan/GopherProxy/internal/register"
	"github.com/BoynChan/GopherProxy/pkg"
	_ "github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Start a simple real http server.
// Author:Boyn
// Date:2020/8/31

const (
	zkRegisterPrefix = "/http_real_server"
)

type HttpServer struct {
	Addr        string
	ServiceName string
	r           *gin.Engine
	certFile    string
	keyFile     string
}

func NewHttpServer(addr, serviceName string, file ...string) *HttpServer {
	if len(file) == 2 {
		return &HttpServer{
			r:           gin.Default(),
			Addr:        addr,
			certFile:    file[0],
			keyFile:     file[1],
			ServiceName: serviceName,
		}
	}
	return &HttpServer{
		r:           gin.Default(),
		Addr:        addr,
		ServiceName: serviceName,
	}
}

func (h *HttpServer) Get(path string, f ...gin.HandlerFunc) {
	h.r.GET(path, f...)
}
func (h *HttpServer) Post(path string, f ...gin.HandlerFunc) {
	h.r.POST(path, f...)
}

func (h *HttpServer) register() error {
	zkHostIp := os.Getenv("ZK_HOST_IP")
	if zkHostIp == "" {
		return errors.New("no zookeeper running")
	}
	go func() {
		manager := register.NewZkManager(zkHostIp)
		_ = manager.GetConnect()
		servicePrefix := fmt.Sprintf("%s/%s", zkRegisterPrefix, h.ServiceName)
		err := manager.RegistServerTmpNode(servicePrefix, h.Addr, []byte(fmt.Sprintf("http://%s", h.Addr))...)
		if err != nil {
			logrus.Errorf("Register error: %v  addr: %s", err, h.Addr)
		}
	}()
	return nil
}

func (h *HttpServer) Run() error {
	err := h.register()
	if err != nil {
		return err
	}
	return h.r.Run(h.Addr)
}

func (h *HttpServer) RunTLS() error {
	err := h.register()
	if err != nil {
		return err
	}
	return h.r.RunTLS(h.Addr, h.certFile, h.keyFile)
}

func newGinServer(addr, serviceName string, file ...string) *HttpServer {
	r := NewHttpServer(addr, serviceName, file...)
	r.Get("/*path", func(c *gin.Context) {
		// in this handler, we just simply send some basic info back to proxy response.
		req := c.Request
		urlPath := fmt.Sprintf("https://%s%s", addr, req.RequestURI)
		retData := make(map[string]interface{})
		retData["RemoteAddr"] = req.RemoteAddr
		retData["X-Forwarded-For"] = req.Header.Get("X-Forwarded-For")
		retData["X-Real-Ip"] = req.Header.Get("X-Real-Ip")
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Data(gin.H{
			"path": urlPath,
			"ip":   retData,
		}).Build())
	})
	return r
}

func main() {
	wg := sync.WaitGroup{}
	os.Setenv("ZK_HOST_IP", "localhost:2181")
	addrs := viper.GetStringSlice("Http.RealServer")
	certFile := viper.GetString("Http.Cert")
	keyFile := viper.GetString("Http.Key")
	serviceName := "TEST_HTTP_PROXY"
	wg.Add(len(addrs))
	for _, a := range addrs {
		go func(addr string) {
			hostPortAddr := strings.Split(addr, "/")[2]
			if err := newGinServer(hostPortAddr, serviceName, certFile, keyFile).RunTLS(); err != nil {
				logrus.Errorf("Run error %v", err)
				wg.Done()
			}
		}(a)
	}
	wg.Wait()
}
