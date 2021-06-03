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

		// TODO: pull out ETMAXDRY = dryET from CS
		// TODO: Create ETMAXDRYAdj and set to ETMAXDRY * Adjustment
		// TODO: RO3 = ETMAXDRY - ETMAXDryAdj * DryETtoRO
		// TODO: DP3 = (ETMAXDRY - ETMAXADJU) - RO3
		// TODO: RO = RO1 + RO3
		// TODO: DP = DP1 + DP3
	}

	return dryParcels, nil
}
