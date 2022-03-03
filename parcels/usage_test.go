package parcels

import (
	"testing"
)

func Test_distributeUsage(t *testing.T) {
	tNir := 2.7
	tMonthlyNir := [12]float64{0, 0, 0, 0, 0.2, 0.4, 0.8, 0.8, 0.5, 0, 0, 0}
	tUsage := 100.0

	p1.distributeUsage(tNir, tMonthlyNir, tUsage)

	tPump := 0.0
	for _, v := range p1.Pump {
		tPump += v
	}

	if tPump < 99.9 {
		t.Errorf("Total should be close to 100, but recieved %f", tPump)
	}

}

func Test_getUsage(t *testing.T) {
	v := dbConnection()
	v.SetYears(1997, 2020)

	usage := getUsage(v)

	records2018 := len(usage[2018])

	if records2018 != 2935 {
		t.Errorf("should have returned 2935 records but got %d records", records2018)
	}
}
