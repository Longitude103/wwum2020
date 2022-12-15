package conveyLoss

import (
	"database/sql"
	"testing"
	"time"
)

func Test_getDiversions(t *testing.T) {
	v := dbConnection()

	v.SYear = 1953
	v.EYear = 1953

	divs, _, err := getDiversions(v)
	if err != nil {
		t.Error("Get Diversions errored")
	}

	for _, div := range divs {
		if div.CanalId == 26 {
			v.Logger.Debugf("Laramie Div: %+v", div)
		}

		if div.CanalId == 272 {
			v.Logger.Debugf("Laramie WY Div: %+v", div)
		}

		if div.CanalId == 13 {
			v.Logger.Debugf("Mitchell Div: %+v", div)
		}

		if div.CanalId == 32 {
			v.Logger.Debugf("Gering Div: %+v", div)
		}

		if div.CanalId == 15 {
			v.Logger.Debugf("Enterprise Div: %+v", div)
		}

		if div.CanalId == 29 {
			v.Logger.Debugf("Minatare Div: %+v", div)
		}
	}

	// if divs[1].DivAmount.Float64 != 12797.0 {
	// 	t.Errorf("Wrong amount being queried: Should be 12797, got %f", divs[1].DivAmount.Float64)
	// }

}

func Test_getDiversionsSS(t *testing.T) {
	v := dbConnection()
	_ = v.SetYears(1953, 2020)
	v.SteadyState = true

	divs, _, err := getDiversions(v)
	if err != nil {
		t.Error("Get Diversions errored")
	}

	for _, d := range divs {
		if d.DivDate.Time.Year() > 1952 || d.DivDate.Time.Year() < 1895 {
			t.Errorf("Div Dates are wrong: %+v", d)
		}

		if d.CanalId == 0 {
			t.Errorf("Bad Canal Id: %+v", d)
		}
	}
}

func Test_getEFDiversions(t *testing.T) {
	v := dbConnection()
	if err := v.SetYears(2016, 2016); err != nil {
		t.Error("Failed to set years")
	}

	_, div, err := getDiversions(v)
	if err != nil {
		t.Error("Error in getDiversions: ", err)
	}

	want := 15
	got := len(div)
	if want != got {
		t.Errorf("getDiversions should have returned %d diversion records but got %d records", want, got)
	}
}

func TestEfPeriod_GetYearAndMonths(t *testing.T) {
	ef := efPeriod{
		CanalId:    1,
		StartDate:  sql.NullTime{Valid: true, Time: time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)},
		EndDate:    sql.NullTime{Valid: true, Time: time.Date(2016, 6, 6, 0, 0, 0, 0, time.UTC)},
		LossPercet: sql.NullFloat64{Valid: true, Float64: 0.56},
	}

	y, months := ef.GetYearAndMonths()

	if y != 2016 {
		t.Errorf("Was expecting year 2016, got year %d", y)
	}

	if len(months) != 3 {
		t.Errorf("was expecting a return of 3 months, but got %d months", len(months))
	}

	for _, m := range months {
		if m < 4 && m > 6 {
			t.Errorf("months should be 4, 5, or 6, but got %d", m)
		}
	}
}
