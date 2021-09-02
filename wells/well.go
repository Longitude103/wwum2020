package wells

import (
	"errors"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/parcels"
)

// WriteWELResults is a function that gets the pumping amounts for the parcel and assigns them to a well or wells that
// supply that parcel.
func WriteWELResults(v database.Setup, parcels *[]parcels.Parcel) error {
	v.Logger.Info("Starting WriteWELResults...")
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
	for p := 0; p < len(*parcels); p++ {
		// no GW, skip
		if !(*parcels)[p].Gw.Bool {
			continue
		}

		wls, count, err := filterWells(wellParcels, (*parcels)[p].ParcelNo, (*parcels)[p].Nrd, (*parcels)[p].Yr)
		if err != nil {
			return err
		}

		if count == 0 {
			continue
		} else {
			// can be one or multiple wells
			for _, w := range wls {
				if welResult, err = addToResults(wellNode, welResult, w, (*parcels)[p], count); err != nil {
					return err
				}
			}
		}
	}

	welDB, err := database.ResultsWelDB(v.SlDb)
	if err != nil {
		return err
	}

	// save groupedResult to DB
	for i := 0; i < len(welResult); i++ {
		err = welDB.Add(welResult[i])
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
func addToResults(wellNode []database.WellNode, r []database.WelResult, well int, p parcels.Parcel, count int) ([]database.WelResult, error) {
	node, err := getNode(wellNode, well, p.Nrd)
	if err != nil {
		return r, err
	}

	// if the well is there, then just add the value
	if found, local := findResult(r, well, p.Yr); found {
		// use local and call add
		r[local].AddPumping(p.Pump, float64(count))
	} else {
		// Otherwise, create a new well and add it to the slice
		ft, err := p.SetWelFileType()
		if err != nil {
			return r, err
		}

		var result [12]float64
		for i, d := range p.Pump {
			result[i] = d / float64(count)
		}

		r = append(r, database.WelResult{Wellid: well, Node: node, Yr: p.Yr, FileType: ft, Result: result})
	}

	return r, nil
}

// findResult is a function to find if there is a slice of database.WelResult that has a well and year and returns a
// bool if it is found and a location in the slice that it is located
func findResult(r []database.WelResult, well int, yr int) (found bool, location int) {
	for i := 0; i < len(r); i++ {
		if r[i].Wellid == well && r[i].Yr == yr {
			return true, i
		}
	}

	return false, 0
}
