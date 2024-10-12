package server

import (
	"cms/clients"
	"cms/handler"
	"cms/middlewares"
	"cms/pb"
	"cms/services/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func NewServer(
	systemDb *gorm.DB,
	log clients.Logger,
	config clients.Config,
) *grpc.Server {
	tokenMaker := clients.NewTokenMaker(log, config)
	s := &handler.Server{
		Config:     config,
		Log:        log,
		SystemDb:   systemDb,
		TokenMaker: tokenMaker,
		UserSvc:    users.NewUserService(log, systemDb, tokenMaker),
	}

	sv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middlewares.AuthInterceptor,
			middlewares.CheckUserPreSetup(s.Config, s.Log),
		))

	pb.RegisterSchemaServiceServer(sv, s)
	pb.RegisterRecordServiceServer(sv, s)
	pb.RegisterUserServiceServer(sv, s)

	reflection.Register(sv)

	return sv
}
