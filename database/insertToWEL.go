package database

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

// TODO: create documentation for this file

type WelResult struct {
	Wellid   int       `db:"well_id"`
	Node     int       `db:"cell_node"`
	Dt       time.Time `db:"dt"`
	FileType int       `db:"file_type"`
	Result   float64   `db:"result"`
}

type WelDB struct {
	sql    *sqlx.DB
	stmt   *sqlx.Stmt
	buffer []WelResult
}

func ResultsWelDB(sqlDB *sqlx.DB) (*WelDB, error) {
	insertSQL := `INSERT INTO results (well_id, cell_node, dt, file_type, result) VALUES (?, ?, ?, ?, ?)`

	stmt, err := sqlDB.Preparex(insertSQL)
	if err != nil {
		return nil, err
	}

	db := WelDB{
		sql:    sqlDB,
		stmt:   stmt,
		buffer: make([]WelResult, 0, 1024),
	}

	return &db, nil
}

func (db *WelDB) Add(value WelResult) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("conveyance loss buffer is full")
	}

	db.buffer = append(db.buffer, value)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return fmt.Errorf("unable to flush conveyance loss: %w", err)
		}
	}

	return nil
}

func (db *WelDB) Flush() error {
	tx, err := db.sql.Beginx()
	if err != nil {
		return err
	}

	for _, cl := range db.buffer {
		_, err := tx.Stmtx(db.stmt).Exec(cl.Node, cl.Dt, cl.FileType, cl.Result)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	db.buffer = db.buffer[:0]
	return tx.Commit()
}

func (db *WelDB) Close() error {
	defer func() {
		_ = db.stmt.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}
