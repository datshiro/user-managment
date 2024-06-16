package database

import (
	"app/internal/interfaces/repositories/models"
	"log"

	"gorm.io/gorm"
)

type DatabaseConnection interface {
	Connect() (*gorm.DB, error)
}

func NewPostgresConnection(opts ...OptFunc) *gorm.DB {
	o := Opts{}
	for _, optFunc := range opts {
		optFunc(&o)
	}

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
