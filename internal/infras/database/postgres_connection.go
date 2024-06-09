package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDbc *gorm.DB

type PostgresConnection struct {
	Opts
}

func (p *PostgresConnection) Connect() (*gorm.DB, error) {
	dsn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
	p.dialector = postgres.Open(fmt.Sprintf(dsn, p.Host, p.User, p.Password, p.DB, p.Port, p.SSLMode, p.TimeZone))
	client, err := p.MakeConnect()
	return client, err
}
