package actions

import (
	"fmt"
	"wwum2020/database"
	"wwum2020/rchFiles"
	//"wwum2020/rchFiles"
)

func RechargeFiles(debug *bool, startYr *int, endYr *int, CSDir *string) {
	db := database.GetSqlite()
	pgDb := database.PgConnx()

	_ = db
	// load up data with cell acres
	cells := database.GetCells(pgDb)

	// loop through the cells
	for _, cell := range cells[:5] {
		fmt.Println(cell.CellId)
	}

	// will also need parcel sw delivery, gw pumping (if available), distributed nir, rf, eff precip for the required crops

	// Natural Veg 102
	//rchFiles.NaturalVeg(db, pgDb, debug, startYr, endYr)

	// Irr Cells
	rchFiles.GetCellsIrr(pgDb, 2014)
}
