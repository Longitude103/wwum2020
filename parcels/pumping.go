package parcels

import (
	"github.com/heath140/wwum2020/database"
)

// estimatePumping is a method that is called on parcels that have metered == false and gw == true so that we can estimate
// the amount of pumping that was done at the parcel since a well is present, but not metered. Usually FA area before
// 2017 and other parcels that are not metered. It fills the Pump field of the Parcel struct
func (p *Parcel) estimatePumping(cCrops []database.CoeffCrop) {
	nirAdj := adjustmentFactor(p, cCrops)

	// get application efficiency
	var swAvailableCU, nirRemaining [12]float64
	if p.Sw.Bool == true {
		for i := 0; i < 12; i++ {
			swAvailableCU[i] = p.SWDel[i] * nirAdj * p.AppEff
		}
	}

	// set nirRemaining to nir - swAvailableCU if positive, then divide by AppEff to arrive at pumping
	for m := 0; m < 12; m++ {
		nirRemaining[m] = (p.Nir[m]*p.Area)/12.0 - swAvailableCU[m]
		if nirRemaining[m] > 0 {
			p.Pump[m] = nirRemaining[m] / p.AppEff
		}
	}
}

// adjustmentFactor function calculates the Parcel adjustment factor by weighting the crops and distribution of the
// crops in a Parcel by calling the nirFactor and then weighting it based on crop distribution
func adjustmentFactor(p *Parcel, cCrops []database.CoeffCrop) float64 {
	var c1, c2, c3, c4 float64
	c1 = nirFactor(cCrops, p.CoeffZone, int(p.Crop1.Int64)) * p.Crop1Cov.Float64

	if p.Crop2.Valid {
		c2 = nirFactor(cCrops, p.CoeffZone, int(p.Crop2.Int64)) * p.Crop2Cov.Float64
	}
	if p.Crop3.Valid {
		c3 = nirFactor(cCrops, p.CoeffZone, int(p.Crop3.Int64)) * p.Crop3Cov.Float64
	}

	if p.Crop4.Valid {
		c4 = nirFactor(cCrops, p.CoeffZone, int(p.Crop4.Int64)) * p.Crop4Cov.Float64
	}

	return c1 + c2 + c3 + c4
}

// nirFactor is a filter function that returns the NirAdjFactor from the CoeffCrop slice and limits it to the zone of the
// Parcel and the crop type.
func nirFactor(cCrops []database.CoeffCrop, zone int, crop int) (nf float64) {
	for _, v := range cCrops {
		if v.Zone == zone && v.Crop == crop {
			nf = v.NirAdjFactor
		}
	}

	return nf
}
