package Utils

import (
	"math"
	"testing"
)

func TestRoundTo(t *testing.T) {
	pi := math.Pi

	if RoundTo(pi, 3) != 3.142 {
		t.Errorf("round function not working: got %g, expected 3.142", RoundTo(pi, 3))
	}

}
