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
	for yr := v.SYear; yr < v.EYear+1; yr++ {
		if yr < 1997 {
			for _, well := range wells {
				// TODO: Only run with wells that start pre 1997
				constMIWell(well, yr, welDB)
			}
		}

	}

	// add pumping from the actual pumping table to results db

	return nil
}

func constMIWell(well database.MIWell, yr int, db *database.WelDB) []database.WelResult {
	var wrList []database.WelResult
	annVolume := -1.0 * float64(well.Rate) * float64(DaysInYear(yr)) / 43560
	for i := 0; i < 12; i++ {
		monthVol := annVolume / float64(DaysInMonth(time.Date(yr, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)))
		wl := database.WelResult{Wellid: well.WellId, Node: well.Node, FileType: MIFileType(well),
			Dt: time.Date(yr, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC), Result: monthVol}
		wrList = append(wrList, wl)
	}

	return wrList
}

func MIFileType(well database.MIWell) int {
	if well.MuniWell {
		return 210
	} else if well.IndustWell {
		return 211
	} else {
		return 212
	}
}

func EndOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	beginMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)

	return beginMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// EndOfYear end of year
func EndOfYear(t time.Time) time.Time {
	y, _, _ := t.Date()
	beginYear := time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC)

	return beginYear.AddDate(1, 0, 0).Add(-time.Nanosecond)
}

func DaysInMonth(t time.Time) int {
	_, _, d := EndOfMonth(t).Date()

	return d
}

func DaysInYear(yr int) int {
	ey := EndOfYear(time.Date(yr, 1, 1, 0, 0, 0, 0, time.UTC))
	return ey.YearDay()
}
