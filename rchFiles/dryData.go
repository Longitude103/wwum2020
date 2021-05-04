package rchFiles

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// DryCell is a stuct that holds the data for each cell and the parcel data asscociated with it including crops and
// the amount of crop that is included.
type DryCell struct {
	CellId   int             `db:"cellid"`
	CellArea float64         `db:"c_area"`
	DryArea  float64         `db:"d_area"`
	Crop1    sql.NullString  `db:"crop1"`
	Crop2    sql.NullString  `db:"crop2"`
	Crop3    sql.NullString  `db:"crop3"`
	Crop4    sql.NullString  `db:"crop4"`
	Crop1Cov sql.NullFloat64 `db:"crop1_cov"`
	Crop2Cov sql.NullFloat64 `db:"crop2_cov"`
	Crop3Cov sql.NullFloat64 `db:"crop3_cov"`
	Crop4Cov sql.NullFloat64 `db:"crop4_cov"`
}

// GetDryCells is a function that returns a struct of cells with parcels that are within it including the crops and acres
// within each cell. If there are more then one parcel within a cell, the cell will be listed multiple times.
func GetDryCells(db *sqlx.DB, yr int) []DryCell {
	query := fmt.Sprintf(`SELECT tfg_cellid cellid, st_area(c.geom)/43560 c_area, st_area(st_intersection(c.geom, d.geom))/43560 d_area,
       crop1, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov from public.act_cells c
    inner join (SELECT crop1, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, geom from np.t%d_dry UNION
                select crop1, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, geom from sp.t%d_dry) d
        on st_intersects(c.geom, d.geom);`, yr, yr)

	var dryCells []DryCell
	err := db.Select(&dryCells, query)
	if err != nil {
		fmt.Println("Error", err)
	}

	//fmt.Println("First Dryland Cell:")
	//fmt.Println(dryCells[0])

	return dryCells
}
