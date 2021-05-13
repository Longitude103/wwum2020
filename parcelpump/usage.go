package parcelpump

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

// distributeUsage function takes in a map of parcel id with monthly values and a usage total and distributes the usage
// by those nir values. Handles multiple parcels and returns the usage in a map of parcel ids with 12 monthly values.
func distributeUsage(p map[int][12]float64, u float64) map[int][12]float64 {
	var totalNIR float64
	var totalMonthlyNIR [12]float64
	for k := range p {
		for i := 0; i < 12; i++ {
			totalMonthlyNIR[i] += p[k][i]
			totalNIR += p[k][i]
		}
	}

	distUsage := map[int][12]float64{}
	for k := range p {
		monthDistribution := [12]float64{}
		for i := 0; i < 12; i++ {
			if totalMonthlyNIR[i] > 0 { // protect from division by zero
				monthPercent := totalMonthlyNIR[i] / totalNIR
				monthDistribution[i] = (p[k][i] / totalMonthlyNIR[i]) * (monthPercent * u)
			}
		}
		distUsage[k] = monthDistribution
	}

	return distUsage
}
