package main

import (
	"fmt"
	"log"
	"net"

	"github.com/BoynChan/GopherProxy/internal/loadbalance"
	"github.com/BoynChan/GopherProxy/internal/proxy"
)

// Author:Boyn
// Date:2020/9/2

const (
	port = ":50051"
)

func main() {
	serviceName := "TEST_GRPC_PROXY"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s, err := proxy.NewGrpcProxyServer(serviceName, loadbalance.Random)
	if err != nil {
		panic(err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
