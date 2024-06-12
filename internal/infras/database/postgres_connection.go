package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDbc *gorm.DB

type PostgresConnection struct {
	Opts
}

func (p *PostgresConnection) Connect() (*gorm.DB, error) {
	p.dialector = postgres.Open(p.GetDSN())
	client, err := gorm.Open(p.dialector, &gorm.Config{})
	return client, err
}
