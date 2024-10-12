package internal

import (
	"cms/clients"
	"cms/models"
	"cms/services"
	"context"
	"gorm.io/gorm"
)

type svc struct {
	log  clients.Logger
	repo *repo
}

func NewInternalService(log clients.Logger, userDb *gorm.DB, config clients.Config) services.InternalService {
	return &svc{
		log:  log,
		repo: newRepo(userDb, log, config),
	}
}

func (s *svc) FetchMetaData(
	ctx context.Context, schemaName string, isSystemName bool,
) ([]models.SchemaMetaData, error) {
	return s.repo.fetchMetaDataInfo(ctx, schemaName, isSystemName)
}
