package qc

import (
	"github.com/Longitude103/wwum2020/database"
)

type QC struct {
	v         *database.Setup
	fileName  string
	Graph     bool
	Year      int
	GJson     bool
	Monthly   bool
	WellGJson bool
	grid      int
}

type Option func(*QC)

func WithGraph() Option {
	return func(q *QC) { q.Graph = true }
}

func WithYear(year int) Option {
	return func(q *QC) { q.Year = year }
}

func WithGJson() Option {
	return func(q *QC) { q.GJson = true }
}

func WithMonthly() Option {
	return func(q *QC) { q.Monthly = true }
}

func WithWellGJson() Option {
	return func(q *QC) { q.WellGJson = true }
}

func NewQC(v *database.Setup, fileName string, options ...Option) *QC {
	q := &QC{v: v, fileName: fileName, Year: 1997}
	for _, option := range options {
		option(q)
	}

	q.grid, _ = database.GetGrid(q.v.SlDb)

	return q
}

func StartQcRMain(q *QC) error {
	//fmt.Printf("q: %+v\n", q)
	if err := q.rechargeBalance(); err != nil {
		return err
	}

	if q.GJson {
		if err := q.rechargeGeoJson(); err != nil {
			return err
		}
	}

	if q.WellGJson {
		if err := q.WellPumpingGJson(); err != nil {
			return err
		}
	}

	return nil
}

func (q *QC) getNodeDataFromSqlite(node int) {
	//q.v.SlDb.Select()
}
