package controller

import (
	"cms/pb"
	"cms/services"
	"cms/services/record"
	"cms/services/schema"
	"cms/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedSchemaServiceServer
	pb.UnimplementedRecordServiceServer

	schemaSvc services.SchemaService
	recordSvc services.RecordService
}

func NewServer(
	db *gorm.DB,
	log utils.Logger,
) *grpc.Server {
	s := &server{
		schemaSvc: schema.NewSchemaService(log, db),
		recordSvc: record.NewRecordService(log, db),
	}

	sv := grpc.NewServer()

	pb.RegisterSchemaServiceServer(sv, s)
	pb.RegisterRecordServiceServer(sv, s)

	reflection.Register(sv)

	return sv
}
