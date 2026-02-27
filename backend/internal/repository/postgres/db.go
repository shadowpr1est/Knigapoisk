package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shadowpr1est/Knigapoisk/config"
)

func NewDB(cfg config.DBConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host%=s port%=s user%=s password%=s dbname%=s sslmode%=s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	return db, nil
}
