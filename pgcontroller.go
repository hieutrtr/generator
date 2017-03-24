package postgresql_generator

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// PGCtrl contain information of PG Connection
type PGCtrl struct {
	conn *sql.DB
}
type Config struct {
	pgUser    string
	pgPass    string
	pgDB      string
	pgHost    string
	pgTimeOut int
}

func (c *Config) parse() error {
	c.pgUser = os.Getenv("PG_USER")
	c.pgPass = os.Getenv("PG_PASS")
	c.pgDB = os.Getenv("PG_DB")
	c.pgHost = os.Getenv("PG_HOST")
	if c.pgHost == "" {
		return errors.New("Missing PG_HOST")
	}
	if c.pgUser == "" {
		return errors.New("Missing PG_USER")
	}
	if c.pgPass == "" {
		return errors.New("Missing PG_PASS")
	}
	if c.pgDB == "" {
		return errors.New("Missing PG_DB")
	}
	return nil
}

// NewPG new PostgreSQL connection
// Need export
func NewPG() (*PGCtrl, error) {
	return &PGCtrl{}, nil
}

func (ctrl *PGCtrl) connect(cfg *Config) error {
	var err error
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.pgUser, cfg.pgPass, cfg.pgHost, cfg.pgDB)
	ctrl.conn, err = sql.Open("postgres", connStr)
	return err
}

// Execute executing of query string
// Query string should be validated
func (ctrl *PGCtrl) Execute(q string) error {
	_, err := ctrl.conn.Exec(q)
	return err
}
