package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

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
	conn, err := parsePostgresConnectionString(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	if err := conn.createDatabaseIfNotExists(); err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	if err := conn.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL: %s:%d, database: %s",
		conn.Host, conn.Port, conn.DatabaseName)
	return conn, nil
}

func parsePostgresConnectionString(connectionString string) (*PostgresConnection, error) {
	conn := &PostgresConnection{
		Host:         "localhost",
		Port:         5432,
		DatabaseName: "kineticafs",
		Username:     "postgres",
		Password:     "",
		SSLMode:      "disable",
	}

	if connectionString == "" {
		return conn, nil
	}

	u, err := url.Parse(connectionString)
	if err != nil {
		return nil, fmt.Errorf("invalid connection string format: %w", err)
	}

	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return nil, fmt.Errorf("connection string must start with postgres:// or postgresql://")
	}

	conn.Host = u.Hostname()
	if port := u.Port(); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			conn.Port = p
		}
	}

	if u.Path != "" && len(u.Path) > 1 {
		conn.DatabaseName = u.Path[1:] // Remove leading slash
	}

	if u.User != nil {
		conn.Username = u.User.Username()
		if password, ok := u.User.Password(); ok {
			conn.Password = password
		}
	}

	if sslmode := u.Query().Get("sslmode"); sslmode != "" {
		conn.SSLMode = sslmode
	}

	return conn, nil
}

func (pc *PostgresConnection) createDatabaseIfNotExists() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		pc.Host, pc.Port, pc.Username, pc.Password, pc.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer db.Close()

	var exists bool
	checkQuery := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)"
	if err := db.QueryRow(checkQuery, pc.DatabaseName).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		createQuery := fmt.Sprintf("CREATE DATABASE %s", pc.DatabaseName)
		if _, err := db.Exec(createQuery); err != nil {
			return fmt.Errorf("failed to create database %s: %w", pc.DatabaseName, err)
		}
		log.Printf("Database '%s' created", pc.DatabaseName)
	} else {
		log.Printf("Database '%s' already exists", pc.DatabaseName)
	}

	return nil
}

func (pc *PostgresConnection) Connect() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pc.Host, pc.Port, pc.Username, pc.Password, pc.DatabaseName, pc.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	pc.DB = db
	return nil
}

func (pc *PostgresConnection) HealthCheck() error {
	if pc.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	if err := pc.DB.Ping(); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	return nil
}

func (pc *PostgresConnection) Close() error {
	if pc.DB != nil {
		if err := pc.DB.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		log.Println("PostgreSQL connection closed")
	}
	return nil
}

func (pc *PostgresConnection) ExecuteQuery(query string, args ...interface{}) error {
	log.Printf("Executing PostgreSQL query: %s, args: %v", query, args)
	_, err := pc.DB.Exec(query, args...)
	return err
}

func (pc *PostgresConnection) ExecuteQueryRow(query string, args ...interface{}) *sql.Row {
	log.Printf("Executing PostgreSQL query row: %s, args: %v", query, args)
	return pc.DB.QueryRow(query, args...)
}

func (pc *PostgresConnection) ExecuteQueryRows(query string, args ...interface{}) (*sql.Rows, error) {
	log.Printf("Executing PostgreSQL query rows: %s, args: %v", query, args)
	return pc.DB.Query(query, args...)
}

func (pc *PostgresConnection) BeginTransaction() (*sql.Tx, error) {
	return pc.DB.Begin()
}

func (pc *PostgresConnection) GetDB() *sql.DB {
	return pc.DB
}
