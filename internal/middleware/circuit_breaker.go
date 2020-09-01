package middleware

import (
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
)

// Author:Boyn
// Date:2020/9/1

type CircuitBreaker struct {
	command string
}

func (cir *CircuitBreaker) GinMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := hystrix.Do(cir.command, func() error {
			c.Next()
			statusCode := c.Writer.Status()
			if statusCode != 200 {
				return errors.New("downstream error")
			}
			return nil
		}, nil)
		if err != nil {
			//加入自动降级处理，如获取缓存数据等
			c.JSON(http.StatusOK, pkg.NewMessageBuilder().Message(err.Error()).Build())
			c.Abort()
		}
	}
}

func NewCircuitBreaker(command string, openStream bool) *CircuitBreaker {
	hystrix.ConfigureCommand(command, hystrix.CommandConfig{
		Timeout:                1000, // 单次请求 超时时间
		MaxConcurrentRequests:  1,    // 最大并发量
		SleepWindow:            5000, // 熔断后多久去尝试服务是否可用
		RequestVolumeThreshold: 1,    // 验证熔断的 请求数量, 10秒内采样
		ErrorPercentThreshold:  1,    // 验证熔断的 错误百分比
	})

	if openStream {
		hystrixStreamHandler := hystrix.NewStreamHandler()
		hystrixStreamHandler.Start()
		go func() {
			err := http.ListenAndServe(net.JoinHostPort("", "2018"), hystrixStreamHandler)
			log.Fatal(err)
		}()
	}

	return &CircuitBreaker{command: command}
}
