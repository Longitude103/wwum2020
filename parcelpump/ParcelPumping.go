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

	// 3. pumping amounts / parcel
	// 1. load parcels
	for y := sYear; y < eYear+1; y++ {
		parcels := getParcels(pgDB, y)

		bar := progressbar.Default(int64(len(parcels)), "Parcels")

		for i := 0; i < len(parcels); i++ {
			(&parcels[i]).parcelNIR(slDB, y, wStations, *csResults) // must be a pointer to work
			_ = bar.Add(1)
		}

		_ = bar.Close()

		// add usage to parcel

		// add sw delivery to parcel

		// calculate / recalculate RO and DP for the parcel & estimate pumping for years without usage

		for _, v := range parcels[:10] {
			fmt.Printf("Parce No: %d, NIR is: %v\n", v.ParcelNo, v.Nir[sYear])
		}
	}

	// 4. parcel recharge / acre
}
