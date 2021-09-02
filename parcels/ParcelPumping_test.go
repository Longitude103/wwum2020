package parcels

import (
	"testing"
)

var u1 = Usage{Yr: 2014, Nrd: "np", CertNum: "3456", UseAF: 100.0}
var u2 = Usage{Yr: 2014, Nrd: "np", CertNum: "3459", UseAF: 240.0}
var u3 = Usage{Yr: 2014, Nrd: "np", CertNum: "3457", UseAF: 100.0}
var u4 = Usage{Yr: 2014, Nrd: "np", CertNum: "3458", UseAF: 200.0}

var testUsageSlice = []Usage{u1, u2, u3, u4}

func Test_distUsage(t *testing.T) {
	for i := 0; i < 12; i++ {
		testParcelSlice[0].Pump[i] = 0
	}

	err := distUsage(testUsageSlice, &testParcelSlice)
	if err != nil {
		t.Error("Function returned an error:", err)
	}

	total := 0.0
	for _, f := range testParcelSlice[0].Pump {
		total += f
	}

	if total < 9.89 || total > 9.9 {
		t.Errorf("Total pumping should have been 9.897 but got %f", total)
	}
}
