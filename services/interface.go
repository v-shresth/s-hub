package services

import (
	"cms/models"
	"cms/pb"
	"context"
)

type SchemaService interface {
	CreateSchema(ctx context.Context, schema models.Schema, metaData []models.SchemaMetaData) (string, error)
	ListSchemas(ctx context.Context) ([]models.SchemaDetail, error)
	GetSchema(ctx context.Context, schemaName string, filter models.Filter) (models.GetSchemaResponse, error)
	DropSchema(ctx context.Context, schemaName string) error
}

type RecordService interface {
	CreateRecord(ctx context.Context, req *pb.CreateRecordRequest) ([]models.SchemaMetaData, []map[string]interface{}, error)
	GetRecords(ctx context.Context, schemaName string, filter models.Filter) ([]map[string]interface{}, error)
	GetRecord(ctx context.Context, schemaName string, recordId int) ([]models.SchemaMetaData, map[string]interface{}, error)
	DeleteRecord(ctx context.Context, schemaName string, recordId int) error
	UpdateRecord(ctx context.Context, req *pb.UpdateRecordRequest) ([]models.SchemaMetaData, map[string]interface{}, error)
}

type InternalService interface {
	FetchMetaData(ctx context.Context, schemaName string, isSystemName bool) ([]models.SchemaMetaData, error)
}
