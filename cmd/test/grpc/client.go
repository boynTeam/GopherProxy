package main

import (
	"fmt"

	"github.com/BoynChan/GopherProxy/internal/grpc/client"
	"google.golang.org/grpc"
)

// Author:Boyn
// Date:2020/9/2

func main() {
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	echoClient := client.NewEchoClient(conn)
	echo, err := echoClient.UnaryEcho("Hi Grpc")
	fmt.Println(echo)
}
