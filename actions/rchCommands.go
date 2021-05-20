package actions

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"

	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
	"github.com/heath140/wwum2020/parcels"
	"github.com/heath140/wwum2020/rchFiles"
	//"wwum2020/rchFiles"
)

func RechargeFiles(debug *bool, CSDir *string) {
	logger, _ := NewLogger()
	sugar := logger.Sugar()

	sugar.Infow("Setting Up Results database, getting postgres DB Connection.")
	slDb := database.GetSqlite(sugar)
	pgDb := database.PgConnx()

	csResults := fileio.LoadTextFiles(*CSDir, sugar)

	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			sugar.Errorf("Invalid number %s, error: %s", input, err)
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
		sugar.Errorf("Prompt failed %v", err)
		return
	}

	startYr, _ := strconv.Atoi(result)

	prompt = promptui.Prompt{
		Label:    "End Year of Model Run",
		Validate: validate,
	}

	result, err = prompt.Run()
	if err != nil {
		sugar.Errorf("Prompt failed %v", err)
		return
	}

	endYr, _ := strconv.Atoi(result)

	// parcel pumping
	irrParcels, err := parcels.ParcelPump(pgDb, slDb, startYr, endYr, &csResults, sugar)
	if err != nil {
		sugar.Errorf("Error in Parcel Pumping: %s", err)
	}
	_ = irrParcels

	dryParcels := parcels.DryLandParcels(pgDb, slDb, startYr, endYr, &csResults, sugar)
	_ = dryParcels

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

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	path := fmt.Sprintf("./results%s.log", time.Now().Format(time.RFC3339))

	cfg.OutputPaths = []string{
		path,
	}

	return cfg.Build()
}
