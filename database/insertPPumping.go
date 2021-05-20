package database

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Pumping struct {
	ParcelID int       `db:"parcelID"`
	Nrd      string    `db:"nrd"`
	Dt       time.Time `db:"dt"`
	Pump     float64   `db:"pump"`
}

type PPDB struct {
	sql    *sqlx.DB
	stmt   *sqlx.Stmt
	buffer []Pumping
}

func ParcelPumpDB(sqlDB *sqlx.DB) (*PPDB, error) {
	insertSql := `INSERT INTO parcelPumping (parcelID, nrd, dt, pump) VALUES (?, ?, ?, ?)`

	stmt, err := sqlDB.Preparex(insertSql)
	if err != nil {
		return nil, err
	}

	db := PPDB{sql: sqlDB, stmt: stmt, buffer: make([]Pumping, 0, 1024)}

	return &db, nil
}

func (db *PPDB) Add(pPump Pumping) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("parcel pumping buffer is full")
	}

	db.buffer = append(db.buffer, pPump)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return fmt.Errorf("unable to flush parcel pumping: %w", err)
		}
	}

	return nil
}

func (db *PPDB) Flush() error {
	tx, err := db.sql.Beginx()
	if err != nil {
		return err
	}

	for _, ppump := range db.buffer {
		if ppump.Pump > 0 {
			_, err := tx.Stmtx(db.stmt).Exec(ppump.ParcelID, ppump.Nrd, ppump.Dt, ppump.Pump)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}

	db.buffer = db.buffer[:0]
	return tx.Commit()
}

func (db *PPDB) Close() error {
	defer func() {
		_ = db.stmt.Close()
		_ = db.sql.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}
