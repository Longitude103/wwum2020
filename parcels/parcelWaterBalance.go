package parcels

import (
	"fmt"
	"github.com/Longitude103/wwum2020/database"
	"math"
)

// waterBalanceWSPP method takes all the parcel information (SW delivery and GW Pumping) and creates a water balance to
// determine the amount of Runoff and Deep Percolation that occurs off of each parcel and sets those values within the
// parcel struct. This uses the methodology that is within the WSPP program.
func (p *Parcel) waterBalanceWSPP(cCrops []database.CoeffCrop) error {
	// TODO: Check on this for an infinite loop, seems stuck.

	// determine GIRFactor and Fsl_co
	// GIRFactor = Gross irrigation Requirement factor
	girFactor, fsl := setGirFact(p.AppEff)
	fmt.Printf("AppEff: %g, girFactor: %g, fsl: %g\n", p.AppEff, girFactor, fsl)

	totalNir := sumAnnual(p.Nir)
	appWAT, _, pslIrr := setAppWat(p.SWDel, p.Pump, fsl)
	roDpWt := setRoDpWt(p.Ro, p.Dp)

	ro2, dp2 := setInitialRoDp(p.Nir, appWAT, pslIrr, roDpWt)

	gainApWat, gainPsl, gainIrrEt, gainDryEt := setPreGain(p.Et, p.DryEt, appWAT, pslIrr)

	cIR := math.Max(gainIrrEt-gainDryEt, 0.0001)
	gIR := totalNir * girFactor

	eTGain := setEtGain(cIR, gainPsl, gIR, gainApWat, p.AppEff, gainIrrEt, gainDryEt)

	_, _, et1, _, etIrrGain := DistEtCOGain(eTGain, pslIrr, p.Et, p.DryEt)
	fmt.Printf("et1: %g\n", et1)
	fmt.Printf("etIrrGain: %g\n", etIrrGain)

	et := [12]float64{}
	det := [12]float64{}
	for i := 0; i < 12; i++ {
		et[i] = et1[i] * adjustmentFactor(p, cCrops, database.NirEt)
		det[i] = et1[i] * adjustmentFactor(p, cCrops, database.NirEt)
	}

	var (
		ro3, dp3 [12]float64
	)
	for i := 0; i < 12; i++ {

		ro3[i] = det[i] * roDpWt[i]
		//dp3[i] = ro3[i] * (1-roDpWt[i])/roDpWt[i]
		dp3[i] = det[i] - ro3[i]

		if p.Nir[i] > 0 {
			ro2[i] = math.Max((pslIrr[i]-etIrrGain[i])*roDpWt[i], -(p.Ro[i] + ro3[i]))
			dp2[i] = math.Max((pslIrr[i]-etIrrGain[i])*(1-roDpWt[i]), -(p.Dp[i] + dp3[i]))
		}
	}

	// Line 762 in WSPP
	for i := 0; i < 12; i++ {
		p.Ro[i] += ro2[i] + ro3[i]
		p.Dp[i] += dp2[i] + dp3[i]
	}

	return nil
}

// DistEtCOGain is a function that is called by the waterBalanceWSPP method that distributes the ETGain by month.
func DistEtCOGain(totalEtGain float64, pslIrr [12]float64, irrEt [12]float64, dryEt [12]float64) (ro [12]float64,
	dp [12]float64, eT1 [12]float64, eTBase [12]float64, etIrrGain [12]float64) {
	// file at Line 2077

	var (
		flag                                  [12]bool // change flag to all true
		eTIrrGain                             [12]float64
		roDpwt                                [12]float64
		irrEtSum, dryEtSum, irrPSLSum, eTGain float64
	)

	eTGain = totalEtGain
	for i := 0; i < 12; i++ {
		if pslIrr[i] <= 0 {
			flag[i] = true
		}
	}

	b := true
	for b {
		irrEtSum = 0
		dryEtSum = 0

		for i := 0; i < 12; i++ {
			if !flag[i] {
				irrEtSum += irrEt[i]
				dryEtSum += dryEt[i]
			}
		}

		for i := 0; i < 12; i++ {
			if !flag[i] {
				eTIrrGain[i] = totalEtGain * (irrEt[i] - dryEt[i]) / (irrEtSum - dryEtSum)
				if eTIrrGain[i] > pslIrr[i] {
					eTIrrGain[i] = pslIrr[i]
					flag[i] = true
					eTGain -= pslIrr[i]
					break
				} else {
					b = false
				}
			}
		}
	}

	for i := 0; i < 12; i++ {
		eTGain += eTIrrGain[i]
		flag[i] = false
	}

	if eTGain < totalEtGain {
		eTGain = totalEtGain - eTGain

		if eTGain > 0.00001 {
			for i := 0; i < 12; i++ {
				if pslIrr[i] <= 0 {
					flag[i] = true
				}
			}

			b := true
			for b {
				irrPSLSum = 0
				for i := 0; i < 12; i++ {
					if !flag[i] {
						irrPSLSum += pslIrr[i]
					}
				}

				for i := 0; i < 12; i++ {
					eTIrrGain[i] = eTGain * pslIrr[i] / irrPSLSum
					if eTIrrGain[i] > pslIrr[i] {
						eTIrrGain[i] = pslIrr[i]
						flag[i] = true
						eTGain -= pslIrr[i]
						break
					} else {
						b = false
					}
				}
			}
		}
	}

	for i := 0; i < 12; i++ {
		eTGain += eTIrrGain[i]
		flag[i] = false
	}

	if eTGain < totalEtGain {
		eTGain = totalEtGain - eTGain

		if eTGain > 0.00001 {
			for i := 0; i < 12; i++ {
				if pslIrr[i] <= 0 {
					flag[i] = true
				}
			}

			irrEtSum = 0
			for i := 0; i < 12; i++ {
				if !flag[i] {
					irrEtSum += irrEt[i]
				}
			}

			for i := 0; i < 12; i++ {
				if !flag[i] {
					eTIrrGain[i] = eTGain * irrEt[i] / irrEtSum
					ro[i] = math.Max((pslIrr[i]-eTIrrGain[i])*roDpwt[i], 0)
					dp[i] = ro[i] * (1 - roDpwt[i]) / roDpwt[i]
				}
			}
		}
	}

	for i := 0; i < 12; i++ {
		eT1[i] = eTIrrGain[i] + irrEt[i]
		eTBase[i] = irrEt[i]
	}

	return
}

