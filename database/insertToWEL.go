package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// WelResult is a struct for the final result to be saved to db that is a value per well, per month
type WelResult struct {
	Wellid   int       `db:"well_id"`
	Node     int       `db:"cell_node"`
	Dt       time.Time `db:"dt"`
	FileType int       `db:"file_type"`
	Result   float64   `db:"result"`
}

// WelAnnualResult is a struct that is used to construct the well result and includes a 12-month array
type WelAnnualResult struct {
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
	buffer []WelResult
}

// ResultsWelDB is a function that creates the WelDB struct and contains the SQL statement to insert the records, it
// also accepts a slice of WelResult used for the buffer
func ResultsWelDB(sqlDB *sqlx.DB) (*WelDB, error) {
	insertSQL := `INSERT INTO wel_results (well_id, cell_node, dt, file_type, result) VALUES (?, ?, ?, ?, ?)`

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

// Add is a method of WelDB that adds a record to the buffer, but we accept a WelAnnualResult or WelResult struct and
// create 12 WelResult records for WelAnnualResult and remove the zeros. If there is a WelResult sent it will just save
// that value directly to the buffer. If the buffer is full it calls the Flush method.
func (db *WelDB) Add(value interface{}) error {
	if len(db.buffer) == cap(db.buffer) {
		return errors.New("WEL buffer is full")
	}

	switch value.(type) {
	case WelAnnualResult:
		wr := value.(WelAnnualResult)
		// take in WelAnnualResult and reformat to WelResult and don't save zero result values
		for i, v := range wr.Result {
			if v > 0 {
				db.buffer = append(db.buffer, WelResult{Wellid: wr.Wellid, Node: wr.Node, Dt: time.Date(wr.Yr,
					time.Month(i+1), 1, 0, 0, 0, 0, time.UTC), FileType: wr.FileType, Result: v})
				if len(db.buffer) == cap(db.buffer) {
					if err := db.Flush(); err != nil {
						return fmt.Errorf("unable to flush WEL: %w", err)
					}
				}
			}
		}
	case WelResult:
		wr := value.(WelResult)
		db.buffer = append(db.buffer, wr)
		if len(db.buffer) == cap(db.buffer) {
			if err := db.Flush(); err != nil {
				return fmt.Errorf("unable to flush WEL: %w", err)
			}
		}
	default:
		return fmt.Errorf("unable to determine the value struct type")
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
		_, err := tx.Stmtx(db.stmt).Exec(cl.Wellid, cl.Node, cl.Dt, cl.FileType, cl.Result)
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

// AddPumping is a method to add more pumping to the WelAnnualResult by pumping and a welCount that is the number of wells
// that it should be divided by.
func (wr *WelAnnualResult) AddPumping(pump [12]float64, welCount float64) {
	for i, f := range pump {
		wr.Result[i] += f / welCount
	}
}
