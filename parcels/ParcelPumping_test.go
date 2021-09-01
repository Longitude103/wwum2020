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
	// TODO: Finish and clean up this test.
	for i := 0; i < 12; i++ {
		testParcelSlice[0].Pump[i] = 0
	}

	err := distUsage(testUsageSlice, &testParcelSlice)
	if err != nil {
		t.Error("Function returned an error:", err)
	}

	totalNIR := 0.0
	for _, parcel := range testParcelSlice {
		for _, f := range parcel.Nir {
			totalNIR += f
		}
	}

	fmt.Println("TotalNIR:", totalNIR)

	total := 0.0
	parcelNIR := 0.0
	for i, f := range testParcelSlice[0].Pump {
		total += f
		parcelNIR += testParcelSlice[0].Nir[i]
	}

	fmt.Println("Monthly Pumping:", testParcelSlice[0].Pump)
	fmt.Println("Total Parcel 1234 NIR:", parcelNIR)
	fmt.Printf("Parcel 1234 is %f percent of total NIR\n", parcelNIR/totalNIR)
	fmt.Println("TotalPumping:", total)
}
