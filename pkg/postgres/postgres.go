package postgres

import (
	"fmt"
	"github.com/aberniando/xenshield/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// New -.
func New(config config.PG) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s",
		config.Username, config.Password, config.DBName, config.SSLMode, config.Host, config.Port)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
