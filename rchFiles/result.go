package rchFiles

import "time"

type Result struct {
	Node     int       `db:"cell_node"`
	Dt       time.Time `db:"dt"`
	FileType int       `db:"file_type"`
	Result   float64   `db:"result"`
}
