package database

import (
	"database/sql"
	"fmt"
)

const (
	host     = "long103-wwum.clmtjoquajav.us-east-2.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "rQ!461k&Rk8J"
	dbname   = "wwum"
)

func PgConn() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// close database
	//defer db.Close()

	// check db
	err = db.Ping()
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("Connected to Postgres %s database.\n", dbname)
	return db
}
