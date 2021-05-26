package parcels

import (
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
)

func DryLandParcels(v database.Setup, csResults map[string][]fileio.StationResults,
	wStations []database.WeatherStation) (dryParcels []Parcel, err error) {

	v.Logger.Info("Getting parcels")
	for y := v.SYear; y < v.EYear+1; y++ {
		dryParcels = getDryParcels(v, y)

		// method is used to set RO and DP, just poorly named.
		for i := 0; i < len(dryParcels); i++ {
			err = (&dryParcels[i]).parcelNIR(v.PNirDB, y, wStations, csResults, DryLand)
		}
		if err != nil {
			return nil, err
		}
	}

	return dryParcels, nil
}
