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
	p.Dialector = postgres.Open(p.GetDSN())
	client, err := gorm.Open(p.Dialector, &gorm.Config{})
	return client, err
}
