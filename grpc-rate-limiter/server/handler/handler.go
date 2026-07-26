package handler

import (
	"context"
	"grpc-rate-limiter/proto/hello"
)

type Server struct {
	hello.UnimplementedGreeterServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SayHello(ctx context.Context, req *hello.GreeterRequest) (*hello.GreeterResponse, error) {
	return &hello.GreeterResponse{
		Message: "HI there " + req.Name,
	}, nil
}
