package record

import (
	"cms/models"
	"cms/pb"
	"cms/services"
	"cms/services/internal"
	"cms/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type svc struct {
	log         utils.Logger
	repo        *repo
	internalSvc services.InternalService
}

func NewRecordService(log utils.Logger, db *gorm.DB) services.RecordService {
	return &svc{
		log:         log,
		repo:        newRepo(db, log),
		internalSvc: internal.NewInternalService(log, db),
	}
}

func (s *svc) CreateRecord(
	ctx context.Context, req *pb.CreateRecordRequest,
) ([]models.SchemaMetaData, []map[string]interface{}, error) {
	metaData, err := s.internalSvc.FetchMetaData(ctx, req.SchemaName, false)
	if err != nil {
		return nil, nil, err
	}

	if len(metaData) == 0 {
		return nil, nil, fmt.Errorf(fmt.Sprintf("table name not exists: %s", req.SchemaName))
	}

	records, err := utils.ConvertApiRecordsToDbRecords(req.Records, metaData)
	if err != nil {
		return nil, nil, err
	}

	dbRecords, err := s.repo.createRecord(ctx, metaData[0].SystemSchemaName, records)

	return metaData, dbRecords, err
}

func (s *svc) GetRecords(
	ctx context.Context, schemaName string, filter models.Filter,
) ([]map[string]interface{}, error) {
	return s.repo.getRecords(ctx, schemaName, filter)
}

func (s *svc) GetRecord(
	ctx context.Context, schemaName string, recordId int,
) ([]models.SchemaMetaData, map[string]interface{}, error) {
	metaData, err := s.internalSvc.FetchMetaData(ctx, schemaName, false)
	if err != nil {
		return nil, nil, err
	}

	if len(metaData) == 0 {
		return nil, nil, fmt.Errorf(fmt.Sprintf("table name not exists: %s", schemaName))
	}

	record, err := s.repo.getRecord(ctx, metaData[0].SystemSchemaName, recordId)
	if err != nil {
		return nil, nil, err
	}

	return metaData, record, nil
}

func (s *svc) DeleteRecord(
	ctx context.Context, schemaName string, recordId int,
) error {
	metaData, err := s.internalSvc.FetchMetaData(ctx, schemaName, false)
	if err != nil {
		return err
	}

	if len(metaData) == 0 {
		return fmt.Errorf(fmt.Sprintf("table name not exists: %s", schemaName))
	}

	err = s.repo.deleteRecord(ctx, metaData[0].SystemSchemaName, recordId)
	if err != nil {
		return err
	}

	return nil
}

func (s *svc) UpdateRecord(
	ctx context.Context, req *pb.UpdateRecordRequest,
) ([]models.SchemaMetaData, map[string]interface{}, error) {
	metaData, err := s.internalSvc.FetchMetaData(ctx, req.SchemaName, false)
	if err != nil {
		return nil, nil, err
	}

	if len(metaData) == 0 {
		return nil, nil, fmt.Errorf(fmt.Sprintf("table name not exists: %s", req.SchemaName))
	}

	records, err := utils.ConvertApiRecordsToDbRecords([]*pb.Record{req.Record}, metaData)
	if err != nil {
		return nil, nil, err
	}

	err = s.repo.updateRecord(ctx, metaData[0].SystemSchemaName, int(req.RecordId), records[0])

	return metaData, records[0], err
}
