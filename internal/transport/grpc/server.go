package grpc

import (
	"net"
	"google.golang.org/grpc"
	grpcHandler "payment/internal/handler/grpc"
	pb "payment/api/proto"
	"context"
)

type Server struct {
	server *grpc.Server
	port string  
}

func NewServer(port string, paymentHandler *grpcHandler.PaymentHandler) *Server {
	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, paymentHandler) 
	return &Server{
		server: s,
		port: port,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	return s.server.Serve(listener)
}

func (s *Server) Stop(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()
	select {
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	case <-done:
		//
	}
	return nil
}
