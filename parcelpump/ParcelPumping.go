package parcelpump

import (
	"fmt"
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
	"github.com/heath140/wwum2020/parcelpump/conveyLoss"
	"github.com/jmoiron/sqlx"
	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"
)

func ParcelPump(pgDB *sqlx.DB, slDB *sqlx.DB, sYear int, eYear int, csResults *map[string][]fileio.StationResults) {
	// cert usage
	fmt.Println("Getting Parcel usage")
	usage := getUsage(pgDB)
	_ = usage

	fmt.Println("Getting Weather Stations")
	wStations := database.GetWeatherStations(pgDB)

	fmt.Println("Getting CoeffCrops Data")
	cCrops := database.GetCoeffCrops(pgDB)

	fmt.Println("Getting Efficiencies")
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
		fmt.Println("Including excess flows")
	}

	fmt.Println("Running Conveyance Loss")
	err = conveyLoss.Conveyance(pgDB, slDB, sYear, eYear, excessFlows)
	if err != nil {
		fmt.Println("Error in Conveyance Loss", err)
	}

	// parcel delivery
	swDelivery := conveyLoss.GetSurfaceWaterDelivery(pgDB, sYear, eYear)
	fmt.Println("First 10 Surface Water Delivery Records")
	for _, v := range swDelivery[:10] {
		fmt.Println(v)
	}

	// 1. load parcels
	for y := sYear; y < eYear+1; y++ {
		parcels := getParcels(pgDB, y)
		filteredDiversions := conveyLoss.FilterSWDeliveryByYear(swDelivery, y)

		bar := progressbar.Default(int64(len(parcels)), "Parcels")

		for i := 0; i < len(parcels); i++ {
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

		// calculate / recalculate RO and DP for the parcel & estimate pumping for years without usage
		for _, v := range parcels[:10] {
			fmt.Printf("Parce No: %d, NIR is: %v, Dp is: %v, Usage is: %v\n", v.ParcelNo, v.Nir, v.Dp, v.Pump)
		}
	}

	// 4. parcel recharge / acre
}
