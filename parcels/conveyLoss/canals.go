package conveyLoss

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CanalCell struct {
	CanalId    int             `db:"id"`
	CanalType  string          `db:"type_2"`
	DistId     int             `db:"district_id"`
	Eff        sql.NullFloat64 `db:"eff"`
	Node       int             `db:"node"`
	StLength   float64         `db:"st_length"`
	CFlag      int             `db:"c_flag"`
	DnrFact    sql.NullFloat64 `db:"dnr_fact"`
	SatFact    sql.NullFloat64 `db:"sat_fact"`
	UsgsFact   sql.NullFloat64 `db:"usgs_fact"`
	CLinkId    int             `db:"clink_id"`
	CanalEff   sql.NullFloat64 `db:"eff2"`
	LatCount   sql.NullInt64   `db:"latcount"`
	TotalLatLn sql.NullFloat64 `db:"tot_lat_ln"`
	TotalCanLn float64         `db:"tot_can_ln"`
}

type Canal struct {
	Id   int             `db:"id"`
	Name string          `db:"name"`
	Eff  float64         `db:"eff"`
	Area sql.NullFloat64 `db:"area"`
	Yr   int             `db:"yr"`
}

func getCanalCells(pgDb *sqlx.DB) []CanalCell {
	query := `-- noinspection SqlResolveForFile
	
	SELECT c.id, c.type_2, c.district_id, c.eff, a.node,
		   ST_Length(ST_Intersection(a.geom, c.geom)), c.c_flag, d.dnr_fact, s.sat_fact, u.usgs_fact, c.clink_id, c1.eff eff2,
		   c2.latcount, c3.tot_lat_ln, c4.tot_can_ln
	FROM public.model_cells a JOIN sw.canals c ON ST_intersects(c.geom, a.geom)
		JOIN(SELECT eff, id from sw.canals where type_2 = 'Canal') c1 on c1.id = c.clink_id
		LEFT JOIN (SELECT count(clink_id) as latcount, clink_id from sw.canals WHERE type_2 = 'Lateral' GROUP BY clink_id) c2 on c2.clink_id = c.clink_id
		LEFT JOIN (SELECT SUM(ST_Length(geom)) as tot_lat_ln, clink_id from sw.canals WHERE type_2 = 'Lateral' GROUP BY clink_id) c3 on c3.clink_id = c.clink_id
		LEFT JOIN (SELECT SUM(ST_Length(geom)) as tot_can_ln, clink_id FROM sw.canals WHERE type_2 = 'Canal' GROUP BY clink_id) c4 on c4.clink_id = c.clink_id
		LEFT OUTER JOIN sw.factors d on d.node = a.node AND c.c_flag = 1
		LEFT OUTER JOIN sw.factors s on s.node = a.node AND c.c_flag = 4
		LEFT OUTER JOIN sw.factors u on u.node = a.node AND c.c_flag = 2 WHERE c.id NOT IN (12,16,17,42,49,54,55,346,347,348,349,350,351,352,353,355) ORDER BY c.clink_id, a.node;`

	var canalCells []CanalCell
	err := pgDb.Select(&canalCells, query)
	if err != nil {
		fmt.Println("Error in getting cells for canals", err)
	}

	return canalCells
}

// getCanals returns a slice of Canal with the canal id, name, efficiency and total acres for that all the years that
// are listed.
func getCanals(pgDb *sqlx.DB, sYear int, eYear int) (canals []Canal) {
	for i := sYear; i < eYear+1; i++ {
		query := fmt.Sprintf(`select id, name, eff, area, %d yr from sw.canals left join (select sw_id, sum(st_area(geom) / 43560) area
				from np.t%d_irr where sw = true and sw_id is not null group by sw_id UNION ALL select sw_id, 
				sum(st_area(geom) / 43560) area from sp.t%d_irr where sw = true and sw_id is not null 
				group by sw_id) a on id = a.sw_id where type_2 = 'Canal' and eff is not null;`, i, i, i)

		err := pgDb.Select(&canals, query)
		if err != nil {
			fmt.Println("Error Getting Canal Data", err)
		}
	}

	return canals
}
