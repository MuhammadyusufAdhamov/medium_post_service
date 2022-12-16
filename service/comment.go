package service

import (
	"context"
	pb "github.com/MuhammadyusufAdhamov/medium_post_service/genproto/post_service"
	grpcPkg "github.com/MuhammadyusufAdhamov/medium_post_service/pkg/grpc_client"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage/repo"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type CommentService struct {
	pb.UnimplementedCommentServiceServer
	storage    storage.StorageI
	grpcClient grpcPkg.GrpcClientI
	logger     *logrus.Logger
}

func NewCommentService(strg storage.StorageI, grpc grpcPkg.GrpcClientI, logger *logrus.Logger) *CommentService {
	return &CommentService{
		storage:    strg,
		grpcClient: grpc,
		logger:     logger,
	}
}

func (s *CommentService) Create(ctx context.Context, req *pb.Comment) (*pb.Comment, error) {
	comment, err := s.storage.Comment().Create(&repo.Comment{
		Description: req.Description,
		UserID:      req.UserId,
		PostID:      req.PostId,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create comment")
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	return &pb.Comment{
		Id:          comment.ID,
		PostId:      comment.PostID,
		UserId:      comment.UserID,
		Description: comment.Description,
		CreatedAt:   comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   comment.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func parseCommentModel(comment *repo.Comment) *pb.Comment {
	return &pb.Comment{
		Id:          comment.ID,
		PostId:      comment.PostID,
		UserId:      comment.UserID,
		Description: comment.Description,
		CreatedAt:   comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   comment.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *CommentService) Delete(ctx context.Context, req *pb.GetCommentRequest) (*emptypb.Empty, error) {
	err := s.storage.Comment().Delete(req.Id)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete comment")
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *CommentService) GetAll(ctx context.Context, req *pb.GetAllCommentsRequest) (*pb.GetAllCommentsResponse, error) {
	result, err := s.storage.Comment().GetAll(&repo.GetAllCommentsParams{
		Limit:  req.Limit,
		Page:   req.Page,
		Search: req.Search,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	response := pb.GetAllCommentsResponse{
		Count:    result.Count,
		Comments: make([]*pb.Comment, 0),
	}

	for _, comment := range result.Comments {
		response.Comments = append(response.Comments, parseCommentModel(comment))
	}

	return &response, nil
}

func (s *CommentService) Update(ctx context.Context, req *pb.Comment) (*pb.Comment, error) {
	comment, err := s.storage.Comment().Update(&repo.Comment{
		UserID:      req.UserId,
		PostID:      req.PostId,
		Description: req.Description,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	return parseCommentModel(comment), nil
}
