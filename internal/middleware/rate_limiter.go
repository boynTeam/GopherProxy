package middleware

import (
	"net/http"

	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	time_rate "golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Author:Boyn
// Date:2020/9/1

type RateLimiter struct {
	*time_rate.Limiter
}

var (
	errGrpcOverRate = status.Errorf(codes.ResourceExhausted, "flow control over rate")
)

func (r *RateLimiter) GrpcMiddleWare() func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !r.Allow() {
			return errGrpcOverRate
		}
		err := handler(srv, ss)
		return err
	}
}

func (r *RateLimiter) GinMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !r.Allow() {
			c.JSON(http.StatusOK, pkg.
				NewMessageBuilder().
				Code(pkg.RateLimitErrorCode).
				Build())
			c.Abort()
			return
		}
		c.Next()
	}
}

// NewRateLimiter new a request speed rate limiter.
// We use a bucket to limit speed.
// @param:rate how many token will produce in one second.
// @param:burst biggest capacity
func NewRateLimiter(rate float64, burst int) *RateLimiter {
	return &RateLimiter{time_rate.NewLimiter(time_rate.Limit(rate), burst)}
}
