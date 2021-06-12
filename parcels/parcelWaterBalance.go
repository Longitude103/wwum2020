package parcels

import (
	"github.com/heath140/wwum2020/database"
	"math"
)

// waterBalanceWSPP method takes all the parcel information (SW delivery and GW Pumping) and creates a water balance to
// determine the amount of Runoff and Deep Percolation that occurs off of each parcel and sets those values within the
// parcel struct. This uses the methodology that is within the WSPP program.
func (p *Parcel) waterBalanceWSPP(cCrops []database.CoeffCrop) error {
	// TODO: Check on this for an infinite loop, seems stuck.
	var (
		totalIrrEt, totalDryEt, totalAppWat, totalPSLIrr, totalNir float64
		ro2, dp2, dap, appWAT, sL, pslIrr                          [12]float64
		// ro2 = Runoff per month from Irrigation dp2= Deep Percolation per month from Irrigation
		//dap = Water delivery required to meet NIR appWAT = total applied water sL = surface loss pslIrr = Post Surface Loss Irrigation Water
	)
	roDpWt := [12]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5} // always the same in DB, Runoff Deep Perc weight

	// determine GIRFactor and Fsl_co
	// GIRFactor = Gross irrigation Requirement factor
	girFactor := 1 / 0.75
	fsl := 0.05

	if p.AppEff >= 0.75 {
		girFactor = 1 / 0.95
		fsl = 0.02
	}

	for i := 0; i < 12; i++ {
		totalNir += p.Nir[i]
		appWAT[i] = p.Pump[i] + p.SWDel[i]
		totalAppWat += appWAT[i]

		if p.Ro[i]+p.Dp[i] > 0 {
			roDpWt[i] = math.Min(math.Max(p.Ro[i]/(p.Ro[i]+p.Dp[i]), 0.2), 0.8)
		}

		sL[i] = fsl * appWAT[i]
		pslIrr[i] = appWAT[i] - sL[i]
		totalPSLIrr += pslIrr[i]

		// Applied water without needing it...
		if p.Nir[i] <= 0 && appWAT[i] > 0 {
			ro2[i] = pslIrr[i] * roDpWt[i]
			dp2[i] = pslIrr[i] - ro2[i]
		} else {
			dap[i] = p.Nir[i] / p.AppEff * (1 - adjustmentFactor(p, cCrops, database.NirEt))
			totalIrrEt += p.Nir[i]
			totalDryEt += p.DryEt[i]
		}
	}
	// RO1irr and DP1irr is RO and DP adjust by the coeffcrops adjustment factor that is always 1 besides native veg handled there.
	cIR := math.Max(totalIrrEt-totalDryEt, 0.0001)
	gIR := totalNir * girFactor
	beta := cIR / gIR
	totalEtGain := 0.0

	if totalPSLIrr < gIR {
		totalEtGain = math.Max(math.Min(cIR*(1-math.Pow(1-totalPSLIrr/gIR, 1/beta)), totalAppWat*p.AppEff), 0)
	} else {
		totalEtGain = math.Max(totalIrrEt-totalDryEt, 0.0)
	}

	_, _, et1, _, etIrrGain := DistEtCOGain(totalEtGain, pslIrr, p.Et, p.DryEt)

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
