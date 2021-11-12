package conveyLoss

import (
	"github.com/Longitude103/wwum2020/database"
)

// GetSurfaceWaterDelivery function returns a slice of Diversion that is a monthly amount of surface water delivered to
// an acre of land. The units of the Diversion are in acre-feet per acre for use in subsequent processes.
func GetSurfaceWaterDelivery(v *database.Setup) ([]Diversion, error) {
	var db *database.SWDelDB
	var err error
	if !v.AppDebug {
		db, err = database.SWDeliveryDB(v.SlDb)
		if err != nil {
			return []Diversion{}, err
		}
	}

	diversions, err := getDiversions(v)
	if err != nil {
		v.Logger.Errorf("Error in getDiversions: %s", err)
		return nil, err
	}

	canals, err := getCanals(v)
	if err != nil {
		v.Logger.Errorf("Error in getCanals: %s", err)
		return nil, err
	}

	for i := 0; i < len(diversions); i++ {
		c := filterCnl(canals, (&diversions[i]).CanalId, (&diversions[i]).DivDate.Time.Year())
		if c.Area.Valid {
			// apply efficiency and convert to AF
			(&diversions[i]).applyEffAcres(c.Eff, c.Area.Float64)
		} else {
			(&diversions[i]).DivAmount.Float64 = 0
		}
	}

	if !v.AppDebug {
		for i := 0; i < len(diversions); i++ {
			if err := db.Add(database.SWDelResult{CanalId: diversions[i].CanalId, Dt: diversions[i].DivDate.Time,
				DelAmount: diversions[i].DivAmount.Float64}); err != nil {
				return []Diversion{}, err
			}
		}

		_ = db.Flush()
		_ = db.Close()
	}

	return diversions, nil
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
