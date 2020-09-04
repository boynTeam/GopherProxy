package proxy

import (
	"context"
	"strings"

	"github.com/BoynChan/GopherProxy/internal/loadbalance"
	"github.com/BoynChan/GopherProxy/internal/middleware"
	"github.com/BoynChan/GopherProxy/internal/urls"
	proxy "github.com/e421083458/grpc-proxy/proxy"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Author:Boyn
// Date:2020/9/3

func NewGrpcProxyServer(urlSli []string, lbType loadbalance.Type) (*grpc.Server, error) {
	registerAddr := viper.GetString("ZookeeperAddr")
	dyUrls, err := urls.NewDynamicUrls(urlSli, lbType, registerAddr)
	if err != nil {
		return nil, err
	}
	director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		nextAddr, err := dyUrls.GetNext("")
		if err != nil {
			return ctx, nil, err
		}
		// 拒绝某些特殊请求
		if strings.HasPrefix(fullMethodName, "/com.example.internal.") {
			return ctx, nil, status.Errorf(codes.Unimplemented,
				"Unknown method")
		}
		c, err := grpc.DialContext(ctx, nextAddr, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
		md, _ := metadata.FromIncomingContext(ctx)
		outCtx, _ := context.WithCancel(ctx)
		outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
		return outCtx, c, err
	}

	limiter := middleware.NewRateLimiter(1, 2)
	return grpc.NewServer(
		grpc.ChainStreamInterceptor(limiter.GrpcMiddleWare()),
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director))), nil
}
