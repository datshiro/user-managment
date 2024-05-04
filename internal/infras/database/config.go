package database

import (
	"app/internal/models"

	"gorm.io/gorm"
)

type DBConfig struct {
	DB        string         `env:"DB"`
	User      string         `env:"USER"`
	Password  string         `env:"PASSWORD"`
	Host      string         `env:"HOST"`
	Port      string         `env:"PORT"`
	SSLMode   string         `env:"SSL_MODE" default:"disable"`
	TimeZone  string         `env:"TIME_ZONE" default:"Asia/Shanghai"`
	dialector gorm.Dialector 
}

func (config *DBConfig) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(config.dialector, &gorm.Config{})
	return db, err
}

func (config *DBConfig) NewPostgresConnection() (*gorm.DB, error) {
	dbConnection := &PostgresConnection{config: config}
	var err error

	postgresDbc, err = dbConnection.Connect()
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	postgresDbc.AutoMigrate(&models.User{})
  // End of migration

	return postgresDbc, nil
}

func ClosePostgresConnection() error {
	db, err := postgresDbc.DB()
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
