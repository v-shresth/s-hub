package record

import (
	"cms/models"
	"cms/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
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

func (r *repo) createRecord(
	ctx context.Context, schemaName string, records []map[string]interface{},
) ([]map[string]interface{}, error) {
	err := r.db.Debug().WithContext(ctx).Table(schemaName).Create(&records).Error
	if err != nil {
		r.log.WithError(err).Error("Failed to create records")
		return nil, err
	}

	return records, nil
}

func (r *repo) getRecords(
	ctx context.Context,
	schemaName string,
	filter models.Filter,
) ([]map[string]interface{}, error) {
	args := []interface{}{
		filter.PageSize,
		filter.PageSize * (filter.PageNumber - 1),
	}
	SQL := fmt.Sprintf(`SELECT * FROM %s WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;`, schemaName)

	var resp []map[string]interface{}
	err := r.db.Debug().WithContext(ctx).Raw(SQL, args...).Scan(&resp).Error
	if err != nil {
		r.log.WithError(err).Error("Error while fetching schema")
		return nil, err
	}

	return resp, nil
}

func (r *repo) getRecord(
	ctx context.Context,
	schemaName string,
	recordId int,
) (map[string]interface{}, error) {
	args := []interface{}{
		recordId,
	}
	SQL := fmt.Sprintf(`SELECT * FROM %s WHERE deleted_at IS NULL AND id = $1;`, schemaName)

	var resp map[string]interface{}
	err := r.db.Debug().WithContext(ctx).Raw(SQL, args...).First(&resp).Error
	if err != nil {
		r.log.WithError(err).Error("Error while fetching record")
		return nil, err
	}

	return resp, nil
}

func (r *repo) deleteRecord(
	ctx context.Context,
	schemaName string,
	recordId int,
) error {
	err := r.db.Debug().WithContext(ctx).Table(schemaName).Where("id = ? AND deleted_at IS NULL", recordId).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
	}).Error
	if err != nil {
		r.log.WithError(err).Error("unable to delete record")
		return err
	}

	return nil
}

func (r *repo) updateRecord(
	ctx context.Context,
	schemaName string,
	recordId int,
	updates map[string]interface{},
) error {
	err := r.db.Debug().WithContext(ctx).Table(schemaName).Where("id=?", recordId).Updates(updates).Error
	if err != nil {
		r.log.WithError(err).Error("Error while updating schema")
		return err
	}

	return nil
}
