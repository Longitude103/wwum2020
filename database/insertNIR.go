package database

import (
	_ "database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type PNir struct {
	ParcelNo  int
	Nrd       string
	ParcelNIR [12]float64
	Year      int
}

type DB struct {
	sql    *sqlx.DB
	stmt   *sqlx.Stmt
	buffer []PNir
}

func PNirDB(sqlDB *sqlx.DB) (*DB, error) {
	insertSQL := `INSERT INTO parcelNIR (parcelID, nrd, dt, nir) VALUES ($1, $2, $3, $4)`

	stmt, err := sqlDB.Preparex(insertSQL)
	if err != nil {
		return nil, err
	}

	db := DB{
		sql:    sqlDB,
		stmt:   stmt,
		buffer: make([]PNir, 0, 1024),
	}

	return &db, nil
}

func (db *DB) Add(pNir PNir) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("nir buffer is full")
	}

	db.buffer = append(db.buffer, pNir)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return fmt.Errorf("unable to flush pnir: %w", err)
		}
	}

	return nil
}

func (db *DB) Flush() error {
	tx, err := db.sql.Beginx()
	if err != nil {
		return err
	}

	for _, pnir := range db.buffer {
		for i := 0; i < 12; i++ {
			if pnir.ParcelNIR[i] > 0 {
				dt := time.Date(pnir.Year, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
				_, err := tx.Stmtx(db.stmt).Exec(pnir.ParcelNo, pnir.Nrd, dt, pnir.ParcelNIR[i])
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

func (db *DB) Close() error {
	defer func() {
		_ = db.stmt.Close()
		_ = db.sql.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}
