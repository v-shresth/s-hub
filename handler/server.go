package handler

import (
	"cms/clients"
	"cms/pb"
	"cms/services"
	"gorm.io/gorm"
)

type Server struct {
	pb.UnimplementedSchemaServiceServer
	pb.UnimplementedRecordServiceServer
	pb.UnimplementedUserServiceServer

	SchemaSvc services.SchemaService
	RecordSvc services.RecordService
	UserSvc   services.UserService

	SystemDb *gorm.DB
	UserDb   *gorm.DB

	Config     clients.Config
	Log        clients.Logger
	TokenMaker clients.TokenMaker
}
