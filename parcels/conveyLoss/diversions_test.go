package conveyLoss

import (
	"testing"
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
