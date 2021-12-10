package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Longitude103/wwum2020/Utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type lg interface {
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

// GetSqlite gets or sets the connection string for the sqlite3 results database.
//  It takes no args, and returns the db object for the connection
func GetSqlite(logger lg, mDesc string) (*sqlx.DB, error) {
	// wd, _ := os.Getwd()
	tn := time.Now()
	fmt.Printf("results%s-%d-%d\n", tn.Format(time.RFC3339)[:len(tn.Format(time.RFC3339))-15], tn.Hour(), tn.Minute())
	dbName := fmt.Sprintf("results%s-%d-%d.sqlite", tn.Format(time.RFC3339)[:len(tn.Format(time.RFC3339))-15], tn.Hour(), tn.Minute())
	oDir := fmt.Sprintf("results%s-%d-%d", tn.Format(time.RFC3339)[:len(tn.Format(time.RFC3339))-15], tn.Hour(), tn.Minute())
	path, err := Utils.MakeOutputDir(oDir)
	if err != nil {
		return nil, err
	}

	logger.Infof("created sqlite results db named: %s", dbName)
	db, err := sqlx.Open("sqlite3", filepath.FromSlash(filepath.Join(path, dbName)))
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
	wd, _ := os.Getwd()

	db, err := sqlx.Open("sqlite3", filepath.Join(wd, fileName))
	if err != nil {
		return nil, err
	}

	return db, nil
}
