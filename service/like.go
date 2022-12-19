package service

import (
	"context"
	pb "github.com/MuhammadyusufAdhamov/medium_post_service/genproto/post_service"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage/repo"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LikeService struct {
	pb.UnimplementedLikeServiceServer
	storage storage.StorageI
	logger  *logrus.Logger
}

func NewLikeService(strg *storage.StorageI, log *logrus.Logger) *LikeService {
	return &LikeService{
		storage: *strg,
		logger:  log,
	}
}

func (s *LikeService) CreateOrUpdate(ctx context.Context, req *pb.Like) (*pb.Like, error) {
	like, err := s.storage.Like().CreateOrUpdate(&repo.Like{
		PostID: req.PostId,
		UserID: req.UserId,
		Status: req.Status,
	})
	if err != nil {
		s.logger.WithError(err).Error("error in create or update like service")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return parseLike(like), nil
}

func (s *LikeService) Get(ctx context.Context, req *pb.GetLikeRequest) (*pb.Like, error) {
	like, err := s.storage.Like().Get(req.UserId, req.PostId)
	if err != nil {
		s.logger.WithError(err).Error("error in get like service")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}
	return parseLike(like), nil
}

func (s *LikeService) GetLikesDislikesCount(ctx context.Context, req *pb.GetLikesRequest) (*pb.LikesDislikesCountResponse, error) {
	counts, err := s.storage.Like().GetLikesDislikesCount(req.PostId)
	if err != nil {
		s.logger.WithError(err).Error("error in get likes and dislike count in like service")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.LikesDislikesCountResponse{
		Likes:    counts.LikesCount,
		Dislikes: int32(counts.DislikesCount),
	}, nil
}

func parseLike(req *repo.Like) *pb.Like {
	return &pb.Like{
		Id:     req.ID,
		UserId: req.UserID,
		PostId: req.PostID,
		Status: req.Status,
	}
}
