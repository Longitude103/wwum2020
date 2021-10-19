package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"time"
)

// GetSqlite gets or sets the connection string for the sqlite3 results database.
//  It takes no args, and returns the db object for the connection
func GetSqlite(logger *zap.SugaredLogger, mDesc string) (*sqlx.DB, error) {
	dbName := fmt.Sprintf("./results%s.sqlite", time.Now().Format(time.RFC3339))
	logger.Infof("created sqlite results db named: %s", dbName)
	db, err := sqlx.Open("sqlite3", dbName)
	if err != nil {
		logger.Errorf("Error in creating SQLite DB: %s", err)
		return nil, err
	}

	// initializes the database with a results table and file_keys table
	err = InitializeDb(db, logger, mDesc)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectSqlite(fileName string) (*sqlx.DB, error) {
	dbPath := fmt.Sprintf("./%s", fileName)

	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
