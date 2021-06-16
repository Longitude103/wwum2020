package parcels

import "github.com/Longitude103/wwum2020/database"

// dryWaterBalanceWSPP is a method that calculates the dryland parcel water balance for RO and DP using the parcel
// information and a slice of database.CoeffCrop and the adjustment factor within it. This method is the WSPP approach
// to the calculation
func (p *Parcel) dryWaterBalanceWSPP(cCrops []database.CoeffCrop) error {
	adjFactor := adjustmentFactor(p, cCrops, database.DryET)

	// Create ETMAXDRYAdj and set to ETMAXDRY * Adjustment
	var (
		etDryAdj [12]float64
	)
	_, EttoRO, _, _, err := database.FilterCCDryLand(cCrops, p.CoeffZone, int(p.Crop1.Int64))
	if err != nil {
		return err
	}

	for i := 0; i < 12; i++ {
		etDryAdj[i] = p.DryEt[i] * adjFactor

		// RO3 = ETMAXDRY - ETMAXDryAdj * DryETtoRO
		p.Ro[i] += ((p.DryEt[i] - etDryAdj[i]) * EttoRO) * p.Area / 12

		// DP3 = (ETMAXDRY - ETMAXADJU) - RO3
		p.Dp[i] += ((p.DryEt[i] - etDryAdj[i]) - (1 - EttoRO)) * p.Area / 12
	}

	return nil
}
