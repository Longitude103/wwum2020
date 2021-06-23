package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PgConnx is a function that returns the sql connection to the postgres database.
func PgConnx(myEnv map[string]string) (*sqlx.DB, error) {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", myEnv["host"], myEnv["port"], myEnv["user"], myEnv["password"], myEnv["dbname"])

	db, err := sqlx.Connect("postgres", psqlConn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
