package parcels

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Usage struct {
	CertNum string  `db:"cert_num"`
	UseAF   float64 `db:"usage_af"`
	Yr      int     `db:"yr"`
	Nrd     string  `db:"nrd"`
}

// getUsage function returns a slice of usage struct for all cert usage in both nrds. It includes all years.
func getUsage(pgDB *sqlx.DB) []Usage {
	query := `SELECT cert_num::varchar, usage_af, yr, 'np' nrd from np.np_usage UNION SELECT cert_num, usage_af, yr, 'sp' nrd from sp.sp_usage;`

	var use []Usage
	err := pgDB.Select(&use, query)
	if err != nil {
		fmt.Println("Error in Usage:", err)
	}

	return use
}

// filterUsage is a function that filters the total usage to a year and returns a slice of Usage
func filterUsage(u []Usage, yr int) (filteredUsage []Usage) {
	for _, v := range u {
		if v.Yr == yr {
			filteredUsage = append(filteredUsage, v)
		}
	}

	return filteredUsage
}

// distributeUsage method receives the total NIR, monthly NIR values and usage for the cert and distributes that pumping
// to the parcel. It also sets the parcel metered property to true.
func (p *Parcel) distributeUsage(totalNIR float64, totalMonthlyNIR [12]float64, u float64) {
	for m := 0; m < 12; m++ {
		if totalMonthlyNIR[m] > 0 { // protect from division by zero
			monthPercent := totalMonthlyNIR[m] / totalNIR
			p.Pump[m] = p.Nir[m] / totalMonthlyNIR[m] * (monthPercent * u)
		}
	}
}
