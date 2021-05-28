package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "long103-wwum.clmtjoquajav.us-east-2.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "rQ!461k&Rk8J"
	dbname   = "wwum"
)

// PgConnx is a function that returns the sql connection to the postgres database.
func PgConnx() (*sqlx.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
