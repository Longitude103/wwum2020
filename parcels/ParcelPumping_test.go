package parcels

import (
	"fmt"
	"testing"
)

var u1 = Usage{Yr: 2014, Nrd: "np", CertNum: "3456", UseAF: 100.0}
var u2 = Usage{Yr: 2014, Nrd: "np", CertNum: "3459", UseAF: 240.0}
var u3 = Usage{Yr: 2014, Nrd: "np", CertNum: "3457", UseAF: 100.0}
var u4 = Usage{Yr: 2014, Nrd: "np", CertNum: "3458", UseAF: 200.0}

var testUsageSlice = []Usage{u1, u2, u3, u4}

func Test_distUsage(t *testing.T) {
	// TODO: Test this function.
	err := distUsage(testUsageSlice, &testParcelSlice)
	if err != nil {
		t.Error("Function returned an error: ", err)
	}

	fmt.Println(testParcelSlice[0].Pump)

}
