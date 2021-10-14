package parcels

import (
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/fileio"
	"sync"
)

// DryLandParcels is a function that returns all the dryland parcels for the years of the simulation and also calls the
// methods to determine parcelNIR and dryland WSPP
func DryLandParcels(v database.Setup, csResults map[string][]fileio.StationResults,
	wStations []database.WeatherStation, cCrop []database.CoeffCrop) (dryParcels []Parcel, err error) {

	v.Logger.Info("Getting parcels")
	for y := v.SYear; y < v.EYear+1; y++ {
		annDryParcels := getDryParcels(v, y)

		// method is used to set RO and DP, just poorly named.
		wg := sync.WaitGroup{}
		for i := 0; i < len(annDryParcels); i++ {
			wg.Add(1)
			go func(d int) {
				err := (&annDryParcels[d]).parcelNIR(v.PNirDB, y, wStations, csResults, DryLand)
				if err != nil {
					v.Logger.Error("error in dry parcel NIR ", err)
				}
			}(i)

			go func(d int) {
				defer wg.Done()
				err := (&annDryParcels[d]).dryWaterBalanceWSPP(cCrop)
				if err != nil {
					v.Logger.Error("error in dry parcel WSPP ", err)
				}
			}(i)
		}
		wg.Wait()

		dryParcels = append(dryParcels, annDryParcels...)
	}

	v.Logger.Info("Finished Dryland parcel operations")
	return dryParcels, nil
}
