package actions

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"
	"strconv"

	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
	"github.com/heath140/wwum2020/parcels"
	"github.com/heath140/wwum2020/rchFiles"
	//"wwum2020/rchFiles"
)

func RechargeFiles(debug *bool, CSDir *string) error {
	v := database.Setup{}
	if err := v.NewSetup(); err != nil {
		return err
	}

	csResults, err := fileio.LoadTextFiles(*CSDir, v.Logger)
	if err != nil {
		fmt.Println("Error in Loading Text Files, check log file")
		return err
	}

	//fmt.Println("CSResults in RCH")
	//fmt.Println(csResults)

	// TODO: Validate should only allow 1953 to current year
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			v.Logger.Errorf("Invalid number %s, error: %s", input, err)
			return errors.New("invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Start Year of Model Run",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		v.Logger.Errorf("Prompt failed %v", err)
		return err
	}

	startYr, _ := strconv.Atoi(result)

	prompt = promptui.Prompt{
		Label:    "End Year of Model Run",
		Validate: validate,
	}

	result, err = prompt.Run()
	if err != nil {
		v.Logger.Errorf("Prompt failed %v", err)
		return err
	}

	endYr, _ := strconv.Atoi(result)

	err = v.SetYears(startYr, endYr)
	if err != nil {
		v.Logger.Errorf("Error Setting Years Error: %s", err)
		return err
	}

	v.Logger.Info("Getting Weather Stations")
	wStations := database.GetWeatherStations(v.PgDb)

	// parcel pumping
	irrParcels, err := parcels.ParcelPump(v, csResults, wStations)
	if err != nil {
		v.Logger.Errorf("Error in Parcel Pumping: %s", err)
	}
	//
	//for i := 0; i < 10; i++ {
	//	fmt.Println(&irrParcels[i])
	//	fmt.Println(irrParcels[i].PrintNIR())
	//}
	_ = irrParcels

	err = v.PNirDB.Flush()
	if err != nil {
		v.Logger.Errorf("Error in flush: %s", err)
	}

	dryParcels, err := parcels.DryLandParcels(v, csResults, wStations)
	if err != nil {
		v.Logger.Errorf("Error in Dry Land Parcels: %s", err)
	}
	_ = dryParcels

	_ = v.PNirDB.Close() // close doesn't close the db, that must be call explicitly so we can keep using it.
	_ = v.SlDb.Close()   // close the db before ending the program

	return nil
	// load up data with cell acres
	cells := database.GetCells(v.PgDb)

	// will also need parcel sw delivery, gw pumping (if available), distributed nir, rf, eff precip for the required crops

	// Natural Veg 102
	//rchFiles.NaturalVeg(db, v.PgDb, debug, startYr, endYr)

	// Irr Cells
	irrCells := rchFiles.GetCellsIrr(v.PgDb, 2014)
	//fmt.Println("First Irrigated Cell:")
	//fmt.Println(irrCells[0])

	// Dry Cells
	dryCells := rchFiles.GetDryCells(v.PgDb, 2014)
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

	return nil
}
