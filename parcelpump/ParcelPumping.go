package parcelpump

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func ParcelPump(pgDB *sqlx.DB, slDB *sqlx.DB, sYear int, eYear int) {
	// 1. load parcels
	parcels := getParcels(pgDB, sYear, eYear)
	for _, v := range parcels[:5] {
		fmt.Println(v)
	}

	// cert usage
	usage := getUsage(pgDB)
	for _, v := range usage[:5] {
		fmt.Println(v)
	}

	// 2. sw deliveries / canal recharge
	// 3. pumping amounts / parcel
	// 4. parcel recharge / acre
}
