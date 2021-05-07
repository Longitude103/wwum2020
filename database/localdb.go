package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// GetSqlite gets or sets the connection string for the sqlite3 results database.
//  It takes no args, and returns the db object for the connection
func GetSqlite() *sqlx.DB {
	dbName := fmt.Sprintf("./results%s.sqlite", time.Now().Format(time.RFC3339))
	db, err := sqlx.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println("Error", err)
	}

	// initializes the database with a results table and file_keys table
	InitializeDb(db)

	return db
}
