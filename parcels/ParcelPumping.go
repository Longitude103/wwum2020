package parcels

import (
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
	"github.com/heath140/wwum2020/parcels/conveyLoss"
	"github.com/jmoiron/sqlx"
	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
	"time"
)

func ParcelPump(pgDB *sqlx.DB, slDB *sqlx.DB, sYear int, eYear int, csResults *map[string][]fileio.StationResults, logger *zap.SugaredLogger) (AllParcels []Parcel) {
	// cert usage
	logger.Info("Getting Cert Usage")
	usage := getUsage(pgDB)
	_ = usage

	logger.Info("Getting Weather Stations")
	wStations := database.GetWeatherStations(pgDB)

	logger.Info("Getting CoeffCrops Data")
	cCrops := database.GetCoeffCrops(pgDB)

	logger.Info("Getting Efficiencies")
	efficiencies := database.GetAppEfficiency(pgDB)

	// 2. sw deliveries / canal recharge
	prompt := promptui.Prompt{
		Label:     "Don't include Excess Flows",
		IsConfirm: true,
		Default:   "y",
	}

	excessFlows := false
	_, err := prompt.Run()
	if err != nil {
		// don't include excess flows
		excessFlows = true
	}

	if excessFlows {
		logger.Info("Including excess flows")
	}

	logger.Info("Running Conveyance Loss")
	err = conveyLoss.Conveyance(pgDB, slDB, sYear, eYear, excessFlows)
	if err != nil {
		logger.Errorf("Error in Conveyance Losses %s", err)
	}

	// parcel delivery
	logger.Info("Getting Surface Water Delivery")
	swDelivery := conveyLoss.GetSurfaceWaterDelivery(pgDB, sYear, eYear)

	var parcels []Parcel
	// 1. load parcels
	for y := sYear; y < eYear+1; y++ {
		parcels = getParcels(pgDB, y, logger)
		filteredDiversions := conveyLoss.FilterSWDeliveryByYear(swDelivery, y)

		bar := progressbar.Default(int64(len(parcels)), "Parcels")

		//for i := 0; i < len(parcels); i++ {
		for i := 0; i < 50; i++ {

			(&parcels[i]).parcelNIR(slDB, y, wStations, *csResults) // must be a pointer to work
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
		var pumpingOutput []Pumping
		for p := 0; p < len(parcels); p++ {
			if parcels[p].Gw.Bool == true {
				// Add data to pumpingStruct and then append
				for m := 1; m < 13; m++ {
					if parcels[p].Pump[m-1] > 0 {
						dt := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
						pmp := Pumping{ParcelID: parcels[p].ParcelNo, Nrd: parcels[p].Nrd, Dt: dt, Pump: parcels[p].Pump[m-1]}
						pumpingOutput = append(pumpingOutput, pmp)
					}
				}

				if p*12 > 500 {
					_ = bulkSaveSqlite(slDB, pumpingOutput, logger)
					pumpingOutput = nil
				}
			}
		}

		// save remaining
		if len(pumpingOutput) > 0 {
			_ = bulkSaveSqlite(slDB, pumpingOutput, logger)
		}

		AllParcels = append(AllParcels, parcels...)
	}

	return AllParcels
}
