package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"go-service-template/config"
)

type DbConnection struct {
	Db               *sqlx.DB
	InsertTimeoutSec int
}

func NewDBConnection(config config.DbConfig) (*DbConnection, error) {
	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.Db)
	db, err := sqlx.Open("postgres", cs)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	return &DbConnection{
		Db:               db,
		InsertTimeoutSec: config.InsertTimeoutSec,
	}, nil
}
