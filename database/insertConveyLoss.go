package database

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type CLResult struct {
	Node     int       `db:"cell_node"`
	Dt       time.Time `db:"dt"`
	FileType int       `db:"file_type"`
	Result   float64   `db:"result"`
}

type CLDB struct {
	sql    *sqlx.DB
	stmt   *sqlx.Stmt
	buffer []CLResult
}

func ConveyLossDB(sqlDB *sqlx.DB) (*CLDB, error) {
	insertSQL := `INSERT INTO results (cell_node, dt, file_type, result) VALUES (?, ?, ?, ?)`

	stmt, err := sqlDB.Preparex(insertSQL)
	if err != nil {
		return nil, err
	}

	db := CLDB{
		sql:    sqlDB,
		stmt:   stmt,
		buffer: make([]CLResult, 0, 1024),
	}

	return &db, nil
}

func (db *CLDB) Add(conveyLoss CLResult) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("conveyance loss buffer is full")
	}

	db.buffer = append(db.buffer, conveyLoss)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return fmt.Errorf("unable to flush conveyance loss: %w", err)
		}
	}

	return nil
}

func (db *CLDB) Flush() error {
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

func (db *CLDB) Close() error {
	defer func() {
		_ = db.stmt.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}