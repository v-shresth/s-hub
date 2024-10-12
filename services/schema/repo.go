package schema

import (
	"cms/clients"
	"cms/models"
	"cms/utils"
	"cms/utils/constants"
	"context"
	"fmt"
	"gorm.io/gorm"
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

func (r *repo) createSchemaMetaData(
	ctx context.Context, tx *gorm.DB, metaData []models.SchemaMetaData,
) error {
	err := tx.Debug().WithContext(ctx).Table(constants.MetadataSchema).Create(&metaData).Error
	if err != nil {
		r.log.WithError(err).Error("Error while creating schema meta data")
		return err
	}
	return nil
}

func (r *repo) createSchema(
	ctx context.Context, tx *gorm.DB, schema models.Schema,
) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
										(
											id                   SERIAL PRIMARY KEY,
											created_at           timestamp default now(),
    										updated_at           timestamp default now(),
											deleted_at          timestamp default null,
										`, schema.SchemaName)

	// Loop through fields and append them to query
	for _, col := range schema.Fields {
		query += fmt.Sprintf("%s %s, ", col.Name, col.Type)
	}

	// Remove the last comma and space, then close the parentheses
	query = query[:len(query)-2] + ");"

	// Execute the query
	err := tx.Debug().WithContext(ctx).Exec(query).Error
	if err != nil {
		r.log.WithError(err).Error("Error while creating schema")
		return err
	}
	return nil
}

func (r *repo) withTransaction(f func(db *gorm.DB) error) error {
	return utils.WithTransaction(r.userDb, func(tx *gorm.DB) error {
		return f(tx)
	})
}

func (r *repo) listSchemas(
	ctx context.Context,
) ([]models.SchemaDetail, error) {
	SQL := fmt.Sprintf(`WITH table_counts AS (SELECT display_schema_name,
												 COUNT(DISTINCT system_field_name) + 4 AS total_fields
										  FROM %s WHERE  deleted_at IS NULL
										  GROUP BY display_schema_name)
					SELECT display_schema_name schema_name,
						   total_fields no_of_fields,
						   (SELECT COUNT(DISTINCT system_schema_name)
							FROM %s WHERE deleted_at IS NULL) AS total_schemas
					FROM table_counts`, constants.MetadataSchema, constants.MetadataSchema)

	var schemaDetails []models.SchemaDetail
	err := r.userDb.Debug().WithContext(ctx).Raw(SQL).Scan(&schemaDetails).Error
	if err != nil {
		return nil, fmt.Errorf("error querying table and field counts: %v", err)
	}

	return schemaDetails, nil
}

func (r *repo) markSchemaArchived(
	ctx context.Context, tx *gorm.DB, schemaName string,
) error {
	err := tx.Debug().WithContext(ctx).Table(constants.MetadataSchema).Where("display_schema_name=?", schemaName).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
	}).Error
	if err != nil {
		r.log.WithError(err).Error("Error while marking schema archived")
		return err
	}

	return nil
}

func (r *repo) dropSchema(
	ctx context.Context,
	db *gorm.DB,
	schemaName string,
) error {
	SQL := fmt.Sprintf("DROP TABLE IF EXISTS %s;", schemaName)
	err := db.Debug().WithContext(ctx).Exec(SQL).Error
	if err != nil {
		r.log.WithError(err).Error("Error while dropping schema")
		return err
	}

	return nil
}
