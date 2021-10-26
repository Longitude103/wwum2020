package qc

import (
	"fmt"
	"github.com/Longitude103/wwum2020/database"
)

type QC struct {
	v        database.Setup
	fileName string
	Graph    bool
	Year     int
	GJson    bool
	Monthly  bool
}

type Option func(*QC)

func WithGraph(graph bool) Option {
	return func(q *QC) { q.Graph = graph }
}

func WithYear(year int) Option {
	return func(q *QC) { q.Year = year }
}

func WithGJson(gj bool) Option {
	return func(q *QC) { q.GJson = gj }
}

func WithMonthly(mon bool) Option {
	return func(q *QC) { q.Monthly = mon }
}

func NewQC(v database.Setup, fileName string, options ...Option) *QC {
	q := &QC{v: v, fileName: fileName, Year: 1997}
	for _, option := range options {
		option(q)
	}

	return q
}

func QcRMain(q *QC) error {
	fmt.Printf("q: %+v\n", q)

	return nil
}

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

	// find a node with a canal inside it
	query := "select node, st_area(m.geom) / 43560 cell_area, soil_code, zone, coeff_zone, name, eff, clink_id" +
		"from public.model_cells m join sw.canals c on st_intersects(m.geom, c.geom)" +
		"where type_2 = 'Canal' order by random() limit 10;"

	if err := q.v.PgDb.Select(&CCResults, query); err != nil {
		return err
	}

	// query the results db about that node

	// graph the values, make a text file output

	return nil
}
