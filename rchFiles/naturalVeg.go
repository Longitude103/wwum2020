package rchFiles

import (
	"database/sql"
	"fmt"
	"github.com/heath140/wwum2020/database"
	"github.com/jmoiron/sqlx"
)

func NaturalVeg(sqliteDB *sql.DB, pgDB *sqlx.DB, debug *bool, sYear *int, eYear *int) error {
	fmt.Println("welcome to NaturalVeg")
	fmt.Println("Debug is set to", *debug)

	cells, err := database.GetCells(pgDB)
	if err != nil {
		return err
	}

	fmt.Println("First Cell ID", cells[0].Node)

	//for yr := *sYear; yr < *eYear; yr++ {
	//	fmt.Println("Year", yr)
	//}

	return nil
}
