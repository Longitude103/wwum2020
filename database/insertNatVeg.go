package database

import (
	_ "database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type NPastCellStruct struct {
	Node int
	RO   [12]float64
	DP   [12]float64
	Yr   int
}

type NvDB struct {
	sql    *sqlx.DB
	stmt   *sqlx.Stmt
	buffer []NPastCellStruct
}

// NatVegDB is a function that returns the NvDB struct setup to insert values into the results database for the Natural
// Vegetation functions.
func NatVegDB(sqlDB *sqlx.DB) (*NvDB, error) {
	insertSQL := `INSERT INTO results (cell_node, dt, file_type, result) VALUES (?, ?, ?, ?)`

	stmt, err := sqlDB.Preparex(insertSQL)
	if err != nil {
		return nil, err
	}

	db := NvDB{
		sql:    sqlDB,
		stmt:   stmt,
		buffer: make([]NPastCellStruct, 0, 1024),
	}

	return &db, nil
}

// Add is a method to NvDB to add a record that will get saved.
func (db *NvDB) Add(c NPastCellStruct) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("nir buffer is full")
	}

	db.buffer = append(db.buffer, c)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return fmt.Errorf("unable to flush pnir: %w", err)
		}
	}

	return nil
}

// Flush is a method to flush the buffer of records and save them to the db.
func (db *NvDB) Flush() error {
	tx, err := db.sql.Beginx()
	if err != nil {
		return err
	}

	for _, cell := range db.buffer {
		for i := 0; i < 12; i++ {
			dt := time.Date(cell.Yr, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
			if cell.RO[i]+cell.RO[i] > 0 {
				_, err := tx.Stmtx(db.stmt).Exec(cell.Node, dt, 102, cell.RO[i]+cell.DP[i])
				if err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}
	}

	db.buffer = db.buffer[:0]
	return tx.Commit()
}

// Close is a method to close the statement, but does not close the DB as it's used in other places around the App.
func (db *NvDB) Close() error {
	defer func() {
		_ = db.stmt.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}
