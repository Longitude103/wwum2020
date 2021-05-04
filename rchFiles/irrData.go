package rchFiles

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type IrrCell struct {
	CellId   int             `db:"cellid"`
	CertNum  sql.NullString  `db:"cert_num"`
	CellArea float64         `db:"c_area"`
	IrrArea  float64         `db:"i_area"`
	Crop1    sql.NullString  `db:"crop1"`
	Crop2    sql.NullString  `db:"crop2"`
	Crop3    sql.NullString  `db:"crop3"`
	Crop4    sql.NullString  `db:"crop4"`
	IrrType  sql.NullString  `db:"irrig_type"`
	SwFac    sql.NullString  `db:"sw_fac"`
	ModelId  sql.NullString  `db:"model_id"`
	Crop1Cov sql.NullFloat64 `db:"crop1_cov"`
	Crop2Cov sql.NullFloat64 `db:"crop2_cov"`
	Crop3Cov sql.NullFloat64 `db:"crop3_cov"`
	Crop4Cov sql.NullFloat64 `db:"crop4_cov"`
	Sw       sql.NullBool    `db:"sw"`
	Gw       sql.NullBool    `db:"gw"`
}

func GetCellsIrr(db *sqlx.DB, yr int) {

	query := fmt.Sprintf(`SELECT tfg_cellid cellid, st_area(c.geom)/43560 c_area, st_area(st_intersection(c.geom, i.geom))/43560 i_area,
       crop1, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, sw, gw, irrig_type, sw_fac, cert_num, model_id
from public.act_cells c inner join (SELECT crop1, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, sw, gw, irrig_type, sw_fac, cert_num::varchar, model_id, geom from np.t%d_irr UNION
                                    SELECT crop1, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, sw, gw, irr_type as irrig_type, sw_fac, id as cert_num, null as model_id, geom from sp.t%d_irr) i on st_intersects(c.geom, i.geom);`, yr, yr)

	irrCells := []IrrCell{}
	err := db.Select(&irrCells, query)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("The First IRR Cell:")
	fmt.Println(irrCells[0])

}
