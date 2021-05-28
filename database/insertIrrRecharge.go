package database

import (
	_ "database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

// TODO: Merge this with Natural Flow and make generic struct and methods for Irr, Dry, Natural Veg

type IrrCellResult struct {
	Node int
	RO   [12]float64
	DP   [12]float64
	Yr   int
	File int
}

type IrrDB struct {
	sql    *sqlx.DB
	stmt   *sqlx.Stmt
	buffer []IrrCellResult
}

// IrrResultDB is a function that returns the IrrDB struct setup to insert values into the results database for the Irrigated
// parcel functions.
func IrrResultDB(sqlDB *sqlx.DB) (*IrrDB, error) {
	insertSQL := `INSERT INTO results (cell_node, dt, file_type, result) VALUES (?, ?, ?, ?)`

	stmt, err := sqlDB.Preparex(insertSQL)
	if err != nil {
		return nil, err
	}

	db := IrrDB{
		sql:    sqlDB,
		stmt:   stmt,
		buffer: make([]IrrCellResult, 0, 1024),
	}

	return &db, nil
}

// Add is a method to NvDB to add a record that will get saved.
func (db *IrrDB) Add(i IrrCellResult) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("nir buffer is full")
	}

	db.buffer = append(db.buffer, i)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return fmt.Errorf("unable to flush pnir: %w", err)
		}
	}

	return nil
}

// Flush is a method to flush the buffer of records and save them to the db.
func (db *IrrDB) Flush() error {
	tx, err := db.sql.Beginx()
	if err != nil {
		return err
	}

	for _, cell := range db.buffer {
		for i := 0; i < 12; i++ {
			dt := time.Date(cell.Yr, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
			if cell.RO[i]+cell.RO[i] > 0 {
				_, err := tx.Stmtx(db.stmt).Exec(cell.Node, dt, cell.File, cell.RO[i]+cell.DP[i])
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
func (db *IrrDB) Close() error {
	defer func() {
		_ = db.stmt.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}
