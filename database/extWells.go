package database

import (
	"errors"
	"time"
)

type ExtWell struct {
	Yr       int     `db:"yr"`
	Mnth     int     `db:"mnth"`
	FileType int     `db:"file_type"`
	Pumping  float64 `db:"pmp"`
	Node     int     `db:"node"`
}

// GetExternalWells is a function to query the external pumping from the database and returns a slice of ExtWell as well
// as includes handling the debug mode.
func GetExternalWells(v Setup) (extWells []ExtWell, err error) {
	const extQuery = "select yr, mnth, file_type, pmp, node from ext_pumping inner join model_cells mc on " +
		"st_contains(mc.geom, ext_pumping.geom);"

	if err := v.PgDb.Select(&extWells, extQuery); err != nil {
		return extWells, errors.New("error getting data from ext_pumping table from DB")
	}

	if v.AppDebug {
		return extWells[:50], nil
	}

	return extWells, nil
}

// Date returns the formatted date of the ExtWell struct with a 1 for the day and zero hour in UTC.
func (w *ExtWell) Date() time.Time {
	return time.Date(w.Yr, time.Month(w.Mnth), 1, 0, 0, 0, 0, time.UTC)
}
