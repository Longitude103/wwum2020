package wells

import (
	"errors"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/parcels"
)

func WriteWELFile(v database.Setup, parcels []parcels.Parcel) error {
	// get a list of the wells and associated parcels
	wellParcels, err := database.GetWellParcels(v)
	if err != nil {
		return err
	}

	wellNode, err := database.GetWellNode(v)
	if err != nil {
		return err
	}

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
		} else if count == 1 {
			// single well, just the pumping for the parcel
			// get well node
			node, err := getNode(wellNode, wls[0], parcels[p].Nrd)
			if err != nil {
				return err
			}

			// write to the db
			_ = node

		} else {
			// multiple wells
		}

	}

	return nil
}

func filterWells(wlPar []database.WellParcel, parcel int, nrd string, yr int) (wells []int, count int, err error) {
	for i := 0; i < len(wlPar); i++ {
		if wlPar[i].Yr == yr && wlPar[i].Nrd == nrd && wlPar[i].ParcelId == parcel {
			wells = append(wells, wlPar[i].WellId)
		}
	}

	count = len(wells)

	return
}

func getNode(wellNodes []database.WellNode, well int, nrd string) (int, error) {
	for i := 0; i < len(wellNodes); i++ {
		if wellNodes[i].WellId == well && wellNodes[i].Nrd == nrd {
			return wellNodes[i].Node, nil
		}
	}

	return 0, errors.New("no well found")
}
