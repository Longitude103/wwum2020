package rchFiles

import (
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/parcels"
)

func IrrigationRCH(v database.Setup, AllParcels []parcels.Parcel) error {

	for y := v.SYear; y < v.EYear+1; y++ {
		irrCells, err := database.GetCellsIrr(v, y)
		if err != nil {
			return err
		}

		// TODO: Divide parcel RO + DP by cell area to get value of cell recharge
		// use the RO + DP from parcel and split by acres to get recharge, will need to keep separate files for the various
		// distributions of scenarios.
		_ = irrCells

		// TODO: Send to results database
	}

	return nil
}
