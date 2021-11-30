package parcels

import (
	"fmt"
	"sync"
	"time"

	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/fileio"
	"github.com/Longitude103/wwum2020/parcels/conveyLoss"
	"github.com/pterm/pterm"
)

// ParcelPump is the main function for the parcels, it gets the usage, efficiencies, operates the surface water convayance
// loss and then calls the surface water delivery. I also creates the parcels then calls the ParcelPumpDB method to set
// the parcel pumping, it then loops through the years for each parcel and sends the diversions, calls parcel NIR, sets the
// efficiency for the parcel, adds SW delivery, adds the known pumping, and then calls the simulate pumping for all other
// parcels. Finally it writes out the pumping per parcel and then operates the WSPP routine to finish the RO and DP.
func ParcelPump(v *database.Setup, csResults map[string][]fileio.StationResults,
	wStations []database.WeatherStation, cCrops []database.CoeffCrop) (AllParcels []Parcel, err error) {

	spinner, _ := pterm.DefaultSpinner.Start("Getting Cert Usage and Efficiencies")
	// cert usage
	v.Logger.Info("Getting Cert Usage")
	usage := getUsage(v)

	v.Logger.Info("Getting Efficiencies")
	efficiencies := database.GetAppEfficiency(v.PgDb)
	spinner.Success()

	v.Logger.Info("Running Conveyance Loss")
	err = conveyLoss.Conveyance(v)
	if err != nil {
		v.Logger.Errorf("Error in Conveyance Losses %s", err)
	}

	spinner, _ = pterm.DefaultSpinner.Start("Getting Surface Water Delivery")
	// parcel delivery
	v.Logger.Info("Getting Surface Water Delivery")
	swDelivery, err := conveyLoss.GetSurfaceWaterDelivery(v)
	if err != nil {
		spinner.Fail("Error in Surface Water Delivery")
		return nil, err
	}
	spinner.Success()

	var parcels []Parcel

	var pPumpDB *database.PPDB
	if !v.AppDebug {
		spinner, _ = pterm.DefaultSpinner.Start("Setting Parcel Pumping")
		v.Logger.Info("Setting Parcel Pumping")
		pPumpDB, err = database.ParcelPumpDB(v.SlDb)
		if err != nil {
			spinner.Fail("Failed Setting Parcel Pumping")
			return nil, err
		}
		spinner.Success()

		defer func(pPumpDB *database.PPDB) {
			err := pPumpDB.Close()
			if err != nil {
				return
			}
		}(pPumpDB)
	}

	// 1. load parcels
	p, _ := pterm.DefaultProgressbar.WithTotal(v.EYear - v.SYear + 1).WithTitle("Parcel Operations").WithRemoveWhenDone(true).Start()
	wg := sync.WaitGroup{}

	for y := v.SYear; y < v.EYear+1; y++ {
		p.UpdateTitle(fmt.Sprintf("Getting %d Parels", y))
		parcels = GetParcels(v, y)

		p.UpdateTitle(fmt.Sprintf("Calculating %d Parcel NIR", y))
		for i := 0; i < len(parcels); i++ {
			wg.Add(1)
			go func(ip int) {
				defer wg.Done()
				err := (&parcels[ip]).ParcelNIR(v, y, wStations, csResults, Irrigated)
				if err != nil {
					v.Logger.Errorf("Parcel NIR Error: %s", err)
					v.Logger.Errorf("Parcel Trace: %+v", parcels[ip])
				}

				(&parcels[ip]).setAppEfficiency(efficiencies, y)

				// add SW Delivery to the parcels
				if parcels[ip].Sw.Bool {
					(&parcels[ip]).parcelSWDelivery(swDelivery[y])
				}
			}(i) // must be a pointer to work

			wg.Wait()
		}

		// add usage to parcel
		p.UpdateTitle(fmt.Sprintf("Calculating %d Parcel Annual Usage", y))
		v.Logger.Info("Setting Annual Usage")
		if err := distUsage(usage[y], &parcels); err != nil {
			return []Parcel{}, err
		}

		// get all parcels simulate pumping if GW == true
		p.UpdateTitle(fmt.Sprintf("Simulating %d Parcel Pumping", y))
		v.Logger.Infof("Simulating Pumping for year %d", y)
		for p := 0; p < len(parcels); p++ {
			if (&parcels[p]).Gw.Bool {
				if err := (&parcels[p]).estimatePumping(v, cCrops); err != nil {
					return []Parcel{}, err
				}
			}
		}

		if !v.AppDebug {
			// write out parcel pumping for each parcel in sqlite results
			p.UpdateTitle(fmt.Sprintf("Saving %d Parcel Pumping Results", y))
			for p := 0; p < len(parcels); p++ {
				if parcels[p].Gw.Bool {
					// Add data to pumpingStruct and then append
					for m := 1; m < 13; m++ {
						if parcels[p].Pump[m-1] > 0 {
							dt := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
							_ = pPumpDB.Add(database.Pumping{ParcelID: parcels[p].ParcelNo, Nrd: parcels[p].Nrd, Dt: dt,
								Pump: parcels[p].Pump[m-1]})
						}
					}
				}
			}
		}

		v.Logger.Infof("Simulating parcel WSPP for year %d", y)
		p.UpdateTitle(fmt.Sprintf("Calculating %d Parcel WSPP", y))
		for p := 0; p < len(parcels); p++ {
			//wg.Add(1)
			//go func(i int) {
			//	defer wg.Done()
			//	err := (&parcels[i]).waterBalanceWSPP(false)
			//	if err != nil {
			//		v.Logger.Errorf("error in parcel WSPP parcel data: %+v", parcels[p])
			//	}
			//}(p)

			err := (&parcels[p]).waterBalanceWSPP(v)
			if err != nil {
				v.Logger.Errorf("error in parcel WSPP: %s\n", err)
				v.Logger.Errorf("Parcel trace: %+v\n", parcels[p])
			}

			AllParcels = append(AllParcels, parcels[p])
		}
		//wg.Wait()
		p.Increment()
	}

	return AllParcels, nil
}

func distUsage(annUsage []Usage, parcels *[]Parcel) error {
	for _, u := range annUsage {
		// filter parcels to this usage cert
		filteredParcels := FilterParcelByCert(parcels, u.CertNum)

		totalNIR := 0.0
		totalMonthlyNIR := [12]float64{}

		for _, pIndex := range filteredParcels {
			for m := 0; m < 12; m++ {
				totalMonthlyNIR[m] += (*parcels)[pIndex].Nir[m]
				totalNIR += (*parcels)[pIndex].Nir[m]
			}
		}

		for _, pIndex := range filteredParcels {
			(*parcels)[pIndex].distributeUsage(totalNIR, totalMonthlyNIR, u.UseAF)
		}
	}

	return nil
}
