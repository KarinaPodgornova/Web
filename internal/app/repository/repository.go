package repository

import (
    "fmt"
    "lab2/internal/app/ds"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Repository struct {
    db *gorm.DB
}

func New(dsn string) (*Repository, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Автомиграция
    err = db.AutoMigrate(&ds.Device{}, &ds.Application{}, &ds.ApplicationDevices{})
    if err != nil {
        return nil, fmt.Errorf("failed to migrate database: %w", err)
    }

    return &Repository{db: db}, nil
}