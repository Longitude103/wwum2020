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

func PgConnx() *sqlx.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	return db
}
