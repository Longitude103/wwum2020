package database

import (
	"errors"
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

type Adjustment int

const (
	DryET Adjustment = 1
	IrrEt Adjustment = 2
	NirEt Adjustment = 3
)

// GetCoeffCrops is a function that calls the database to get the list of coefficients for each crop and zone in the model
// and returns a slice of CoeffCrop.
func GetCoeffCrops(v Setup) (CoeffCrops []CoeffCrop, err error) {
	query := `select * from rswb.coeffcrops;`

	if err := v.PgDb.Select(&CoeffCrops, query); err != nil {
		v.Logger.Errorf("Error Getting CoeffCrops data from Db: %s", err)
		return nil, err
	}

	return CoeffCrops, nil
}

// FilterValues is a method that will return the values that you would use to filter a slice of CoeffCrop struct for a filter
// function, returns the zone and the crop as integers
func (coeff CoeffCrop) FilterValues() (z, c int) {
	return coeff.Zone, coeff.Crop
}

// FilterCCDryLand is a function that returns the dryland values of the coefficient of crops table by giving it a slice
// of CoeffCrop and a zone and crop.
func FilterCCDryLand(cSlice []CoeffCrop, z int, c int) (DryEtAdj float64, DryEtToRo float64, DpAdj float64,
	RoAdj float64, err error) {
	for _, v := range cSlice {
		vZ, vC := v.FilterValues()
		if vZ == z && vC == c {
			return v.DryEtAdj, v.DryEtToRo, v.DpAdj, v.RoAdj, nil
		}
	}

	if c == 15 {
		for _, v := range cSlice {
			vZ, vC := v.FilterValues()
			if vZ == z && vC == 7 {
				return v.DryEtAdj, v.DryEtToRo, v.DpAdj, v.RoAdj, nil
			}
		}
	}

	return 0, 0, 0, 0, errors.New("crop not found in coefficient of crops")
}
