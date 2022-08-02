package fileio

import (
	"errors"
)

func (sr StationResults) AverageAnnual() MonthlyValues {
	var mv MonthlyValues

	totalVals := float64(len(sr.MonthlyData))

	var etSum, effSum, nirSum, dpSum, roSum, precSum float64
	for _, m := range sr.MonthlyData {
		etSum += m.Et
		effSum += m.Eff_precip
		nirSum += m.Nir
		dpSum += m.Dp
		roSum += m.Ro
		precSum += m.Precip
	}

	mv.Et = etSum / totalVals
	mv.Eff_precip = effSum / totalVals
	mv.Nir = nirSum / totalVals
	mv.Dp = dpSum / totalVals
	mv.Ro = roSum / totalVals
	mv.Precip = precSum / totalVals

	return mv
}

func AverageStationResults(stationData map[string][]StationResults, AvgStart, AvgEnd int) (map[string][]StationResults, error) {
	// loop through each station and average the results into 1952 year

	result := make(map[string][]StationResults)

	if len(stationData) == 0 {
		return result, errors.New("no station data provided")
	}

	for k, v := range stationData {
		// fmt.Printf("Station: %s - Data: %+v\n", k, v)
		// fmt.Printf("Station: %s", k)
		// fmt.Println()
		for _, s := range stationSoils(v) {
			// fmt.Printf("Soils: %v\n", s)
			stationsBySoil := filterSoils(s, v)
			for _, c := range stationCrops(stationsBySoil) {
				// fmt.Printf("Crop: %v\n", c)
				stationsByCrop := filterCrops(c, stationsBySoil)
				for _, t := range stationTillage(stationsByCrop) {
					// fmt.Printf("Tillage: %v\n", t)
					stationsByTillage := filterTillage(t, stationsByCrop)
					for _, i := range stationIrr(stationsByTillage) {
						// fmt.Printf("Irrigation: %v\n", i)
						stationsByIrr := filterIrr(i, stationsByTillage)

						stationCount := float64(len(stationsByIrr))
						if stationCount == 0 {
							break
						}

						var etSum, effSum, nirSum, dpSum, roSum, precSum [12]float64
						for _, sr := range stationsByIrr {
							if sr.Yr >= AvgStart && sr.Yr <= AvgEnd {

								for m := 0; m < 12; m++ {
									etSum[m] += sr.MonthlyData[m].Et
									effSum[m] += sr.MonthlyData[m].Eff_precip
									nirSum[m] += sr.MonthlyData[m].Nir
									dpSum[m] += sr.MonthlyData[m].Dp
									roSum[m] += sr.MonthlyData[m].Ro
									precSum[m] += sr.MonthlyData[m].Precip
								}
							}
						}

						var avgMV []MonthlyValues
						for m := 0; m < 12; m++ {
							mv := MonthlyValues{
								Et:         etSum[m] / stationCount,
								Eff_precip: effSum[m] / stationCount,
								Nir:        nirSum[m] / stationCount,
								Dp:         dpSum[m] / stationCount,
								Ro:         roSum[m] / stationCount,
								Precip:     precSum[m] / stationCount,
							}

							avgMV = append(avgMV, mv)
						}

						for yr := 1893; yr < 1953; yr++ {
							newSR := StationResults{
								Station:     k,
								Soil:        s,
								MonthlyData: avgMV,
								Yr:          yr,
								Crop:        c,
								Tillage:     t,
								Irrigation:  i,
							}

							result[k] = append(result[k], newSR)
						}

					}

				}
			}
		}
	}
	return result, nil
}

func stationSoils(results []StationResults) (soils []int) {
	for _, s := range results {
		if !foundItem(soils, s.Soil) {
			soils = append(soils, s.Soil)
		}
	}

	return soils
}

func stationCrops(results []StationResults) (crops []int) {
	for _, s := range results {
		if !foundItem(crops, s.Crop) {
			crops = append(crops, s.Crop)
		}
	}

	return crops
}

func stationTillage(results []StationResults) (till []int) {
	for _, s := range results {
		if !foundItem(till, s.Tillage) {
			till = append(till, s.Tillage)
		}
	}

	return till
}

func stationIrr(results []StationResults) (irr []int) {
	for _, s := range results {
		if !foundItem(irr, s.Irrigation) {
			irr = append(irr, s.Irrigation)
		}
	}

	return irr
}

func filterSoils(soil int, stations []StationResults) (sr []StationResults) {
	for _, s := range stations {
		if s.Soil == soil {
			sr = append(sr, s)
		}
	}

	return
}

func filterCrops(crop int, stations []StationResults) (sr []StationResults) {
	for _, s := range stations {
		if s.Crop == crop {
			sr = append(sr, s)
		}
	}

	return
}

func filterTillage(till int, stations []StationResults) (sr []StationResults) {
	for _, s := range stations {
		if s.Tillage == till {
			sr = append(sr, s)
		}
	}

	return
}

func filterIrr(irr int, stations []StationResults) (sr []StationResults) {
	for _, s := range stations {
		if s.Irrigation == irr {
			sr = append(sr, s)
		}
	}

	return
}

func foundItem(items []int, newValue int) bool {
	for _, i := range items {
		if i == newValue {
			return true
		}
	}

	return false
}
