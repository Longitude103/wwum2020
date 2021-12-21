package database

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Cellrc struct {
	node int `db:"node"`
	rw   int `db:"rw"`
	clm  int `db:"clm"`
}

type CellrcDB struct {
	sql    *sqlx.DB
	stmt   *sqlx.Stmt
	buffer []Cellrc
}

func CellRCDB(sqlDB *sqlx.DB) (*CellrcDB, error) {
	insertSQL := `INSERT INTO cellrc (node, rw, clm) VALUES (?, ?, ?)`

	stmt, err := sqlDB.Preparex(insertSQL)
	if err != nil {
		return nil, err
	}

	db := CellrcDB{
		sql:    sqlDB,
		stmt:   stmt,
		buffer: make([]Cellrc, 0, 1024),
	}

	return &db, nil
}

func (db *CellrcDB) Add(c Cellrc) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("Cellrc buffer is full")
	}

	db.buffer = append(db.buffer, c)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return fmt.Errorf("unable to flush Cellrc buffer: %w", err)
		}
	}

	return nil
}

func (db *CellrcDB) Flush() error {
	tx, err := db.sql.Beginx()
	if err != nil {
		return err
	}

	for _, cl := range db.buffer {
		_, err := tx.Stmtx(db.stmt).Exec(cl.node, cl.rw, cl.clm)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	db.buffer = db.buffer[:0]
	return tx.Commit()
}

func (db *CellrcDB) Close() error {
	defer func() {
		_ = db.stmt.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}
