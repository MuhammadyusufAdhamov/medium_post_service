package main

import (
	"fmt"
	pb "github.com/MuhammadyusufAdhamov/medium_post_service/genproto/post_service"
	"github.com/MuhammadyusufAdhamov/medium_post_service/service"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"log"
	"net"

	"github.com/MuhammadyusufAdhamov/medium_post_service/config"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage"
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

	postService := service.NewPostService(strg)
	categoryService := service.NewCategoryService(strg)

	lis, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterPostServiceServer(s, postService)
	pb.RegisterCategoryServiceServer(s, categoryService)

	log.Println("Grpc server started in port ", cfg.GrpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while listening: %v", err)
	}
}


