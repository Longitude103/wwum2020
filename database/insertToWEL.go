package database

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

// welResult is a stuct for the final result to be saved to db that is a value per well, per month
type welResult struct {
	Wellid   int       `db:"well_id"`
	Node     int       `db:"cell_node"`
	Dt       time.Time `db:"dt"`
	FileType int       `db:"file_type"`
	Result   float64   `db:"result"`
}

// WelResult is a struct that is used to construct the well result per well, per year and the result is 12 month array
type WelResult struct {
	Wellid   int
	Node     int
	Yr       int
	FileType int
	Result   [12]float64
}

// WelDB is a struct for the sql database, statement and buffer that is used to save information in chunks
type WelDB struct {
	sql    *sqlx.DB
	stmt   *sqlx.Stmt
	buffer []welResult
}

// ResultsWelDB is a function that creates the WelDB struct and contains the SQL statement to insert the records, it
// also accepts a slice of welResult used for the buffer
func ResultsWelDB(sqlDB *sqlx.DB) (*WelDB, error) {
	insertSQL := `INSERT INTO wel_results (well_id, cell_node, dt, file_type, result) VALUES (?, ?, ?, ?, ?)`

	stmt, err := sqlDB.Preparex(insertSQL)
	if err != nil {
		return nil, err
	}

	db := WelDB{
		sql:    sqlDB,
		stmt:   stmt,
		buffer: make([]welResult, 0, 1024),
	}

	return &db, nil
}

// Add is a method of WelDB that adds a record to the buffer, but we accept a WelResult and create 12 welResult records
// and remove the zeros for space consideration. If the buffer is full it calls the Flush method.
func (db *WelDB) Add(value WelResult) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("conveyance loss buffer is full")
	}

	// take in WelResult and reformat to welResult and don't save zero result values
	for i, v := range value.Result {
		if v > 0 {
			db.buffer = append(db.buffer, welResult{Wellid: value.Wellid, Node: value.Node, Dt: time.Date(value.Yr,
				time.Month(i+1), 1, 0, 0, 0, 0, time.UTC), FileType: value.FileType, Result: v})
			if len(db.buffer) == cap(db.buffer) {
				if err := db.Flush(); err != nil {
					return fmt.Errorf("unable to flush conveyance loss: %w", err)
				}
			}
		}
	}

	return nil
}

// Flush is a method to empty the buffer by executing the SQL statement, inserting the records and then clears out the
// buffer and commits to the database.
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

// Close is a method to call the Flush method once more and then calls the close for the statement
func (db *WelDB) Close() error {
	defer func() {
		_ = db.stmt.Close()
	}()

	if err := db.Flush(); err != nil {
		return err
	}

	return nil
}

// AddPumping is a method to add more pumping to the WelResult by pumping and a welCount that is the number of wells
// that it should be divided by.
func (wr *WelResult) AddPumping(pump [12]float64, welCount float64) {
	for i, f := range pump {
		wr.Result[i] += f / welCount
	}
}
