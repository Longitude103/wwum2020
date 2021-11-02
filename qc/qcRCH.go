package qc

import (
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
	//fmt.Printf("q: %+v\n", q)
	if err := q.rechargeBalance(); err != nil {
		return err
	}

	if err := q.rechargeGeoJson(); err != nil {
		return err
	}

	return nil
}

func (q *QC) getNodeDataFromSqlite(node int) {
	//q.v.SlDb.Select()
}
