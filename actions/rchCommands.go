package actions

import (
	"wwum2020/database"
	"wwum2020/rchFiles"
)

func RechargeFiles(debug *bool, startYr *int, endYr *int, CSDir *string) {
	db := database.GetSqlite()

	// Natural Veg 102
	rchFiles.NaturalVeg(db, debug)

}
