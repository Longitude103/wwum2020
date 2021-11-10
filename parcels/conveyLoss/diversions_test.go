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

	for i := 0; i < 12; i++ {
		fmt.Printf("Diversion: %+v\n", divs[i])
	}

	if divs[1].DivAmount.Float64 != 12797.0 {
		t.Errorf("Wrong amount being queried: Should be 12797, got %f", divs[1].DivAmount.Float64)
	}

}
