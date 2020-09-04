package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/BoynChan/GopherProxy/internal/grpc/proto"
	"github.com/BoynChan/GopherProxy/internal/grpc/server"
	_ "github.com/BoynChan/GopherProxy/pkg"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// Author:Boyn
// Date:2020/9/2

func main() {
	addrs := viper.GetStringSlice("Grpc.RealServer")
	wg := sync.WaitGroup{}
	wg.Add(len(addrs))
	for _, addr := range addrs {
		go func(address string) {
			lis, err := net.Listen("tcp", address)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			fmt.Printf("server listening at %v\n", lis.Addr())
			s := grpc.NewServer()
			pb.RegisterEchoServer(s, &server.EchoServerImpl{Addr: address})
			if err := s.Serve(lis); err != nil {
				panic(err)
			}
			wg.Done()
		}(addr)
	}
	wg.Wait()
}
