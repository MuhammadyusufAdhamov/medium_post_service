package service

import (
	pb "github.com/MuhammadyusufAdhamov/medium_post_service/genproto/post_service"
	"github.com/MuhammadyusufAdhamov/medium_post_service/storage"
)

type LikeService struct {
	pb.UnimplementedLikeServiceServer
	storage storage.StorageI
}

func NewLikeService(strg storage.StorageI) *LikeService {
	return &LikeService{
		storage: strg,
	}
}

//
//func (s *LikeService) Create(ctx context.Context, req *pb.Like) (*pb.Like, error) {
//	like, err := s.storage.Like().CreateOrUpdate(&repo.Like{
//		PostID: req.PostId,
//		UserID: req.UserId,
//		Status: req.Status,
//	})
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
//	}
//	return parsePostModel(like), nil
//}
