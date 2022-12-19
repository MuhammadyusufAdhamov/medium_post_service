package main

import (
	"fmt"
	"github.com/MuhammadyusufAdhamov/medium_post_service/config"
	pb "github.com/MuhammadyusufAdhamov/medium_post_service/genproto/post_service"
	grpcPkg "github.com/MuhammadyusufAdhamov/medium_post_service/pkg/grpc_client"
	"github.com/MuhammadyusufAdhamov/medium_post_service/pkg/logger"
	"github.com/MuhammadyusufAdhamov/medium_post_service/service"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.Load(".")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	strg := storage.NewStoragePg(psqlConn)

	grpcConn, err := grpcPkg.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	log := logger.New()

	postService := service.NewPostService(strg)
	categoryService := service.NewCategoryService(strg, log)
	likeService := service.NewLikeService(&strg, log)
	commentService := service.NewCommentService(strg, grpcConn, log)

	lis, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, postService)
	pb.RegisterCategoryServiceServer(s, categoryService)
	pb.RegisterLikeServiceServer(s, likeService)
	pb.RegisterCommentServiceServer(s, commentService)

	log.Println("Grpc server started in port ", cfg.GrpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while listening: %v", err)
	}
}
