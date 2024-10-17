package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/IRXCI/chat-server/config"
	desc "github.com/IRXCI/chat-server/pkg/chat"
	sq "github.com/Masterminds/squirrel"
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

func (s *server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {

	builderCreateChat := sq.Insert("chat").
		PlaceholderFormat(sq.Dollar).
		Columns("user_names").
		Values(pq.Array(req.GetUsernames())).
		Suffix("RETURNING id")

	query, args, err := builderCreateChat.ToSql()
	if err != nil {
		log.Printf("failes to build query: %v", err)
		return nil, err
	}

	var ChatID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&ChatID)
	if err != nil {
		log.Printf("failed to insert chat: %v", err)
		return nil, err
	}

	log.Printf("Insert chat with id: %d", ChatID)

	return &desc.CreateChatResponse{
		Id: ChatID,
	}, nil
}

func (s *server) DeleteChat(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {

	builderDeleteChat := sq.Delete("chat").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDeleteChat.ToSql()
	if err != nil {
		log.Printf("failes to build query: %v", err)
		return nil, err
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failes to delete chat: %v", err)
		return nil, err
	}

	log.Printf("Deleted chat with id: %v", req.GetId())

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(_ context.Context, _ *desc.SendChatRequest) (*emptypb.Empty, error) {

	log.Printf("SendMessage is working...")
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

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
