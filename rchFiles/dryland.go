package rchFiles

import (
	"errors"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/parcels"
	"time"
)

// Dryland gets a slice of dryland parcels and writes out the values to the results' database.
func Dryland(v database.Setup, dryParcels []parcels.Parcel) error {
	for y := v.SYear; y < v.EYear+1; y++ {

		dryCells := database.GetDryCells(v, y) // will need to iterate through years

		var preResults []database.RchResult
		for i := 0; i < len(dryCells); i++ {
			parcelArea, rf, err := parcelValues(dryParcels, int(dryCells[i].PId), dryCells[i].Nrd)
			if err != nil {
				return err
			}

			for m := 0; m < 12; m++ {
				if rf[m] > 0 {
					preResults = append(preResults,
						database.RchResult{Node: dryCells[i].Node,
							Dt:       time.Date(y, time.Month(m+1), 1, 0, 0, 0, 0, time.UTC),
							FileType: 101, Result: rf[m] * dryCells[i].DryArea / parcelArea})
				}
			}
		}

		results := groupResults(preResults)

		for i := 0; i < len(results); i++ {
			err := v.RchDb.Add(results[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// parcelValues is a function that returns the area of the parcel and the monthly return flow (rf) values for processing.
func parcelValues(p []parcels.Parcel, id int, nrd string) (area float64, rf [12]float64, err error) {
	for i := 0; i < len(p); i++ {
		if p[i].ParcelNo == id && p[i].Nrd == nrd {
			for m := 0; m < 12; m++ {
				rf[m] += p[i].Ro[m] + p[i].Dp[m]
				return p[i].Area, rf, nil
			}
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
			results = append(results, database.RchResult{Node: r[i].Node, Dt: r[i].Dt, FileType: r[i].FileType, Result: r[i].Result})
		}
	}

	return
}
