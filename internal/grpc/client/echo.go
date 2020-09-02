package client

import (
	"context"
	"io"

	pb "github.com/BoynChan/GopherProxy/internal/grpc/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Author:Boyn
// Date:2020/9/2

type EchoClient struct {
	c pb.EchoClient
}

func NewEchoClient(c *grpc.ClientConn) *EchoClient {
	client := pb.NewEchoClient(c)
	return &EchoClient{c: client}
}

func (e *EchoClient) UnaryEcho(message string) (string, error) {
	resp, err := e.c.UnaryEcho(context.TODO(), &pb.EchoRequest{Message: message})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}

func (e *EchoClient) ServerStreamingEcho(message string) ([]string, error) {
	stream, err := e.c.ServerStreamingEcho(context.TODO(), &pb.EchoRequest{Message: message})
	if err != nil {
		return nil, err
	}
	messages := make([]string, 0)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return messages, nil
		}
		if err != nil {
			return nil, err
		}
		messages = append(messages, resp.Message)
	}
}

func (e *EchoClient) ClientStreamingEcho(messages []string) (string, error) {
	stream, err := e.c.ClientStreamingEcho(context.TODO())
	if err != nil {
		return "", err
	}
	for _, msg := range messages {
		err := stream.Send(&pb.EchoRequest{Message: msg})
		if err != nil {
			return "", err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}

func (e *EchoClient) BidirectionalStreamingEcho(messages []string) ([]string, error) {
	stream, err := e.c.BidirectionalStreamingEcho(context.TODO())
	if err != nil {
		return nil, err
	}
	go func() {
		// Send all requests to the server.
		for _, msg := range messages {
			if err := stream.Send(&pb.EchoRequest{Message: msg}); err != nil {
				logrus.Errorf("failed to send streaming: %v\n", err)
			}
		}
		stream.CloseSend()
	}()

	respMessages := make([]string, 0)
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return respMessages, nil
		}
		if err != nil {
			return nil, err
		}
		respMessages = append(respMessages, r.Message)
	}
}
