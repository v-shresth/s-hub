package handler

import (
	"cms/pb"
	"cms/utils"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateSchema(
	ctx context.Context, req *pb.CreateSchemaRequest,
) (*emptypb.Empty, error) {
	schemaName := req.SchemaName
	metaData, err := utils.ValidateCreateSchemaRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	schema, err := utils.ConvertCreateSchemaApiReqToDbModel(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	userSchemaName, err := s.SchemaSvc.CreateSchema(ctx, schema, metaData)
	if err != nil {
		if schemaName != "" {
			if schemaName == userSchemaName {
				return nil, status.Errorf(codes.InvalidArgument, "Table already exists")
			} else {
				return nil, status.Error(codes.InvalidArgument, fmt.Sprintf(`Table signature is same as the existing table "%s". Please change table name.`, userSchemaName))
			}
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return nil, nil
}

func (s *Server) ListSchemas(
	ctx context.Context, req *pb.ListSchemasRequest,
) (*pb.ListSchemasResponse, error) {
	schemas, err := s.SchemaSvc.ListSchemas(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var resp = &pb.ListSchemasResponse{}
	for _, schema := range schemas {
		resp.Schemas = append(resp.Schemas, &pb.SchemaDetail{
			SchemaName: schema.SchemaName,
			NoOfFields: int32(schema.NoOfFields),
		})
	}
	if len(schemas) > 0 {
		resp.TotalSchemas = int32(schemas[0].TotalSchemas)
	}
	return resp, nil
}

func (s *Server) GetSchema(
	ctx context.Context, req *pb.GetSchemaRequest,
) (*pb.GetSchemaResponse, error) {
	filter, err := utils.ValidateGetSchemaRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	resp, err := s.SchemaSvc.GetSchema(ctx, req.SchemaName, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	apiResp, err := utils.ConvertGetSchemaDbRespToApiResp(resp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return apiResp, nil
}

func (s *Server) DropSchema(
	ctx context.Context, req *pb.DropSchemaRequest,
) (*emptypb.Empty, error) {
	err := utils.ValidateDropSchemaRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.SchemaSvc.DropSchema(ctx, req.SchemaName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return nil, nil
}
