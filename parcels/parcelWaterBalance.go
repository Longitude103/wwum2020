package parcels

import (
	"errors"
	"fmt"
	"math"
)

// waterBalanceWSPP method takes all the parcel information (SW delivery and GW Pumping) and creates a water balance to
// determine the amount of Runoff and Deep Percolation that occurs off of each parcel and sets those values within the
// parcel struct. This uses the methodology that is within the WSPP program.
func (p *Parcel) waterBalanceWSPP(verbose bool) error {
	girFactor, fsl := setGirFact(p.AppEff)
	if verbose {
		fmt.Printf("GIR: %g, fsl: %g\n", girFactor, fsl)
	}

	totalNir := sumAnnual(p.Nir)
	if totalNir <= 0 {
		return errors.New("total nir cannot be zero")
	}

	appWAT, _, pslIrr := setAppWat(p.SWDel, p.Pump, fsl)
	roDpWt, err := setRoDpWt(p.Ro, p.Dp)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Println("AppWat:", appWAT)
		fmt.Println("pslIrr:", pslIrr)
		fmt.Println("RoDpWt:", roDpWt)
	}

	ro1, dp1 := setInitialRoDp(p.Ro, p.Dp, 1, 1)
	if verbose {
		fmt.Println("RO1:", ro1)
		fmt.Println("DP1:", dp1)
	}

	gainApWat, gainPsl, gainIrrEt, gainDryEt := setPreGain(p.Et, p.DryEt, appWAT, pslIrr)
	cIR := math.Max(gainIrrEt-gainDryEt, 0.0001)
	gIR := totalNir * girFactor
	if verbose {
		fmt.Printf("gainApWat: %g, gainPsl: %g, gainIrrEt: %g, gainDryEt: %g\n", gainApWat, gainPsl, gainIrrEt, gainDryEt)
		fmt.Printf("cIR: %g, gIR: %g\n", cIR, gIR)
	}

	eTGain, err := setEtGain(cIR, gainPsl, gIR, gainApWat, p.AppEff, gainIrrEt, gainDryEt)
	if err != nil {
		return err
	}

	distGain, err := distEtGain(eTGain, pslIrr, p.Et, p.DryEt)
	if err != nil {
		return err
	}
	etBase := setEtBase(pslIrr, p.Et, p.DryEt)
	et := setET(etBase, distGain)
	deltaET := setDeltaET(et, 0.95)

	if verbose {
		fmt.Println("etGain:", eTGain)
		fmt.Println("distGain:", distGain)
		fmt.Println("EtBase:", etBase)
		fmt.Println("ET:", et)
		fmt.Println("deltaET:", deltaET)
	}

	ro3, dp3 := distDeltaET(deltaET, roDpWt)
	if verbose {
		fmt.Println("RO3:", ro3)
		fmt.Println("DP3:", dp3)
	}

	ro2, dp2 := excessIrrReturnFlow(pslIrr, distGain, roDpWt)
	if verbose {
		fmt.Println("RO2:", ro2)
		fmt.Println("DP2:", dp2)
	}

	p.Ro = sumReturnFlows(ro1, ro2, ro3)
	p.Dp = sumReturnFlows(dp1, dp2, dp3)

	return nil
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
func setRoDpWt(ro [12]float64, dp [12]float64) ([12]float64, error) {
	wt := [12]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5} // always the same in DB, Runoff Deep Perc weight

	for i := 0; i < 12; i++ {
		if ro[i]+dp[i] > 0 {
			wt[i] = math.Min(math.Max(ro[i]/(ro[i]+dp[i]), 0.2), 0.8)
		}
	}

	return wt, nil
}

// setInitialRoDp is a function to set the initial run off (Ro2) and Deep Perc (Dp2) from irrigation in the model of zero and handle the
// condition where water was applied but no nir was calculated so that all the water goes back to Ro and DP.
func setInitialRoDp(csRo [12]float64, csDp [12]float64, adjRo float64, adjDp float64) (ro [12]float64, dp [12]float64) {
	for i := 0; i < 12; i++ {
		ro[i] = csRo[i] * adjRo
		dp[i] = csDp[i] * adjDp
	}

	return
}

