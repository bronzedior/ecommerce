package grpc

import (
	"context"
	"user/cmd/user/usecase"
	"user/proto/userpb"
)

type GRPCServer struct {
	userpb.UnimplementedUserServiceServer
	UserUsecase usecase.UserUsecase
}

func (s *GRPCServer) GetUserInfoByUserID(ctx context.Context, req *userpb.GetUserInfoRequest) (*userpb.GetUserInfoResult, error) {
	userInfo, err := s.UserUsecase.GetUserByUserID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserInfoResult{
		Id:    userInfo.ID,
		Name:  userInfo.Name,
		Email: userInfo.Email,
		Role:  userInfo.Role,
	}, nil
}
