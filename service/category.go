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

func (s *CategoryService) GetAll(ctx context.Context, req *pb.GetAllCategoriesRequest) (*pb.GetAllCategoriesResponse, error) {
	result, err := s.storage.Category().GetAll(&repo.GetAllCategoriesParams{
		Limit:  req.Limit,
		Page:   req.Page,
		Search: req.Search,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	response := pb.GetAllCategoriesResponse{
		Count:      result.Count,
		Categories: make([]*pb.Category, 0),
	}

	for _, category := range result.Categories {
		response.Categories = append(response.Categories, parseCategoryModel(category))
	}

	return &response, nil
}

func (s *CategoryService) Update(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	user, err := s.storage.Category().Update(&repo.Category{
		ID:    req.Id,
		Title: req.Title,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	return parseCategoryModel(user), nil
}

func (s *CategoryService) Delete(ctx context.Context, req *pb.GetCategoryRequest) error {
	err := s.storage.Category().Delete(req.Id)
	if err != nil {
		return err
	}

	return nil
}