// setGirFact is a function that sets the gross irrigation factor for the WSPP program and the fraction of surface loss
// amount depending on the efficiency passed in. It returns two float64 values used within the app.
func setGirFact(eff float64) (gir float64, fsl float64) {
	if eff >= 0.75 {
		gir = 1 / 0.95
		fsl = 0.02
	} else {
		gir = 1 / 0.75
		fsl = 0.05
	}

	return
}

// sumAnnual is a function to get the annual amount from a 12 month array of float64s, it returns a float64 total
func sumAnnual(data [12]float64) (total float64) {
	for _, d := range data {
		total += d
	}

	return
}

// setAppWat is a function that sets the applied water (appWat), surface loss of water (sL) and
// post surface loss of water (pSL) for each month of the parcel. It takes in surface water applied (sw),
// ground water applied (gw) and fraction of surface loss (fsl) and returns three arrays of monthly results.
func setAppWat(sw [12]float64, gw [12]float64, fsl float64) (appWat [12]float64, sL [12]float64, pSL [12]float64) {
	for i := 0; i < 12; i++ {
		appWat[i] = sw[i] + gw[i]
		sL[i] = appWat[i] * fsl
		pSL[i] = appWat[i] - sL[i]
	}

	return
}

// setRoDpWt sets the weight of the runoff to deep percolation values for each month but is bound by 0.2 to 0.8. It returns
// a monthly array of percent that is runoff of the total of runoff + deep percolation; has a default value of 0.5.
func setRoDpWt(ro [12]float64, dp [12]float64) [12]float64 {
	wt := [12]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5} // always the same in DB, Runoff Deep Perc weight

	for i := 0; i < 12; i++ {
		if ro[i]+dp[i] > 0 {
			wt[i] = math.Min(math.Max(ro[i]/(ro[i]+dp[i]), 0.2), 0.8)
		}
	}

	return wt
}

// setInitialRoDp is a function to set the initial run off (Ro2) and Deep Perc (Dp2) from irrigation in the model of zero and handle the
// condition where water was applied but no nir was calculated so that all the water goes back to Ro and DP.
func setInitialRoDp(nir [12]float64, appWat [12]float64, pslIrr [12]float64, RoDpWt [12]float64) (ro [12]float64, dp [12]float64) {
	for i := 0; i < 12; i++ {
		if nir[i] <= 0 && appWat[i] > 0 {
			ro[i] = pslIrr[i] * RoDpWt[i]
			dp[i] = pslIrr[i] - ro[i]
		}
	}

	return
}

// setPreGain is a function to set some total variables if there is a presance of ETGain where irrEt > DryEt. This sums the
// irrigated ET, Dry ET, Applied Water, and Post Surface Loss Water during those months where the condition is met.
func setPreGain(et [12]float64, dryEt [12]float64, appWat [12]float64, pslIrr [12]float64) (gainApWat float64, gainPsl float64, gainIrrEt float64, gainDryEt float64) {
	for i := 0; i < 12; i++ {
		if et[i] > dryEt[i] {
			// it's etgain
			gainIrrEt += et[i]
			gainDryEt += dryEt[i]
			gainApWat += appWat[i]
			gainPsl += pslIrr[i]
		}
	}

	return
}

// setEtGain sets the annual gain for the parcel using a diminishing returns production function. Returns the amount of gain
func setEtGain(cIR float64, psl float64, gir float64, appWat float64, eff float64, irrEt float64, dryEt float64) (gain float64) {
	beta := cIR / gir
	if psl < gir {
		gain = math.Max(math.Min(cIR*(1-math.Pow(1-psl/gir, 1/beta)), appWat*eff), 0)
	} else {
		gain = irrEt - dryEt
	}

	return
}
