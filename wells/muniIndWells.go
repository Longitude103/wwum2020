package wells

import (
	"github.com/Longitude103/wwum2020/database"
	"time"
)

// MunicipalIndWells is a function that adds the municipal and industrial wells from postgresql to the results database
// and uses either assumed pumping rates or actual pumping numbers.
func MunicipalIndWells(v database.Setup) error {
	// go get the wells data
	wells, err := database.GetMIWells(v)
	if err != nil {
		return err
	}

	// process the data for the monthly amounts for average data
	welDB, err := database.ResultsWelDB(v.SlDb)
	if err != nil {
		return err
	}

	// start97 == false then use the "rate" to create the monthly pumping
	var wlResult []database.WelResult

	for yr := v.SYear; yr < v.EYear+1; yr++ {
		if yr < 1997 {
			for _, well := range wells {
				if well.Start97 == false {
					wlResult = append(wlResult, constMIWell(well, TimeExt{y: yr})...)
				}
			}
		}

		if yr >= 1997 {
			for _, well := range wells {
				if well.Stop97 == false && well.Start97 == false {
					wlResult = append(wlResult, constMIWell(well, TimeExt{y: yr})...)
				}

				if well.Start97 {
					wlResult = append(wlResult, pumpMIWell(well, TimeExt{y: yr})...)
				}
			}
		}

		for i := 0; i < len(wlResult); i++ {
			if err := welDB.Add(wlResult[i]); err != nil {
				return err
			}
		}

	}

	return nil
}

func constMIWell(well database.MIWell, yr TimeExt) []database.WelResult {
	var wrList []database.WelResult
	annVolume := -1.0 * float64(well.Rate) * float64(yr.DaysInYear()) / 43560
	for i := 0; i < 12; i++ {
		dInMon := TimeExt{t: time.Date(yr.y, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)}
		monthVol := annVolume / float64(dInMon.DaysInMonth())
		wl := database.WelResult{Wellid: well.WellId, Node: well.Node, FileType: well.MIFileType(),
			Dt: time.Date(yr.y, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC), Result: monthVol}
		wrList = append(wrList, wl)
	}

	return wrList
}

func pumpMIWell(well database.MIWell, yr TimeExt) []database.WelResult {
	var wrList []database.WelResult
	for _, p := range well.Pumping {
		if p.PumpDate.Year() == yr.y {
			wl := database.WelResult{Wellid: well.WellId, Node: well.Node, FileType: well.MIFileType(),
				Dt: p.PumpDate, Result: p.Pump}
			wrList = append(wrList, wl)
		}
	}

	return wrList
}

type TimeExt struct {
	t time.Time
	y int
}

func (tm TimeExt) EndOfMonth() time.Time {
	y, m, _ := tm.t.Date()
	beginMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)

	return beginMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// EndOfYear end of year
func (tm TimeExt) EndOfYear() time.Time {
	y, _, _ := tm.t.Date()
	beginYear := time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC)

	return beginYear.AddDate(1, 0, 0).Add(-time.Nanosecond)
}

func (tm TimeExt) DaysInMonth() int {
	_, _, d := tm.EndOfMonth().Date()

	return d
}

func (tm TimeExt) DaysInYear() int {
	t := TimeExt{t: time.Date(tm.y, 1, 1, 0, 0, 0, 0, time.UTC)}
	ey := t.EndOfYear()
	return ey.YearDay()
}
