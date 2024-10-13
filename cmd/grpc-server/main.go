package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/IRXCI/chat/config"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	desc "github.com/IRXCI/chat-server/pkg/chat"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	desc.UnimplementedChatAPIServer
	pool *pgxpool.Pool
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "../../.env", "path to config file")
}

func (s *server) CreateChat(_ context.Context, _ *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {

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

	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatAPIServer(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())
}
