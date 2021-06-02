package parcels

import (
	"github.com/heath140/wwum2020/database"
	"math"
)

// waterBalance method takes all the parcel information (SW delivery and GW Pumping) and creates a water balance to
// determine the amount of Runoff and Deep Percolation that occurs off of each parcel and sets those values within the
// parcel struct. This uses the method that is within the WSPP program.
func (p *Parcel) waterBalanceWSPP(cCrops []database.CoeffCrop) error {
	// determine GIRFactor and Fsl_co
	// GIRFactor = Gross irrigation Requirement factor
	// girFactor := 1 / 0.75
	fsl := 0.05

	if p.AppEff >= 0.75 {
		//girFactor = 1 / 0.95
		fsl = 0.02
	}

	ro2 := [12]float64{}
	dp2 := [12]float64{}
	dap := [12]float64{}
	appWAT := [12]float64{}                                                           // total applied water
	sL := [12]float64{}                                                               // surface loss
	pslIrr := [12]float64{}                                                           // Post Surface Loss Irrigation water
	roDpWt := [12]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5} // always the same in DB, Runoff Deep Perc weight
	for i := 0; i < 12; i++ {
		appWAT[i] = p.Pump[i] + p.SWDel[i]

		if p.Ro[i]+p.Dp[i] > 0 {
			roDpWt[i] = math.Min(math.Max(p.Ro[i]/(p.Ro[i]+p.Dp[i]), 0.2), 0.8)
		}

		sL[i] = fsl * appWAT[i]
		pslIrr[i] = appWAT[i] - sL[i]

		// Applied water without needing it...
		if p.Nir[i] <= 0 && appWAT[i] > 0 {
			ro2[i] = pslIrr[i] * roDpWt[i]
			dp2[i] = pslIrr[i] - ro2[i]
		} else {
			dap[i] = p.Nir[i] / p.AppEff * (1 - adjustmentFactor(p, cCrops))
			// TODO: add irrigated ET here
			// TODO: go to parcelNIR and change to add in the dryland version of the crop and add it's ET

		}

	}

	// RO1irr and DP1irr is RO and DP adjust by the coeffcrops adjustment factor that is always 1 besides native veg handled there.

	return nil
}

// waterBalanceESC is a method that will be implemented for the escape modeling at a later date. It will not be used at this
// point but we needed a stub placeholder. It might need to be refactored to it's own file.
func (p *Parcel) waterBalanceESC() error {

	return nil
}
