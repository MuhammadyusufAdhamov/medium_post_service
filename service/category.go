package service

import (
	"context"
	pb "github.com/MuhammadyusufAdhamov/medium_post_service/genproto/post_service"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type CategoryService struct {
	storage storage.StorageI
	pb.UnimplementedCategoryServiceServer
}

func NewCategoryService(strg storage.StorageI) *CategoryService {
	return &CategoryService{
		storage: strg,
	}
}

func (s *CategoryService) Create(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	category, err := s.storage.Category().Create(&repo.Category{
		Title: req.Title,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}
	return parseCategoryModel(category), nil
}

func parseCategoryModel(c *repo.Category) *pb.Category {
	return &pb.Category{
		Id:        c.ID,
		Title:     c.Title,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
	}
}
