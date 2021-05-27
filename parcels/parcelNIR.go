package parcels

import (
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
)

// parcelNIR is a method that adds the NIR, RO, and DP for each parcel from the CSResults and weather station data.
// It produces an intermediate results table of NIR in local sqlite for review and adds three maps to the parcel struct
func (p *Parcel) parcelNIR(pNirDB *database.DB, Year int, wStations []database.WeatherStation,
	csResults map[string][]fileio.StationResults, it IrrType) error {
	var parcelNIR, parcelRo, parcelDp [12]float64

	dist := database.Distances(p, wStations)
	for _, st := range dist {
		var annData []fileio.StationResults
		for _, data := range csResults[st.Station] {
			if data.Yr == Year && data.Soil == p.SoilCode && data.Irrigation == int(it) {
				annData = append(annData, data)
			}
		}

		var cropsNir, cropsRo, cropsDp [4][12]float64 // 4 crops X 12 months
		var cropCov [4]float64                        // crop_coverage
		if p.Crop1.Valid {
			cropsNir[0], cropsRo[0], cropsDp[0] = crop(p.Crop1.Int64, annData)
			cropCov[0] = p.Crop1Cov.Float64
		}

		if p.Crop2.Valid {
			cropsNir[1], cropsRo[1], cropsDp[1] = crop(p.Crop2.Int64, annData)
			cropCov[1] = p.Crop2Cov.Float64
		}

		if p.Crop3.Valid {
			cropsNir[2], cropsRo[2], cropsDp[2] = crop(p.Crop3.Int64, annData)
			cropCov[2] = p.Crop3Cov.Float64
		}

		if p.Crop4.Valid {
			cropsNir[3], cropsRo[3], cropsDp[3] = crop(p.Crop4.Int64, annData)
			cropCov[3] = p.Crop4Cov.Float64
		}

		// weight the crops based on crop_cov and weather station weight
		parcelNIR = pValues(parcelNIR, cropsNir, cropCov, st.Weight, p.Area)
		parcelRo = pValues(parcelRo, cropsRo, cropCov, st.Weight, p.Area)
		parcelDp = pValues(parcelDp, cropsDp, cropCov, st.Weight, p.Area)
	}
	//fmt.Printf("Weighted Parcel ID: %d, NIR is: %v\n", p.ParcelNo, parcelNIR)

	// save parcelNIR to sqlite only for irrigated
	if it == Irrigated {
		err := pNirDB.Add(database.PNir{ParcelNo: p.ParcelNo, Nrd: p.Nrd, ParcelNIR: parcelNIR, Year: Year, IrrType: int(it)})
		if err != nil {
			return err
		}
	}

	p.Nir = parcelNIR
	p.Ro = parcelRo
	p.Dp = parcelDp

	return nil
}

// crop function filters the results to the integer crop that is included and returns the NIR, RunOff and Deep Percolation from those
// crops as three arrays.
func crop(c int64, aData []fileio.StationResults) (nir [12]float64, ro [12]float64, dp [12]float64) {
	var data fileio.StationResults

	for _, d := range aData {
		if int64(d.Crop) == c {
			data = d
		}
	}

	for i, monthly := range data.MonthlyData {
		nir[i] = monthly.Nir
		ro[i] = monthly.Ro
		dp[i] = monthly.Dp
	}

	return nir, ro, dp
}

// pValues creates the parcel nir, ro, or dp by multiplying the cropNir with crop coverage by station weight to
// return a nir, ro, or dp portion for that parcel.
func pValues(parcelValues [12]float64, cropsNIR [4][12]float64, cropCov [4]float64, stWeight float64, area float64) (values [12]float64) {
	//fmt.Printf("parcelValues: %+v, cropCov: %+v, stWeight: %+v\n", parcelValues, cropCov, stWeight)
	for month, v := range parcelValues {
		values[month] += v
		for crop, cropNir := range cropsNIR {
			values[month] += cropNir[month] * cropCov[crop] * stWeight * area / 12
		}
	}

	return values
}
