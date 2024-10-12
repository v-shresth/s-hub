package handler

import (
	"cms/pb"
	"cms/utils"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(
	ctx context.Context,
	req *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {
	err := utils.ValidateRegisterRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dbUser := utils.ConvertApiRegisterRequestToDbRegisterRequest(req)

	user, err := s.UserSvc.Register(ctx, dbUser)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		User: utils.ConvertDbRegisterRequestToApiRegisterRequest(user),
	}, nil
}

func (s *Server) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (*pb.LoginResponse, error) {
	err := utils.ValidateLoginRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokenPair, err := s.UserSvc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *Server) Logout(
	ctx context.Context,
	req *empty.Empty,
) (*empty.Empty, error) {
	metadata := utils.ExtractMetadata(ctx)

	err := s.UserSvc.Logout(ctx, metadata.AuthedSessionId)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
