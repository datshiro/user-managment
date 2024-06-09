package database

import (
	"app/internal/models"
	"log"

	"gorm.io/gorm"
)

type DatabaseConnection interface {
	Connect() (*gorm.DB, error)
}

func NewPostgresConnection() *gorm.DB {
	o := loadOpts()

	dbConnection := &PostgresConnection{Opts: o}

	postgresDbc, err := dbConnection.Connect()
	if err != nil {
		log.Fatalf("Failed to connect database ; %v", err)
	}

	// Migrate the schema
	postgresDbc.AutoMigrate(&models.User{})
	// End of migration

	return postgresDbc
}

func ClosePostgresConnection() error {
	db, err := postgresDbc.DB()
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
