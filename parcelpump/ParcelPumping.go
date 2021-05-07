package parcelpump

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/schollz/progressbar"
	"wwum2020/database"
	"wwum2020/fileio"
)

func ParcelPump(pgDB *sqlx.DB, slDB *sqlx.DB, sYear int, eYear int, csResults *map[string][]fileio.StationResults) {
	// cert usage
	usage := getUsage(pgDB)
	_ = usage

	wStations := database.GetWeatherStations(pgDB)

	// 2. sw deliveries / canal recharge
	// 3. pumping amounts / parcel
	// 1. load parcels
	for y := sYear; y < eYear+1; y++ {
		parcels := getParcels(pgDB, y)

		bar := progressbar.Default(int64(len(parcels)), "Parcel NIR")

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
