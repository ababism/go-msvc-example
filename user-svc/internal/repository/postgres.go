package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	usersTable    = "users"
	sonsTable     = "sons"
	albumsTable   = "albums"
	releasesTable = "releases"
)

type Config struct {
	Host     string
	Port     string
	Username string
	DBName   string
	Password string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}

//func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
//	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
//		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
//	if err != nil {
//		return nil, err
//	}
//
//	err = db.Ping()
//	if err != nil {
//		return nil, err
//	}
//
//	return db, nil
//}
