package parcels

import (
	"github.com/Longitude103/wwum2020/database"
	"testing"
)

var (
	eff1         = database.Efficiency{Yr: 2014, AeFlood: 0.65, AeSprinkler: 0.85}
	eff2         = database.Efficiency{Yr: 2013, AeFlood: 0.60, AeSprinkler: 0.80}
	efficiencies = []database.Efficiency{eff1, eff2}
)

func Test_filterEff(t *testing.T) {
	f, s := filterEff(efficiencies, 2014)

	if f != 0.65 || s != 0.85 {
		t.Errorf("filter Efficiencies not correct for 2014: got %f, expected 0.65; got %f, expected 0.85", f, s)
	}
}

func TestParcel_setAppEfficiency(t *testing.T) {
	p1.setAppEfficiency(efficiencies, 2013)
	p3.setAppEfficiency(efficiencies, 2013)

	if p1.AppEff != 0.8 {
		t.Errorf("AppEff not set correctly: got %f, expected 0.8", p1.AppEff)
	}

	if p3.AppEff != 0.6 {
		t.Errorf("AppEff not set correctly: got %f, expected 0.8", p3.AppEff)
	}
}
