package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Efficiency struct {
	Yr          int     `db:"yr"`
	AeFlood     float64 `db:"ae_flood"`
	AeSprinkler float64 `db:"ae_sprinkler"`
}

func GetAppEfficiency(pgDB *sqlx.DB) (efficiencies []Efficiency) {
	query := `-- noinspection SqlResolve
	
	select * from rswb.appeff;`

	err := pgDB.Select(&efficiencies, query)
	if err != nil {
		fmt.Println("Cannot Get Efficiencies from DB", err)
	}

	return efficiencies
}
