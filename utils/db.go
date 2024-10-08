package utils

import (
	"cms/models"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(log Logger, config Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.GetDbUrl()))
	if err != nil {
		log.WithError(err).Fatal("Error creating database connection")
	}

	err = db.AutoMigrate(&models.SchemaMetaData{})
	if err != nil {
		log.WithError(err).Fatal("Error applying migrations")
	}

	return db
}
