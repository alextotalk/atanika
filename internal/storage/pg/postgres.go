package pg

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string
	Port     string
	Username string
	DBName   string
	Password string
	SSLMode  string
}

func New(cfg Config) (*sqlx.DB, error) {

	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)

	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		op := "storage.pg.New"
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
