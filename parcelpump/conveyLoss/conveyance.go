package conveyLoss

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func Conveyance(pgDB *sqlx.DB) {
	spRates := map[string]float64{"interstate": 0.4869, "highline": 0.2617, "lowline": 0.2513}
	_ = spRates

	canalCells := getCanalCells(pgDB)
	for _, v := range canalCells[:10] {
		fmt.Println(v)
	}

}
