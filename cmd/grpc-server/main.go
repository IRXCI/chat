package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	desc "github.com/IRXCI/chat-server/pkg/chat"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedChatAPIServer
}

func (s *server) CreateChat(_ context.Context, _ *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	log.Printf("CreateChat working...")

	return &desc.CreateChatResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) DeleteChat(_ context.Context, _ *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	log.Printf("DeleteChat working...")
	return nil, nil
}

func (s *server) SendMessage(_ context.Context, _ *desc.SendChatRequest) (*emptypb.Empty, error) {
	log.Printf("SendMessage working...")
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatAPIServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
