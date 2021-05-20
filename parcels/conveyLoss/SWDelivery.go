package conveyLoss

import (
	"github.com/jmoiron/sqlx"
)

// GetSurfaceWaterDelivery function returns a slice of Diversion that is a monthly amount of surface water delivered to
// an acre of land. The units of the Diversion are in acre-feet per acre for use in subsequent processes.
func GetSurfaceWaterDelivery(pgDb *sqlx.DB, sYear int, eYear int) []Diversion {
	diversions := getDiversions(pgDb, sYear, eYear, false)
	canals := getCanals(pgDb, sYear, eYear)

	for i := 0; i < len(diversions); i++ {
		c := filterCnl(canals, (&diversions[i]).CanalId, (&diversions[i]).DivDate.Time.Year())
		if c.Area.Valid {
			(&diversions[i]).applyEffAcres(c.Eff, c.Area.Float64)
		} else {
			(&diversions[i]).DivAmount.Float64 = 0
		}
	}

	return diversions
}

// filterCnl filters the list of canals to a specific one.
func filterCnl(canals []Canal, canal int, yr int) (c Canal) {
	for _, v := range canals {
		if v.Id == canal && v.Yr == yr {
			c = v
		}
	}

	return c
}

func FilterSWDeliveryByYear(divs []Diversion, y int) (diversions []Diversion) {
	for _, v := range divs {
		if v.DivDate.Time.Year() == y {
			diversions = append(diversions, v)
		}
	}

	return diversions
}
