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

	// TODO: finish pass through of methods starting here
	_ = eTGain
	_ = ro2
	_ = dp2

	// TODO: ended at bullet 12 in example at equation 51

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

// roundTo rounds a float number to a specified number of decimal places.
func roundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
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
		if et[i] > dryEt[i] { // it's ET Gain
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

// distEtGain distributes the ET Gain by the monthly gain listed by post surface loss water, and if there are any
// remaining, it apportions it again to months without PSL but with ET differences.
func distEtGain(etGain float64, psl [12]float64, etIrr [12]float64, etDry [12]float64) (distEtGain [12]float64) {
	// three criteria, leftover falls to next distribution
	var (
		totalDiff       float64 // total difference when psl > 0
		totalNonPslDiff float64 // total difference when psl <= 0
		diffMonths      []int   // months when psl > 0
		nonPslMonths    []int   // months when psl <= 0
		remainGain      float64 // gain after first distribution
	)

	remainGain = etGain

	// find total difference
	for i := 0; i < 12; i++ {
		if psl[i] > 0 {
			totalDiff += etIrr[i] - etDry[i]
			diffMonths = append(diffMonths, i)
		} else {
			totalNonPslDiff += etIrr[i] - etDry[i]
			nonPslMonths = append(nonPslMonths, i)
		}
	}

	if len(diffMonths) > 0 {
		for _, v := range diffMonths {
			distEtGain[v] = math.Min(etGain*(etIrr[v]-etDry[v])/totalDiff, psl[v])
			remainGain -= distEtGain[v]
		}
	}

	if remainGain > 0.001 {
		// psl = 0 but ETirr > ETdry || remainingGain left
		for _, v := range nonPslMonths {
			distEtGain[v] += remainGain * (etIrr[v] - etDry[v]) / totalNonPslDiff
		}
	}

	return
	// psl > 0 && ETdry > ETirr ?? Strange, not covered

}

// setEtBase is a function that uses post surface loss irrigation to determine the etBase from etIrr and etDry and returns
// a monthly etBase value
func setEtBase(psl [12]float64, etIrr [12]float64, etDry [12]float64) (etBase [12]float64) {
	for i := 0; i < 12; i++ {
		if psl[i] <= 0 {
			etBase[i] = etIrr[i]
		} else {
			if etIrr[i] > etDry[i] {
				etBase[i] = etDry[i]
			} else {
				etBase[i] = etIrr[i]
			}
		}
	}

	return
}

// setET combines the distributed ET Gain with the base ET for a final ET Value
func setET(etBase [12]float64, distEtGain [12]float64) (et [12]float64) {
	for i, _ := range etBase {
		et[i] = etBase[i] + distEtGain[i]
	}

	return
}

// setDeltaET returns the monthly amount of adjustment of ET that is created from the adjustment factor application
func setDeltaET(et [12]float64, adjFactor float64) (deltaET [12]float64) {
	for i, v := range et {
		deltaET[i] = v * (1 - adjFactor)
	}

	return
}
