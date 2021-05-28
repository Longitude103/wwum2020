package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

type CoeffCrop struct {
	Zone         int     `db:"zone"`
	Crop         int     `db:"crop"`
	DryEtAdj     float64 `db:"dryetadj"`
	IrrEtAdj     float64 `db:"irretadj"`
	NirAdjFactor float64 `db:"niradjfactor"`
	FslGW        float64 `db:"fslgw"`
	DryEtToRo    float64 `db:"dryettoro"`
	FslSW        float64 `db:"fslsw"`
	PerToRch     float64 `db:"pertorch"`
	DpAdj        float64 `db:"dpadj"`
	RoAdj        float64 `db:"roadj"`
}

// GetCoeffCrops is a function that calls the database to get the list of coefficients for each crop and zone in the model
// and returns a slice of CoeffCrop.
func GetCoeffCrops(pgDB *sqlx.DB) (CoeffCrops []CoeffCrop, err error) {
	query := `select * from rswb.coeffcrops;`

	if err := pgDB.Select(&CoeffCrops, query); err != nil {
		return nil, err
	}

	return CoeffCrops, nil
}

// FilterValues is a method that will return the values that you would use to filter a slice of CoeffCrop struct for a filter
// function
func (coeff CoeffCrop) FilterValues() (z, c int) {
	return coeff.Zone, coeff.Crop
}

func FilterCCDryLand(cSlice []CoeffCrop, z int, c int) (DryEtAdj float64, DryEtToRo float64, DpAdj float64,
	RoAdj float64, err error) {
	for _, v := range cSlice {
		vZ, vC := v.FilterValues()
		if vZ == z && vC == c {
			return v.DryEtAdj, v.DryEtToRo, v.DpAdj, v.RoAdj, nil
		}
	}

	return 0, 0, 0, 0, errors.New("crop not found in coefficient of crops")
}
