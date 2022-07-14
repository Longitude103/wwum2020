package database

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Note struct {
	Nt string `db:"note"`
}

type NoteDB struct {
	sql    *sqlx.DB
	stmt   *sqlx.Stmt
	buffer []Note
}

// ResultsNoteDB is a function that returns a db struct that includes a sql connection, a insert statement, and
// buffer to add notes to the output database.
func ResultsNoteDB(sqlDB *sqlx.DB) (*NoteDB, error) {
	insertSQL := `INSERT INTO results_notes (note) VALUES (?)`

	stmt, err := sqlDB.Preparex(insertSQL)
	if err != nil {
		return nil, err
	}

	db := NoteDB{
		sql:    sqlDB,
		stmt:   stmt,
		buffer: make([]Note, 0, 1024),
	}

	return &db, nil
}

// Add is a method on NoteDB that will add a text message to the output db results_notes table
func (db *NoteDB) Add(n Note) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("notes buffer is full")
	}

	db.buffer = append(db.buffer, n)
	if len(db.buffer) == cap(db.buffer) {
		if err := db.Flush(); err != nil {
			return fmt.Errorf("unable to flush notes buffer: %w", err)
		}
	}

	return nil
}

func (db *NoteDB) Flush() error {
	tx, err := db.sql.Beginx()
	if err != nil {
		return err
	}

	for _, cl := range db.buffer {
		_, err := tx.Stmtx(db.stmt).Exec(cl.Nt)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	db.buffer = db.buffer[:0]
	return tx.Commit()
}

func (db *NoteDB) Close() error {
	defer func() {
		_ = db.stmt.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}
