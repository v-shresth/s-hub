package internal

import (
	"cms/models"
	"cms/services"
	"cms/utils"
	"context"
	"gorm.io/gorm"
)

type svc struct {
	log  utils.Logger
	repo *repo
}

func NewInternalService(log utils.Logger, db *gorm.DB) services.InternalService {
	return &svc{
		log:  log,
		repo: newRepo(db, log),
	}
}

func (s *svc) FetchMetaData(
	ctx context.Context, schemaName string, isSystemName bool,
) ([]models.SchemaMetaData, error) {
	return s.repo.fetchMetaDataInfo(ctx, schemaName, isSystemName)
}
