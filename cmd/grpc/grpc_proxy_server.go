package main

import (
	"fmt"
	"log"
	"net"

	"github.com/BoynChan/GopherProxy/internal/loadbalance"
	"github.com/BoynChan/GopherProxy/internal/proxy"
	"github.com/spf13/viper"
)

// Author:Boyn
// Date:2020/9/2

const (
	port = ":50051"
)

func main() {
	addrs := viper.GetStringSlice("Grpc.RealServer")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s, err := proxy.NewGrpcProxyServer(addrs, loadbalance.Random)
	if err != nil {
		panic(err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
