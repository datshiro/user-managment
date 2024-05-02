package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDbc *gorm.DB

type PostgresConnection struct {
	config *DBConfig
}

func (p *PostgresConnection) Connect() (*gorm.DB, error) {
	dsn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
	p.config.dialector = postgres.Open(
    fmt.Sprintf(dsn, p.config.Host, p.config.User, p.config.Password, p.config.DB, p.config.Port, p.config.SSLMode, p.config.TimeZone))
	client, err := p.config.Connect()
	return client, err
}
