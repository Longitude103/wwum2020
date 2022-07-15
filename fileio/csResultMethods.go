package fileio

import "errors"

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
		return result, errors.New("No station data provided")
	}

	for k, v := range stationData {
		for _, s := range stationSoils(v) {
			stationsBySoil := filterSoils(s, v)
			for _, c := range stationCrops(stationsBySoil) {
				stationsByCrop := filterCrops(c, stationsBySoil)
				for _, t := range stationTillage(stationsByCrop) {
					stationsByTillage := filterTillage(t, stationsByCrop)
					for _, i := range stationIrr(stationsByTillage) {
						stationsByIrr := filterIrr(i, stationsByTillage)

						totalVals := float64(len(stationsByIrr))
						if totalVals == 0 {
							break
						}

						var etSum, effSum, nirSum, dpSum, roSum, precSum float64

						for _, sr := range stationsByIrr {
							if sr.Yr >= AvgStart && sr.Yr <= AvgEnd {
								avg := sr.AverageAnnual()

								etSum += avg.Et
								effSum += avg.Eff_precip
								nirSum += avg.Nir
								dpSum += avg.Dp
								roSum += avg.Ro
								precSum += avg.Precip
							}
						}

						mv := MonthlyValues{
							Et:         etSum / totalVals,
							Eff_precip: effSum / totalVals,
							Nir:        nirSum / totalVals,
							Dp:         dpSum / totalVals,
							Ro:         roSum / totalVals,
							Precip:     precSum / totalVals,
						}

						newSR := StationResults{
							Station:     k,
							Soil:        s,
							MonthlyData: []MonthlyValues{mv},
							Yr:          1952,
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
