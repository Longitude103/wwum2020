package database

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type SWDelResult struct {
	CanalId   int       `db:"canalId"`
	Dt        time.Time `db:"dt"`
	DelAmount float64   `db:"delAmount"`
}

type SWDelDB struct {
	sql    *sqlx.DB
	stmt   *sqlx.Stmt
	buffer []SWDelResult
}

func SWDeliveryDB(sqlDB *sqlx.DB) (*SWDelDB, error) {
	insertSQL := `INSERT INTO swDelivery (canalId, dt, delAmount) VALUES (?, ?, ?)`

	stmt, err := sqlDB.Preparex(insertSQL)
	if err != nil {
		return nil, err
	}

	db := SWDelDB{
		sql:    sqlDB,
		stmt:   stmt,
		buffer: make([]SWDelResult, 0, 1024),
	}

	return &db, nil
}

func (db *SWDelDB) Add(delivery SWDelResult) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("conveyance loss buffer is full")
	}

	db.buffer = append(db.buffer, delivery)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return fmt.Errorf("unable to flush conveyance loss: %w", err)
		}
	}

	return nil
}

func (db *SWDelDB) Flush() error {
	tx, err := db.sql.Beginx()
	if err != nil {
		return err
	}

	for _, cl := range db.buffer {
		_, err := tx.Stmtx(db.stmt).Exec(cl.CanalId, cl.Dt, cl.DelAmount)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	db.buffer = db.buffer[:0]
	return tx.Commit()
}

func (db *SWDelDB) Close() error {
	defer func() {
		_ = db.stmt.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}
