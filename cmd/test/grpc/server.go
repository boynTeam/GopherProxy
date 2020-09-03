package main

import (
	"fmt"
	pb "github.com/BoynChan/GopherProxy/internal/grpc/proto"
	"github.com/BoynChan/GopherProxy/internal/grpc/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

// Author:Boyn
// Date:2020/9/2

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50055))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server.EchoServerImpl{})
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
