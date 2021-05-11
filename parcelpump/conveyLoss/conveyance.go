package conveyLoss

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func Conveyance(pgDB *sqlx.DB, sYear int, eYear int, excessFlow bool) {
	spRates := map[string]float64{"interstate": 0.4869, "highline": 0.2617, "lowline": 0.2513}
	_ = spRates

	canalCells := getCanalCells(pgDB)
	fmt.Println("First 10 Canal Cells")
	for _, v := range canalCells[:10] {
		fmt.Println(v)
	}

	diversions := getDiversions(pgDB, sYear, eYear, excessFlow)
	fmt.Println("First 10 Diversions")
	for _, v := range diversions[:10] {
		fmt.Println(v)
	}

}
