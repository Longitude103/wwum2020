package rchFiles

import (
	"errors"
	"fmt"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/parcels"
	"github.com/pterm/pterm"
	"time"
)

// Dryland gets a slice of dryland parcels and writes out the values to the results' database.
func Dryland(v database.Setup, dryParcels []parcels.Parcel) error {
	p, _ := pterm.DefaultProgressbar.WithTotal(v.EYear - v.SYear + 1).WithTitle("Dryland Recharge Results").WithRemoveWhenDone(true).Start()

	for y := v.SYear; y < v.EYear+1; y++ {
		p.UpdateTitle(fmt.Sprintf("Getting %d cells and filtering them", y))
		dryCells := database.GetDryCells(v, y) // will need to iterate through years
		annParcels, err := parcelFilterByYear(dryParcels, y)
		if err != nil {
			v.Logger.Errorf("Didn't return parcels from filter by year for year: %d", y)
			return err
		}

		p.UpdateTitle(fmt.Sprintf("Writing %d values to DB", y))
		var preResults []database.RchResult
		for i := 0; i < len(dryCells); i++ {
			parcelArea, rf, err := parcelValues(annParcels, int(dryCells[i].PId), dryCells[i].Nrd)
			if err != nil {
				v.Logger.Errorf("Dryland RCH parcelValues error, year: %d, parcel_id: %d, nrd: %s", y,
					int(dryCells[i].PId), dryCells[i].Nrd)
				return err
			}

			for m := 0; m < 12; m++ {
				if rf[m] > 0 {
					preResults = append(preResults,
						database.RchResult{Node: dryCells[i].Node, Size: dryCells[i].CellArea,
							Dt:       time.Date(y, time.Month(m+1), 1, 0, 0, 0, 0, time.UTC),
							FileType: 101, Result: rf[m] * dryCells[i].DryArea / parcelArea})
				}
			}
		}

		//results := groupResults(preResults)

		for i := 0; i < len(preResults); i++ {
			err := v.RchDb.Add(preResults[i])
			if err != nil {
				return err
			}
		}
		p.Increment()
	}

	return nil
}

// parcelValues is a function that returns the area of the parcel and the monthly return flow (rf) values for processing.
func parcelValues(p []parcels.Parcel, id int, nrd string) (area float64, rf [12]float64, err error) {
	for i := 0; i < len(p); i++ {
		if p[i].ParcelNo == id && p[i].Nrd == nrd {
			for m := 0; m < 12; m++ {
				rf[m] += p[i].Ro[m] + p[i].Dp[m]
			}
			return p[i].Area, rf, nil
		}
	}

	return 0, rf, errors.New("parcel not found")
}

// findResult takes a slice of RchResult and returns the one that matches the node number and date
func findResult(r []database.RchResult, node int, dt time.Time) (found bool, location int) {
	for i := 0; i < len(r); i++ {
		if r[i].Node == node && r[i].Dt == dt {
			return true, i
		}
	}

	return false, 0
}

// inGrouped is a function that looks to see if the node is present in the slice
func inGrouped(g []database.RchResult, node int) bool {
	for _, i := range g {
		if i.Node == node {
			return true
		}
	}

	return false
}

// groupResults is a function to gorup the results together to make a smaller results set so that if there are more than
// one node results, they are added together and made into one value.
func groupResults(r []database.RchResult) (results []database.RchResult) {
	for i := 0; i < len(r); i++ {
		if inGroup := inGrouped(results, r[i].Node); inGroup {
			// add to the result this value
			found, resultLocal := findResult(results, r[i].Node, r[i].Dt)
			if found {
				results[resultLocal].Result += r[i].Result
			}
		} else {
			results = append(results, database.RchResult{Node: r[i].Node, Size: r[i].Size, Dt: r[i].Dt,
				FileType: r[i].FileType, Result: r[i].Result})
		}
	}

	return
}
