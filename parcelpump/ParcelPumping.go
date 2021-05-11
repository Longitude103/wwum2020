package parcelpump

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"

	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
	"github.com/heath140/wwum2020/parcelpump/conveyLoss"
)

func ParcelPump(pgDB *sqlx.DB, slDB *sqlx.DB, sYear int, eYear int, csResults *map[string][]fileio.StationResults) {
	// cert usage
	usage := getUsage(pgDB)
	_ = usage

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

	conveyLoss.Conveyance(pgDB)

	os.Exit(0)

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

		for _, v := range parcels[:10] {
			fmt.Printf("Parce No: %d, NIR is: %v\n", v.ParcelNo, v.Nir[sYear])
		}
	}

	// 4. parcel recharge / acre
}
