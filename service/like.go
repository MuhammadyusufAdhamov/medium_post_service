package service

import (
	pb "github.com/MuhammadyusufAdhamov/medium_post_service/genproto/post_service"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage/repo"
	"github.com/sirupsen/logrus"
)

type LikeService struct {
	pb.UnimplementedLikeServiceServer
	storage storage.StorageI
	logger *logrus.Logger
}

func NewLikeService(strg storage.StorageI, logger *logrus.Logger) *LikeService {
	return &LikeService{
		storage: strg,
		logger: logger,
	}
}

//func (s *LikeService) CreateOrUpdate(ctx context.Context, req *pb.Like) (*pb.Like, error) {
//	like, err := s.storage.Like().CreateOrUpdate(&repo.Like{
//		PostID: req.PostId,
//		UserID: req.UserId,
//		Status: req.Status,
//	})
//	if err != nil {
//		s.logger.WithError(err).Error("failed in create or update like")
//		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
//	}
//
//	return parseLikeModel(like), nil
//}
//
//func (s *LikeService) Get(ctx context.Context, req *pb.GetLikeRequest) (*pb.Like, error) {
//	like, err := s.storage.Like().Get(req.)
//}

func parseLikeModel(req *repo.Like) *pb.Like {
	return &pb.Like{
		Id: req.ID,
		UserId: req.UserID,
		PostId: req.PostID,
		Status: req.Status,
	}
}
