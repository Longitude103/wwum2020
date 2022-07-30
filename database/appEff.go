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

// GetSSAppEfficiency is a function that returns the application efficiency for use in calculations throughout the
// Steady State portion of the app, it is fixed and doesn't use the database. The values are fixed and are made for
// years of 1893 through 1952.
func GetSSAppEfficiency() (efficiencies []Efficiency) {
	for yr := 1893; yr < 1953; yr++ {
		efficiencies = append(efficiencies, Efficiency{Yr: yr, AeFlood: 0.65, AeSprinkler: 0.7})
	}

	return efficiencies
}
