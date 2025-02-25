package handler

import (
	"cms/pb"
	"cms/utils"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateRecord(
	ctx context.Context, req *pb.CreateRecordRequest,
) (*pb.Record, error) {
	err := utils.ValidateCreateRecordRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	metaData, records, err := s.RecordSvc.CreateRecord(ctx, req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := utils.ConvertDbRecordsToApiRecords(records, metaData)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp[0], nil
}

func (s *Server) GetRecord(
	ctx context.Context, req *pb.GetRecordRequest,
) (*pb.Record, error) {
	err := utils.ValidateGetRecordRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	metaData, record, err := s.RecordSvc.GetRecord(ctx, req.SchemaName, int(req.RecordId))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := utils.ConvertDbRecordsToApiRecords([]map[string]interface{}{record}, metaData)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp[0], nil
}

func (s *Server) DeleteRecord(
	ctx context.Context, req *pb.DeleteRecordRequest,
) (*emptypb.Empty, error) {
	err := utils.ValidateDeleteRecordRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.RecordSvc.DeleteRecord(ctx, req.SchemaName, int(req.RecordId))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return nil, nil
}

func (s *Server) UpdateRecord(
	ctx context.Context, req *pb.UpdateRecordRequest,
) (*pb.Record, error) {
	err := utils.ValidateUpdateRecordRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	metaData, record, err := s.RecordSvc.UpdateRecord(ctx, req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := utils.ConvertDbRecordsToApiRecords([]map[string]interface{}{
		record,
	}, metaData)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp[0], nil
}
