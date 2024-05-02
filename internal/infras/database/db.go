package database

import "gorm.io/gorm"


type DatabaseConnection interface {
    Connect() (*gorm.DB, error)
}


