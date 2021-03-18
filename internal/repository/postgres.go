package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/nedostupno/seidhr/internal/config"
)

func NewPostgreDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Name, cfg.DB.Password, "disable"))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
