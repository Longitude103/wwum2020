package rchFiles

import (
	"database/sql"
	"fmt"
	"github.com/heath140/wwum2020/database"
	"github.com/jmoiron/sqlx"
)

func NaturalVeg(sqliteDB *sql.DB, pgDB *sqlx.DB, debug *bool, sYear *int, eYear *int) {
	fmt.Println("welcome to NaturalVeg")
	fmt.Println("Debug is set to", *debug)

	cells := database.GetCells(pgDB)

	fmt.Println("First Cell ID", cells[0].CellId)

	//for yr := *sYear; yr < *eYear; yr++ {
	//	fmt.Println("Year", yr)
	//}
}
