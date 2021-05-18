package parcelpump

import "github.com/heath140/wwum2020/database"

func (p *Parcel) setAppEfficiency(efficiencies []database.Efficiency, year int) {
	f, s := filterEff(efficiencies, year)

	if p.IrrType.String == "Flood" {
		p.AppEff = f
	} else {
		p.AppEff = s
	}
}

func filterEff(efficiencies []database.Efficiency, year int) (f float64, s float64) {
	for _, v := range efficiencies {
		if v.Yr == year {
			f = v.AeFlood
			s = v.AeSprinkler
		}
	}

	return f, s
}
