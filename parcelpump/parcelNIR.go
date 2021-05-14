package parcelpump

import (
	"fmt"
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
	"sort"
	"time"

	"github.com/heath140/gisUtils"
	"github.com/jmoiron/sqlx"
)

// parcelNIR is a method that adds the NIR, RO, and DP for each parcel from the CSResults and weather station data.
// It produces an intermediate results table of NIR in local sqlite for review and adds three maps to the parcel struct
func (p *Parcel) parcelNIR(slDB *sqlx.DB, Year int, wStations []database.WeatherStation, csResults map[string][]fileio.StationResults) {
	var parcelNIR, parcelRo, parcelDp [12]float64
	if p.Nir == nil {
		p.Nir = map[int][12]float64{}
	}

	if p.Ro == nil {
		p.Ro = map[int][12]float64{}
	}

	if p.Dp == nil {
		p.Dp = map[int][12]float64{}
	}

	dist := distances(*p, wStations)
	for _, st := range dist {
		var annData []fileio.StationResults
		for _, data := range csResults[st.Station] {
			if data.Yr == Year && data.Soil == p.SoilCode && data.Irrigation == 3 {
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
		parcelNIR = pValues(parcelNIR, cropsNir, cropCov, st.Weight)
		parcelRo = pValues(parcelRo, cropsRo, cropCov, st.Weight)
		parcelDp = pValues(parcelDp, cropsDp, cropCov, st.Weight)
	}
	//fmt.Printf("Weighted Parcel ID: %d, NIR is: %v\n", parcel.ParcelNo, parcelNIR)

	// save parcelNIR to sqlite
	saveSqlite(slDB, p.ParcelNo, p.Nrd, parcelNIR, Year)
	p.Nir[Year] = parcelNIR
	p.Ro[Year] = parcelRo
	p.Dp[Year] = parcelDp
}

// distances is a function that that returns the top three weather stations from the list with the appropriate weighting
// factor. Used to make CSResults Distribution.
func distances(parcel Parcel, wStations []database.WeatherStation) []database.StDistances {
	var dist []database.StDistances
	var lengths []float64
	for _, v := range wStations {
		var stDistance database.StDistances
		d := gisUtils.Distance(parcel.PointY, parcel.PointX, v.Cor.Coordinates[1], v.Cor.Coordinates[0])
		lengths = append(lengths, d)
		stDistance.Distance = d
		stDistance.Station = v.Code
		dist = append(dist, stDistance)
	}

	sort.Slice(dist, func(i, j int) bool {
		return dist[i].Distance < dist[j].Distance
	})

	sort.Float64s(lengths)

	idw, err := gisUtils.InverseDW(lengths[:3])
	if err != nil {
		fmt.Println("Error", err)
	}

	for i, v := range idw {
		dist[i].Weight = v
	}

	return dist[:3]
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
func pValues(parcelValues [12]float64, cropsNIR [4][12]float64, cropCov [4]float64, stWeight float64) (values [12]float64) {
	for month, v := range parcelValues {
		values[month] += v
		for crop, cropNir := range cropsNIR {
			values[month] += cropNir[month] * cropCov[crop] * stWeight
		}
	}

	return values
}

// saveSqlite function saves the data for the parcel into local sqlite so that additional error checking can be preformed
// without loosing the data.
func saveSqlite(slDB *sqlx.DB, parcelID int, nrd string, pNIR [12]float64, yr int) {
	tx := slDB.MustBegin()

	for i, v := range pNIR {
		dt := time.Date(yr, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
		tx.MustExec("INSERT INTO parcelNIR (parcelID, nrd, dt, nir) VALUES ($1, $2, $3, $4)", parcelID, nrd, dt.Format(time.RFC3339), v)
	}

	err := tx.Commit()
	if err != nil {
		fmt.Println("Error in SQLite Commit", err)
	}
}
