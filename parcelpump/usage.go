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
