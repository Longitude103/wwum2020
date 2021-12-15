package database

import (
	"database/sql"
	"fmt"
	"math"
)

// DryCell is a struct that holds the data for each cell and the parcel data associated with it including crops and
// the amount of crop that is included.
type DryCell struct {
	Node     int             `db:"node"`
	Mtg      sql.NullFloat64 `db:"mtg"`
	CellArea float64         `db:"c_area"`
	DryArea  float64         `db:"d_area"`
	PId      int64           `db:"parcel_id"`
	Nrd      string          `db:"nrd"`
}

// GetDryCells is a function that returns a struct of cells with parcels that are within it including the crops and acres
// within each cell. If there are more then one parcel within a cell, the cell will be listed multiple times.
func GetDryCells(v *Setup, yr int) []DryCell {
	query := fmt.Sprintf(`SELECT node, mtg, st_area(c.geom)/43560 c_area, st_area(st_intersection(c.geom, d.geom))/43560 d_area, parcel_id, nrd
        from public.model_cells c inner join (SELECT parcel_id, 'np' nrd, geom from np.t%d_dry UNION select parcel_id, 'sp' nrd, 
		geom from sp.t%d_dry) d on st_intersects(c.geom, d.geom) where cell_type = %d;`, yr, yr, v.CellType())

	var dryCells []DryCell
	err := v.PgDb.Select(&dryCells, query)
	if err != nil {
		fmt.Println("Error", err)
	}

	return dryCells
}

func (d DryCell) GetLossFactor() float64 {
	if !d.Mtg.Valid || d.Mtg.Float64 == 0 {
		return 0.5
	}

	return math.Min(1-math.Exp(-0.02*d.Mtg.Float64), 1)
}
