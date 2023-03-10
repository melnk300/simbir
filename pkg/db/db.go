package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db"`
}

func InitDB(cfg Config) (*sqlx.DB, error) {
	return sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"user=%s password=%s dbname=%s sslmode=disable",
			cfg.User,
			cfg.Password,
			cfg.DBName,
		),
	)
}
