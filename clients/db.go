package clients

import (
	"cms/models"
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbUserMap map[uint]*gorm.DB

func init() {
	dbUserMap = make(map[uint]*gorm.DB)
}

func GetSystemDB(log Logger, config Config) *gorm.DB {
	systemDb, err := gorm.Open(postgres.Open(config.GetSystemDbUrl()))
	if err != nil {
		log.WithError(err).Fatal("Error creating database connection")
	}

	err = systemDb.Debug().AutoMigrate(&models.Users{}, &models.UserConnections{}, &models.Session{})
	if err != nil {
		log.WithError(err).Fatal("Error creating database connection")
	}

	return systemDb
}

func GetUserDb(ctx context.Context, userId uint, config Config, log Logger) (*gorm.DB, error) {
	if conn, ok := dbUserMap[userId]; ok {
		return conn, nil
	}

	userDb, err := gorm.Open(postgres.Open(config.GetUserDbUrl()))
	if err != nil {
		log.WithError(err).Fatal("Error creating database connection")
		return nil, err
	}

	schemaName := GetUserSchemaName(userId)

	createSchemaQuery := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)
	if err := userDb.Debug().WithContext(ctx).Exec(createSchemaQuery).Error; err != nil {
		log.WithError(err).Fatal("Error creating schema")
		return nil, err
	}

	setSearchPathQuery := fmt.Sprintf("SET search_path TO %s", schemaName)
	if err := userDb.Debug().WithContext(ctx).Exec(setSearchPathQuery).Error; err != nil {
		log.WithError(err).Fatal("Error setting search path")
		return nil, err
	}

	// Auto-migrate the tables within the schema
	err = userDb.Debug().WithContext(ctx).AutoMigrate(&models.SchemaMetaData{})
	if err != nil {
		log.WithError(err).Fatal("Error auto migrating tables")
		return nil, err
	}

	dbUserMap[userId] = userDb

	return userDb, nil
}

func GetUserSchemaName(userId uint) string {
	return fmt.Sprintf("shub_%d", userId)
}
