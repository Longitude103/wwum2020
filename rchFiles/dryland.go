package rchFiles

import (
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/parcels"
)

func dryland(v database.Setup, dryParcels []parcels.Parcel) error {
	// get cells with parcels within it
	dryCells := database.GetDryCells(v, 2014)
	_ = dryCells

	// TODO: Finish this section.
	// find parcel in each cell
	// add cell Dp + Ro to Output Database
	// filetype: 101

	return nil
}
