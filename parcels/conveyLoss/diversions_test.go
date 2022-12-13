package conveyLoss

import (
	"fmt"
	"testing"
	"time"
)

func Test_getDiversions(t *testing.T) {
	v := dbConnection()

	v.SYear = 1953
	v.EYear = 1953

	divs, err := getDiversions(v)
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
	v.SetYears(1953, 2020)
	v.SteadyState = true

	divs, err := getDiversions(v)
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

	withOutEFDivs, err := getDiversions(v)
	if err != nil {
		t.Error("Get Diversions errored")
	}

	v.ExcessFlow = true
	efDivs, err := getDiversions(v)
	if err != nil {
		t.Error("Get Diversions errored")
	}

	canalFilter := 29

	var filteredWO, filteredWEF []Diversion
	for _, div := range withOutEFDivs {
		if div.CanalId == canalFilter {
			filteredWO = append(filteredWO, div)
		}
	}

	for _, div := range efDivs {
		if div.CanalId == canalFilter {
			filteredWEF = append(filteredWEF, div)
		}
	}

	var WOEFMay, WEFMay Diversion
	for _, div := range filteredWO {
		if div.DivDate.Time.Month() == time.May {
			WOEFMay = div
		}
	}

	for _, div := range filteredWEF {
		if div.DivDate.Time.Month() == time.May {
			WEFMay = div
		}
	}

	if WOEFMay.DivAmount.Float64 >= WEFMay.DivAmount.Float64 {
		t.Error("May 2016 Canal 29 Excess Flow is not being removed")
		fmt.Println("Without Excess Flow: ", filteredWO)
		fmt.Println("With Excess Flow: ", filteredWEF)
	}
}
