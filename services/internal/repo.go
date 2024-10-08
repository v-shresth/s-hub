package internal

import (
	"cms/models"
	"cms/utils"
	"cms/utils/constants"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type repo struct {
	db  *gorm.DB
	log utils.Logger
}

func newRepo(db *gorm.DB, log utils.Logger) *repo {
	return &repo{
		db:  db,
		log: log,
	}
}

func (r *repo) fetchMetaDataInfo(ctx context.Context, schemaName string, isSystemName bool) ([]models.SchemaMetaData, error) {
	matchWithTableField := "system_schema_name"
	if !isSystemName {
		matchWithTableField = "display_schema_name"
	}
	var metaData []models.SchemaMetaData
	err := r.db.Debug().WithContext(ctx).Table(constants.MetadataSchema).
		Select("id, system_schema_name, display_schema_name, system_field_name, display_field_name, display_field_type").
		Where(fmt.Sprintf("%s = ? AND deleted_at IS NULL", matchWithTableField), schemaName).Scan(&metaData).Error
	if err != nil {
		return metaData, fmt.Errorf("error checking table existence: %v", err)
	}

	if len(metaData) > 0 {
		for _, data := range constants.DefaultMetaDataColumns {
			metaData = append(metaData, models.SchemaMetaData{
				SystemSchemaName:  metaData[0].SystemSchemaName,
				DisplaySchemaName: metaData[0].DisplaySchemaName,
				DisplayFieldType:  data.DisplayFieldType,
				SystemFieldName:   data.SystemFieldName,
				DisplayFieldName:  data.DisplayFieldName,
			})
		}
	}

	return metaData, nil
}
