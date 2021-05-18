package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CoeffCrop struct {
	Zone         int     `db:"zone"`
	Crop         int     `db:"crop"`
	DryEtAdj     float64 `db:"dryetadj"`
	IrrEtAdj     float64 `db:"irretadj"`
	NirAdjFactor float64 `db:"niradjfactor"`
	FslGW        float64 `db:"fslgw"`
	DryEtToro    float64 `db:"dryettoro"`
	FslSW        float64 `db:"fslsw"`
	PerToRch     float64 `db:"pertorch"`
	DpAdj        float64 `db:"dpadj"`
	RoAdj        int     `db:"roadj"`
}

func GetCoeffCrops(pgDB *sqlx.DB) (CoeffCrops []CoeffCrop) {
	query := `select * from rswb.coeffcrops;`

	err := pgDB.Select(&CoeffCrops, query)
	if err != nil {
		fmt.Println("Error in Getting CoeffCrops data", err)
	}

	return CoeffCrops
}
