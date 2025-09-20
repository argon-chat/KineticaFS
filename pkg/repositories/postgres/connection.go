package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	DB           *sql.DB
	DatabaseName string
	Host         string
	Port         int
	Username     string
	Password     string
	SSLMode      string
}

func NewPostgresConnection(connectionString string) (*PostgresConnection, error) {
	return nil, nil
}

func (pc *PostgresConnection) Connect() error {
	return nil
}

func (pc *PostgresConnection) Close() error {
	return nil
}
