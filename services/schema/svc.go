package schema

import (
	"cms/models"
	"cms/services"
	"cms/services/internal"
	"cms/services/record"
	"cms/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type svc struct {
	log         utils.Logger
	repo        *repo
	recordSvc   services.RecordService
	internalSvc services.InternalService
}

func NewSchemaService(log utils.Logger, db *gorm.DB) services.SchemaService {
	return &svc{
		log:         log,
		repo:        newRepo(db, log),
		recordSvc:   record.NewRecordService(log, db),
		internalSvc: internal.NewInternalService(log, db),
	}
}

func (s *svc) CreateSchema(
	ctx context.Context, schema models.Schema, metaData []models.SchemaMetaData,
) (string, error) {
	schemaMetaData, err := s.internalSvc.FetchMetaData(ctx, schema.SchemaName, true)
	if err != nil {
		return "", err
	}

	if len(schemaMetaData) > 0 {
		return schemaMetaData[0].DisplaySchemaName, fmt.Errorf("table name already exists")
	}

	txErr := s.repo.withTransaction(func(tx *gorm.DB) error {
		err = s.repo.createSchemaMetaData(ctx, tx, metaData)
		if err != nil {
			return err
		}

		return s.repo.createSchema(ctx, tx, schema)
	})

	return "", txErr
}

func (s *svc) ListSchemas(
	ctx context.Context,
) ([]models.SchemaDetail, error) {
	return s.repo.listSchemas(ctx)
}

func (s *svc) GetSchema(
	ctx context.Context, schemaName string, filter models.Filter,
) (models.GetSchemaResponse, error) {
	var result models.GetSchemaResponse
	var err error
	result.MetaData, err = s.internalSvc.FetchMetaData(ctx, schemaName, false)
	if err != nil {
		return result, err
	}

	if len(result.MetaData) == 0 {
		return result, fmt.Errorf(fmt.Sprintf("table name not exists: %s", schemaName))
	}

	result.Data, err = s.recordSvc.GetRecords(ctx, result.MetaData[0].SystemSchemaName, filter)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *svc) DropSchema(
	ctx context.Context, schemaName string,
) error {
	schemaMetaData, err := s.internalSvc.FetchMetaData(ctx, schemaName, false)
	if err != nil {
		return err
	}

	if len(schemaMetaData) == 0 {
		return fmt.Errorf("table name not exists")
	}

	txErr := s.repo.withTransaction(func(tx *gorm.DB) error {
		err = s.repo.markSchemaArchived(ctx, tx, schemaName)
		if err != nil {
			return err
		}

		return s.repo.dropSchema(ctx, tx, schemaMetaData[0].SystemSchemaName)

	})

	return txErr
}
