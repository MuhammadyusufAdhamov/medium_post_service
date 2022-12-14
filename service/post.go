package service

import (
	"context"
	pb "github.com/MuhammadyusufAdhamov/medium_post_service/genproto/post_service"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type PostService struct {
	pb.UnimplementedPostServiceServer
	storage storage.StorageI
}

func NewPostService(strg storage.StorageI) *PostService {
	return &PostService{
		storage: strg,
	}
}

func (s *PostService) Create(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	post, err := s.storage.Post().Create(&repo.Post{
		Title:       req.Title,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		UserID:      req.UserId,
		CategoryID:  req.CategoryId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}
	return parsePostModel(post), nil
}

func parsePostModel(post *repo.Post) *pb.Post {
	return &pb.Post{
		Id:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		ImageUrl:    post.ImageUrl,
		UserId:      post.UserID,
		CategoryId:  post.CategoryID,
		CreatedAt:   post.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   post.UpdatedAt.Format(time.RFC3339),
		ViewsCount:  post.ViewsCount,
	}
}

func (s *PostService) Get(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	post, err := s.storage.Post().Get(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	return parsePostModel(post), nil
}


func (s *PostService) GetAll(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	result, err := s.storage.User().GetAll(&repo.GetAllUsersParams{
		Limit:  req.Limit,
		Page:   req.Page,
		Search: req.Search,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	response := pb.GetAllUsersResponse{
		Count: result.Count,
		Users: make([]*pb.User, 0),
	}

	for _, post := range result.Users {
		response.Users = append(response.Users, parseUserModel(user))
	}

	return &response, nil
}

func (s *PostService) Update(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	post, err := s.storage.Post().(&repo.User{
		ID:              req.Id,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     req.PhoneNumber,
		Gender:          req.Gender,
		Username:        req.Username,
		ProfileImageUrl: req.ProfileImageUrl,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	return parsePostModel(post), nil
}

func (s *PostService) Delete(ctx context.Context, req *pb.GetPostRequest) (*emptypb.Empty, error) {
	err := s.storage.Post().Delete(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	return &emptypb.Empty{}, nil
}
