package server

import (
	"context"
	"io"

	pb "github.com/BoynChan/GopherProxy/internal/grpc/proto"
	"github.com/sirupsen/logrus"
)

// Author:Boyn
// Date:2020/9/2
const (
	streamingCount = 10
)

type EchoServerImpl struct {
}

func (e *EchoServerImpl) UnaryEcho(ctx context.Context, request *pb.EchoRequest) (*pb.EchoResponse, error) {
	logrus.Infof("UnaryEcho: receive message: %s", request.Message)
	return &pb.EchoResponse{Message: request.Message}, nil
}

func (e *EchoServerImpl) ServerStreamingEcho(request *pb.EchoRequest, server pb.Echo_ServerStreamingEchoServer) error {
	logrus.Infof("ServerStreamingEcho: receive message: %s", request.Message)
	for i := 0; i < streamingCount; i++ {
		err := server.Send(&pb.EchoResponse{Message: request.Message})
		if err != nil {
			logrus.Errorf("ServerStreamingEcho: receive message error:%v", err)
			return err
		}
	}
	return nil
}

func (e *EchoServerImpl) ClientStreamingEcho(server pb.Echo_ClientStreamingEchoServer) error {
	logrus.Infof("ClientStreamingEcho: state receive message")
	var message string
	for {
		in, err := server.Recv()
		if err == io.EOF {
			logrus.Infof("ClientStreamingEcho: ")
			return server.SendAndClose(&pb.EchoResponse{Message: message})
		}
		if err != nil {
			return err
		}
		message = in.Message
		logrus.Infof("ClientStreamingEcho: receive message %s", message)
	}
}

func (e *EchoServerImpl) BidirectionalStreamingEcho(server pb.Echo_BidirectionalStreamingEchoServer) error {
	logrus.Infof("BidirectionalStreamingEcho: state receive message")
	for {
		in, err := server.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		logrus.Infof("ClientStreamingEcho: receive message %s", in.Message)
		if err = server.Send(&pb.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}
