package conveyLoss

import (
	"fmt"
	"testing"
)

func Test_getDiversions(t *testing.T) {
	v := dbConnection()

	divs, err := getDiversions(v)
	if err != nil {
		t.Error("Get Diversions errored")
	}

	for _, div := range divs {
		if div.CanalId == 26 {
			fmt.Printf("Laramie Div: %+v\n", div)
		}

		if div.CanalId == 13 {
			fmt.Printf("Mitchell Div: %+v\n", div)
		}

		if div.CanalId == 32 {
			fmt.Printf("Gering Div: %+v\n", div)
		}

		if div.CanalId == 15 {
			fmt.Printf("Enterprise Div: %+v\n", div)
		}
	}

	if divs[1].DivAmount.Float64 != 12797.0 {
		t.Errorf("Wrong amount being queried: Should be 12797, got %f", divs[1].DivAmount.Float64)
	}

}
