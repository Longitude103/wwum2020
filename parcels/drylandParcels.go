package parcels

import (
	"fmt"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/fileio"
	"github.com/pterm/pterm"
	"sync"
)

// DryLandParcels is a function that returns all the dryland parcels for the years of the simulation and also calls the
// methods to determine parcelNIR and dryland WSPP
func DryLandParcels(v *database.Setup, csResults map[string][]fileio.StationResults,
	wStations []database.WeatherStation, cCrop []database.CoeffCrop) (dryParcels []Parcel, err error) {

	p, _ := pterm.DefaultProgressbar.WithTotal(v.EYear - v.SYear + 1).WithTitle("Dryland Parcel Operations").WithRemoveWhenDone(true).Start()
	v.Logger.Info("Getting Dryland parcels")
	for y := v.SYear; y < v.EYear+1; y++ {
		p.UpdateTitle(fmt.Sprintf("Getting %d Dryland Parcels", y))
		annDryParcels := getDryParcels(v, y)

		// method is used to set RO and DP, just poorly named.
		p.UpdateTitle(fmt.Sprintf("Calculating %d Dryland Parcels return flows", y))
		wg := sync.WaitGroup{}
		for i := 0; i < len(annDryParcels); i++ {
			wg.Add(1)
			go func(d int) {
				defer wg.Done()
				err := (&annDryParcels[d]).parcelNIR(v.PNirDB, y, wStations, csResults, DryLand)
				if err != nil {
					v.Logger.Errorf("error in dry parcel NIR:", err)
					v.Logger.Errorf("Parcel trace: %+v", annDryParcels[d])
				}
			}(i)

			//go func(d int) {
			//	defer wg.Done()
			//	err := (&annDryParcels[d]).dryWaterBalanceWSPP(cCrop)
			//	if err != nil {
			//		v.Logger.Error("error in dry parcel WSPP ", err)
			//		v.Logger.Errorf("Parcel trace: %+v", annDryParcels[d])
			//	}
			//}(i)
		}
		wg.Wait()
		dryParcels = append(dryParcels, annDryParcels...)
		p.Increment()
	}

	v.Logger.Info("Finished Dryland parcel operations")
	return dryParcels, nil
}
