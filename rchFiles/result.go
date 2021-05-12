package rchFiles

import "time"

type Result struct {
	Node     int
	Dt       time.Time
	FileType int
	Result   float64
}
