package database

import (
	"fmt"
	"time"
)

type ExtRch struct {
	Node     int     `db:"node"`
	Size     float64 `db:"cell_size"`
	Yr       int     `db:"yr"`
	Mnth     int     `db:"mnth"`
	FileType int     `db:"file_type"`
	Rch      float64 `db:"rch"`
}

// GetExtRecharge is a function to return a slice of ExtRch filled with the ext_recharge data from the database and
// used in the model to create the recharge values for the external non-NRD areas.
func GetExtRecharge(v Setup) (eRch []ExtRch, err error) {
	query := fmt.Sprintf("select yr, mnth, file_type, rch, node, st_area(mc.geom)/43560 cell_size from "+
		"ext_recharge inner join model_cells mc on st_contains(mc.geom, ext_recharge.geom) where yr >= %d "+
		"and yr <= %d ;", v.SYear, v.EYear)

	if err := v.PgDb.Select(&eRch, query); err != nil {
		return eRch, err
	}

	return eRch, nil
}

// Date is a method to return a formatted date object for the ExtRch struct using its components and setting the day
// to 1 and 0 hours in UTC.
func (e *ExtRch) Date() time.Time {
	return time.Date(e.Yr, time.Month(e.Mnth), 1, 0, 0, 0, 0, time.UTC)
}
