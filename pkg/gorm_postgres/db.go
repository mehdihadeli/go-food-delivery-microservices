package gorm_postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"dbName"`
	SSLMode  bool   `yaml:"sslMode"`
	Password string `yaml:"password"`
}

func NewGorm(cfg *Config) (*gorm.DB, error) {

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
