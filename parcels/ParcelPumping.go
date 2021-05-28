package parcels

import (
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
	"github.com/heath140/wwum2020/parcels/conveyLoss"
	"github.com/schollz/progressbar/v3"
	"time"
)

func ParcelPump(v database.Setup, csResults map[string][]fileio.StationResults,
	wStations []database.WeatherStation) (AllParcels []Parcel, err error) {
	// cert usage
	v.Logger.Info("Getting Cert Usage")
	usage := getUsage(v.PgDb)
	_ = usage

	v.Logger.Info("Getting CoeffCrops Data")
	cCrops, err := database.GetCoeffCrops(v.PgDb)
	if err != nil {
		v.Logger.Error("Cannot get Coefficient of Crops")
		return nil, err
	}

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
	for y := v.SYear; y < v.EYear+1; y++ {
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
		annUsage := filterUsage(usage, y)
		for _, u := range annUsage {
			//fmt.Printf("Annual Usage in %v is %g\n", u.CertNum, u.UseAF)
			// filter parcels to this usage cert
			filteredParcels := filterParcelByCert(&parcels, u.CertNum)

			totalNIR := 0.0
			totalMonthlyNIR := [12]float64{}

			for i := 0; i < len(filteredParcels); i++ {
				for m := 0; m < 12; m++ {
					totalMonthlyNIR[m] += filteredParcels[i].Nir[m]
					totalNIR += filteredParcels[i].Nir[m]
				}
			}

			for i := 0; i < len(filteredParcels); i++ {
				filteredParcels[i].distributeUsage(totalNIR, totalMonthlyNIR, u.UseAF)
			}
		}

		// get all parcels where Metered == false and simulate pumping if GW == true
		for p := 0; p < len(parcels); p++ {
			if (&parcels[p]).Metered == false && (&parcels[p]).Gw.Bool == true {
				(&parcels[p]).estimatePumping(cCrops)
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

		// TODO: Water Balance the parcel where if SW or GW or Both was over / under applied, then adjust RO and DP
		// to account for the difference. Save that back to the parcel to simplify the RCH file Creation later.

		AllParcels = append(AllParcels, parcels...)
	}

	return AllParcels, nil
}
