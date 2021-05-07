package actions

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar"
	"os"
	"strconv"
	"wwum2020/database"
	"wwum2020/fileio"
	"wwum2020/parcelpump"
	"wwum2020/rchFiles"
	//"wwum2020/rchFiles"
)

func RechargeFiles(debug *bool, CSDir *string) {
	slDb := database.GetSqlite()
	pgDb := database.PgConnx()

	csResults := fileio.LoadTextFiles(*CSDir)

	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Start Year",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	startYr, _ := strconv.Atoi(result)

	prompt = promptui.Prompt{
		Label:    "End Year",
		Validate: validate,
	}

	result, err = prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	endYr, _ := strconv.Atoi(result)

	// parcel pumping
	parcelpump.ParcelPump(pgDb, slDb, startYr, endYr, &csResults)
	os.Exit(0)

	// load up data with cell acres
	cells := database.GetCells(pgDb)

	// will also need parcel sw delivery, gw pumping (if available), distributed nir, rf, eff precip for the required crops

	// Natural Veg 102
	//rchFiles.NaturalVeg(db, pgDb, debug, startYr, endYr)

	// Irr Cells
	irrCells := rchFiles.GetCellsIrr(pgDb, 2014)
	//fmt.Println("First Irrigated Cell:")
	//fmt.Println(irrCells[0])

	// Dry Cells
	dryCells := rchFiles.GetDryCells(pgDb, 2014)
	//fmt.Println("First Dry Cell:")
	//fmt.Println(dryCells[0])

	// progress bar
	bar := progressbar.Default(int64(len(cells)))

	// loop through the cells
	for _, cell := range cells {
		//fmt.Println(cell.CellId)

		var irrCellsResult []rchFiles.IrrCell
		var dryCellResult []rchFiles.DryCell
		totalParcelAcres := 0.0

		// filter irrCells for this cell also get acres and add to total
		for _, ic := range irrCells {
			if ic.CellId == cell.CellId {
				irrCellsResult = append(irrCellsResult, ic)
				totalParcelAcres += ic.IrrArea
			}
		}

		// filter dryCells for this cell also get acres and add to total
		for _, dc := range dryCells {
			if dc.CellId == cell.CellId {
				dryCellResult = append(dryCellResult, dc)
				totalParcelAcres += dc.DryArea
			}
		}

		if cell.CellId == 78585 {
			fmt.Printf("CellId: %d, Total Parcel Acres: %g\n", cell.CellId, totalParcelAcres)
		}

		bar.Add(1)
	}
}
