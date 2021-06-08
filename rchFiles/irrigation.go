package rchFiles

import (
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/parcels"
)

func IrrigationRCH(v database.Setup, AllParcels []parcels.Parcel) error {

	for y := v.SYear; y < v.EYear+1; y++ {
		// filter all parcels to this year only
		parcelList, err := parcelFilterByYear(AllParcels, y)
		if err != nil {
			return err
		}

		irrCells, err := database.GetCellsIrr(v, y)
		if err != nil {
			return err
		}

		// use the RO + DP from parcel and split by acres to get recharge, will need to keep separate files for the various
		// distributions of scenarios.
		for i := 0; i < len(irrCells); i++ {
			p, err := parcelFilterById(parcelList, irrCells[i].ParcelId)
			if err != nil {
				return err
			}

			_ = p

			// TODO: RO + DP by cell area to get value of cell recharge
			// TODO: get p values and send to a function that will assign the correct file integer classification and save to results, prorating for parcel size to parcel size within that cell
		}

		// TODO: Send to results database
	}

	return nil
}
