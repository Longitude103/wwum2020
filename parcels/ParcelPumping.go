package parcels

import (
	"fmt"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/fileio"
	"github.com/Longitude103/wwum2020/parcels/conveyLoss"
	"github.com/schollz/progressbar/v3"
	"os"
	"time"
)

// ParcelPump is the main function for the parcels, it gets the usage, efficiencies, operates the surface water convayance
// loss and then calls the surface water delivery. I also creates the parcels then calls the ParcelPumpDB method to set
// the parcel pumping, it then loops through the years for each parcel and sends the diversions, calls parcel NIR, sets the
// efficiency for the parcel, adds SW delivery, adds the known pumping, and then calls the simulate pumping for all other
// parcels. Finally it writes out the pumping per parcel and then operates the WSPP routine to finish the RO and DP.
func ParcelPump(v database.Setup, csResults map[string][]fileio.StationResults,
	wStations []database.WeatherStation, cCrops []database.CoeffCrop) (AllParcels []Parcel, err error) {
	// cert usage
	v.Logger.Info("Getting Cert Usage")
	usage := getUsage(v.PgDb)

	v.Logger.Info("Getting Efficiencies")
	efficiencies := database.GetAppEfficiency(v.PgDb)

	v.Logger.Info("Running Conveyance Loss")
	err = conveyLoss.Conveyance(v)
	if err != nil {
		v.Logger.Errorf("Error in Conveyance Losses %s", err)
	}

	// parcel delivery
	v.Logger.Info("Getting Surface Water Delivery")
	swDelivery, err := conveyLoss.GetSurfaceWaterDelivery(v)
	if err != nil {
		return nil, err
	}

	var parcels []Parcel

	v.Logger.Info("Setting Parcel Pumping")
	pPumpDB, err := database.ParcelPumpDB(v.SlDb)
	if err != nil {
		return nil, err
	}

	defer func(pPumpDB *database.PPDB) {
		err := pPumpDB.Close()
		if err != nil {
			return
		}
	}(pPumpDB)

	// 1. load parcels
	parcelYearBar := progressbar.Default(int64(v.EYear-v.SYear), "Years of Parcels")
	for y := v.SYear; y < v.EYear+1; y++ {
		_ = parcelYearBar.Add(1)
		parcels = getParcels(v, y)
		filteredDiversions := conveyLoss.FilterSWDeliveryByYear(swDelivery, y)

		bar := progressbar.Default(int64(len(parcels)), "Parcels")
		for i := 0; i < len(parcels); i++ {
			err = (&parcels[i]).parcelNIR(v.PNirDB, y, wStations, csResults, Irrigated) // must be a pointer to work
			if err != nil {
				return nil, err
			}
			(&parcels[i]).setAppEfficiency(efficiencies, y)

			// add SW Delivery to the parcels
			if parcels[i].Sw.Bool == true {
				(&parcels[i]).parcelSWDelivery(filteredDiversions)
			}

			_ = bar.Add(1)
		}

		_ = bar.Close()

		// add usage to parcel
		v.Logger.Info("Setting Annual Usage")
		annUsage := filterUsage(usage, y)
		// TODO: check results of this function and write more tests.
		if err := distUsage(annUsage, &parcels); err != nil {
			return []Parcel{}, err
		}

		// get all parcels simulate pumping if GW == true
		v.Logger.Infof("Simulating Pumping for year %d", y)
		for p := 0; p < len(parcels); p++ {
			if (&parcels[p]).Gw.Bool == true {
				// TODO: Check to make sure this function is working
				if err := (&parcels[p]).estimatePumping(cCrops); err != nil {
					return []Parcel{}, err
				}
			}
		}

		// write out parcel pumping for each parcel in sqlite results
		for p := 0; p < len(parcels); p++ {
			if parcels[p].Gw.Bool == true {
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

		v.Logger.Infof("Simulating parcel WSPP for year %d", y)
		wbBar := progressbar.Default(int64(len(parcels)), "Water Balance Parcels")
		for p := 0; p < len(parcels); p++ {
			_ = wbBar.Add(1)
			if err := (&parcels[p]).waterBalanceWSPP(false); err != nil {
				v.Logger.Errorf("error in parcel WSPP parcel data: %+v", parcels[p])
				return nil, err
			}

			AllParcels = append(AllParcels, parcels[p])
		}
		_ = wbBar.Close()

	}
	_ = parcelYearBar.Close()

	for i := 0; i < len(AllParcels); i++ {
		fmt.Println("All Parcels", AllParcels[i])
		if i > 40 {
			os.Exit(1)
		}
	}

	return AllParcels, nil
}

func distUsage(annUsage []Usage, parcels *[]Parcel) error {
	for _, u := range annUsage {
		//fmt.Printf("Annual Usage in %v is %g\n", u.CertNum, u.UseAF)
		// filter parcels to this usage cert
		filteredParcels := filterParcelByCert(parcels, u.CertNum)

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
