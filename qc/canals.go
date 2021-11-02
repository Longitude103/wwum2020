package qc

type CanalCells struct {
	Node     int     `db:"node"`
	CArea    float64 `db:"cell_area"`
	SoilCode int     `db:"soil_code"`
	Zone     int     `db:"zone"`
	CZone    int     `db:"coeff_zone"`
	Name     string  `db:"name"`
	Eff      float64 `db:"eff"`
	CLink    int     `db:"clink_id"`
}

func nodeWithCanal(q *QC) error {
	var CCResults []CanalCells

	// find 10 cells that have a canal inside it, they will be random
	query := "select node, st_area(m.geom) / 43560 cell_area, soil_code, zone, coeff_zone, name, eff, clink_id" +
		"from public.model_cells m join sw.canals c on st_intersects(m.geom, c.geom)" +
		"where type_2 = 'Canal' order by random() limit 10;"

	if err := q.v.PgDb.Select(&CCResults, query); err != nil {
		return err
	}

	for _, c := range CCResults {
		_ = c
	}
	// query the results db about that node

	// graph the values, make a text file output

	return nil
}