// setPreGain is a function to set some total variables if there is a presence of ETGain where irrEt > DryEt. This sums the
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
func setEtGain(cIR float64, psl float64, gir float64, appWat float64, eff float64, irrEt float64, dryEt float64) (gain float64, err error) {
	if gir == 0 {
		return 0, errors.New("gir cannot be zero in setEtGain")
	}

	beta := cIR / gir
	if psl < gir {
		gain = math.Max(math.Min(cIR*(1-math.Pow(1-psl/gir, 1/beta)), appWat*eff), 0)
	} else {
		gain = irrEt - dryEt
	}

	return gain, nil
}

// distEtGain distributes the ET Gain by the monthly gain listed by post surface loss water, and if there are any
// remaining, it apportions it again to months without PSL but with ET differences.
func distEtGain(etGain float64, psl [12]float64, etIrr [12]float64, etDry [12]float64) (distEtGain [12]float64, err error) {
	// three criteria, leftover falls to next distribution
	var (
		totalDiff       float64 // total difference when psl > 0
		totalNonPslDiff float64 // total difference when psl <= 0
		totalEtIrr      float64
		diffMonths      []int   // months when psl > 0
		nonPslMonths    []int   // months when psl <= 0
		remainGain      float64 // gain after first distribution
	)

	remainGain = etGain

	// find total difference
	for i := 0; i < 12; i++ {
		if psl[i] > 0 {
			totalDiff += etIrr[i] - etDry[i]
			totalEtIrr += etIrr[i]
			diffMonths = append(diffMonths, i)
		} else {
			if etIrr[i]-etDry[i] > 0 {
				totalNonPslDiff += etIrr[i] - etDry[i]
				nonPslMonths = append(nonPslMonths, i)
			}
		}
	}

	if len(diffMonths) > 0 {
		if totalDiff <= 0 {
			return distEtGain, errors.New("totalDiff cannot be zero in distEtGain")
		}
		for _, v := range diffMonths {
			distEtGain[v] = math.Min(etGain*(etIrr[v]-etDry[v])/totalDiff, psl[v])
			remainGain -= distEtGain[v]
		}
	}

	if remainGain > 0.001 {
		if len(nonPslMonths) > 0 {
			if totalNonPslDiff <= 0 {
				return distEtGain, errors.New("totalNonPslDiff cannot be zero in distEtGain")
			}

			// psl = 0 but ETirr > ETdry || remainingGain left
			for _, v := range nonPslMonths {
				distEtGain[v] += remainGain * (etIrr[v] - etDry[v]) / totalNonPslDiff
			}
		} else {
			// no other diff months, add back by weight of ETirr
			if totalEtIrr <= 0 {
				return distEtGain, errors.New("totalEtIrr cannot be zero in distEtGain")
			}

			for _, v := range diffMonths {
				distEtGain[v] += remainGain * (etIrr[v] / totalEtIrr)
			}
		}
	}

	return distEtGain, nil
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
	for i := range etBase {
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

// distDeltaET is a function that returns the run off and deep percolation of the delta ET
func distDeltaET(deltaET [12]float64, roDpWt [12]float64) (ro [12]float64, dp [12]float64) {
	for i, v := range deltaET {
		ro[i] = v * roDpWt[i]
		dp[i] = v - ro[i]
	}

	return
}

// excessIrrReturnFlow is a function that returns the excess irrigation return flows using the post surface loss irrigation
// with et gain and then distributed to ro and dp using the weighted values
func excessIrrReturnFlow(psl [12]float64, distEtGain [12]float64, roDpWt [12]float64) (ro [12]float64, dp [12]float64) {
	for i, v := range psl {
		if v > 0 { // protect against zero psl
			ro[i] = (v - distEtGain[i]) * roDpWt[i]
			dp[i] = (v - distEtGain[i]) - ro[i]
		}
	}

	return
}

// sumReturnFlows is a function to sum the three return flow sub variables into one.
func sumReturnFlows(v1 [12]float64, v2 [12]float64, v3 [12]float64) (sumValues [12]float64) {
	for i := 0; i < 12; i++ {
		sumValues[i] = v1[i] + v2[i] + v3[i]
	}

	return
}
