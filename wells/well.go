package wells

import (
	"errors"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/parcels"
	"time"
)

// TODO: Create Tests for this file.

// WriteWELResults is a function that gets the pumping amounts for the parcel and assigns them to a well or wells that
// supply that parcel.
func WriteWELResults(v database.Setup, parcels []parcels.Parcel) error {
	// get a list of the wells and associated parcels
	wellParcels, err := database.GetWellParcels(v)
	if err != nil {
		return err
	}

	wellNode, err := database.GetWellNode(v)
	if err != nil {
		return err
	}

	var welResult []database.WelResult
	for p := 0; p < len(parcels); p++ {
		// no GW, skip
		if !parcels[p].Gw.Bool {
			continue
		}

		// find wells
		wls, count, err := filterWells(wellParcels, parcels[p].ParcelNo, parcels[p].Nrd, parcels[p].Yr)
		if err != nil {
			return err
		}

		if count == 0 {
			continue
		} else {
			// can be one or multiple wells
			for _, w := range wls {
				if err = addToResults(wellNode, welResult, w, parcels[p], count); err != nil {
					return err
				}
			}
		}
	}

	// sum the same wells together
	groupedResult, err := sumWells(welResult)

	welDB, err := database.ResultsWelDB(v.SlDb)
	if err != nil {
		return err
	}

	// save groupedResult to DB
	for i := 0; i < len(groupedResult); i++ {
		err = welDB.Add(groupedResult[i])
		if err != nil {
			return err
		}
	}

	err = welDB.Flush()
	if err != nil {
		return err
	}

	return nil
}

// filterWells is a function to filter out the wells by an nrd and parcel number and returns a slice of wells that
// supply that parcel
func filterWells(wlPar []database.WellParcel, parcel int, nrd string, yr int) (wells []int, count int, err error) {
	for i := 0; i < len(wlPar); i++ {
		if wlPar[i].Yr == yr && wlPar[i].Nrd == nrd && wlPar[i].ParcelId == parcel {
			wells = append(wells, wlPar[i].WellId)
		}
	}

	count = len(wells)

	return
}

// getNode is a function that gets the node that a well is located in based on well id and nrd string and returns the
// node value in an int.
func getNode(wellNodes []database.WellNode, well int, nrd string) (int, error) {
	for i := 0; i < len(wellNodes); i++ {
		if wellNodes[i].WellId == well && wellNodes[i].Nrd == nrd {
			return wellNodes[i].Node, nil
		}
	}

	return 0, errors.New("no well found")
}

// addToResults is the function that creates another result from the parcel and adds to the result slice.
func addToResults(wellNode []database.WellNode, r []database.WelResult, well int, p parcels.Parcel, count int) error {
	node, err := getNode(wellNode, well, p.Nrd)
	if err != nil {
		return err
	}

	ft, err := p.SetWelFileType()
	if err != nil {
		return err
	}

	for i, d := range p.Pump {
		if d > 0 {
			r = append(r, database.WelResult{Wellid: well, Node: node, Dt: time.Date(p.Yr,
				time.Month(i+1), 1, 0, 0, 0, 0, time.UTC), FileType: ft, Result: d / float64(count)})
		}
	}

	return nil
}

// sumWells is a function to add the same well in the same month together for one pumping value.
func sumWells([]database.WelResult) (result []database.WelResult, err error) {
	// run through and add wells together that are the same well, month and nrd, output a new slice.
	// TODO: Create this function or might want to do this in addToResults as part of that function
	// similar to dryland.go[60:98] file. Look there.

	return result, nil
}
