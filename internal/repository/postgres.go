package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable        = "users"
	invoicesTable     = "invoices"
	checksTable       = "checks"
	transactionsTable = "transactions"

	timeFormat = "2006-01-02" //yyyy-mm-dd
)

type Config struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DatabaseName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
