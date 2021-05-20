package parcels

import "github.com/heath140/wwum2020/database"

// setAppEfficiency is a Parcel method that uses the parcel IrrType and year to determine the application efficiency of the
// parcel and sets that in the Parcel struct. It takes a slice of Efficiency and year as they very through the study period.
func (p *Parcel) setAppEfficiency(efficiencies []database.Efficiency, year int) {
	f, s := filterEff(efficiencies, year)

	if p.IrrType.String == "Flood" {
		p.AppEff = f
	} else {
		p.AppEff = s
	}
}

// filterEff returns the efficiency of flood and sprinkler for a given year.
func filterEff(efficiencies []database.Efficiency, year int) (f float64, s float64) {
	for _, v := range efficiencies {
		if v.Yr == year {
			f = v.AeFlood
			s = v.AeSprinkler
		}
	}

	return f, s
}
