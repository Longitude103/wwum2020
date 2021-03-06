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

// GetAppEfficiency is a function that returns the application efficiency for use in calculations throughout the app.
func GetAppEfficiency(pgDB *sqlx.DB) (efficiencies []Efficiency) {
	query := `select * from rswb.appeff;`

	err := pgDB.Select(&efficiencies, query)
	if err != nil {
		fmt.Println("Cannot Get Efficiencies from DB", err)
	}

	return efficiencies
}
