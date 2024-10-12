package record

import (
	"cms/clients"
	"cms/models"
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type repo struct {
	userDb *gorm.DB
	log    clients.Logger
	config clients.Config
}

func newRepo(userDb *gorm.DB, log clients.Logger, config clients.Config) *repo {
	return &repo{
		userDb: userDb,
		log:    log,
		config: config,
	}
}

func (r *repo) createRecord(
	ctx context.Context, schemaName string, records []map[string]interface{},
) ([]map[string]interface{}, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("no records to insert")
	}

	// Extract the columns from the first record
	columns := make([]string, 0, len(records[0]))
	for key := range records[0] {
		columns = append(columns, key)
	}

	// Start building the query
	var query strings.Builder
	query.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES ", schemaName, strings.Join(columns, ", ")))

	// Add placeholders for values
	values := make([]interface{}, 0)
	placeholders := make([]string, len(records))
	for i, record := range records {
		valuePlaceholders := make([]string, len(columns))
		for j, column := range columns {
			valuePlaceholders[j] = "?"
			values = append(values, record[column]) // Append the actual value for the placeholder
		}
		placeholders[i] = fmt.Sprintf("(%s)", strings.Join(valuePlaceholders, ", "))
	}

	// Append the value placeholders to the query
	query.WriteString(strings.Join(placeholders, ", "))

	// Add the RETURNING clause to get all values back
	query.WriteString(" RETURNING *")

	// Execute the query using GORM's Raw method
	var insertedRecords []map[string]interface{}
	err := r.userDb.Debug().WithContext(ctx).Raw(query.String(), values...).Scan(&insertedRecords).Error
	if err != nil {
		r.log.WithError(err).Error("Failed to create records")
		return nil, err
	}

	return insertedRecords, nil
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
	err := r.userDb.Debug().WithContext(ctx).Raw(SQL, args...).Scan(&resp).Error
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
	err := r.userDb.Debug().WithContext(ctx).Raw(SQL, args...).First(&resp).Error
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
	err := r.userDb.Debug().WithContext(ctx).Table(schemaName).Where("id = ? AND deleted_at IS NULL", recordId).Updates(map[string]interface{}{
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
	err := r.userDb.Debug().WithContext(ctx).Table(schemaName).Where("id=?", recordId).Updates(updates).Error
	if err != nil {
		r.log.WithError(err).Error("Error while updating schema")
		return err
	}

	return nil
}
