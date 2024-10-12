package middlewares

import (
	"cms/clients"
	handler2 "cms/handler"
	"cms/services/record"
	"cms/services/schema"
	"cms/utils"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckUserPreSetup(
	config clients.Config,
	logger clients.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, ok := exemptedMethods[info.FullMethod]; !ok {
			metadata := utils.ExtractMetadata(ctx)

			userDb, err := clients.GetUserDb(ctx, metadata.AuthedUserId, config, logger)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}

			server, ok := info.Server.(*handler2.Server)
			if !ok {
				return nil, status.Error(codes.Internal, "Server Error")
			}

			server.UserDb = userDb
			server.RecordSvc = record.NewRecordService(logger, userDb, config)
			server.SchemaSvc = schema.NewSchemaService(logger, userDb, config)
		}
		h, err := handler(ctx, req)
		return h, err
	}
}
