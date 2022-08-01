package actions

import (
	"fmt"
	"github.com/Longitude103/wwum2020/parcels"
	"time"

	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/fileio"
	"github.com/Longitude103/wwum2020/rchFiles"
	"github.com/pterm/pterm"
)

// RunSteadyState is a function that runs the model in Steady State Mode. This produces the following recharge file, but no
// well file is produced.
func RunSteadyState(mDesc, CSDir string, AvgStart, AvgEnd int, oldGrid, mf640 bool, myEnv map[string]string) error {
	// first stress period is all Natural Veg for whole grid
	// Second stress period is all natural veg and a repeat of the 1st period
	// Third -> end is a monthly stress periods using surface water data starting January 1895 to December 1952
	timeStart := time.Now()

	pterm.Info.Println("Setting up results database")
	var opts []database.Option
	opts = append(opts, database.WithSteadyState())

	if oldGrid {
		opts = append(opts, database.WithOldGrid())
	}

	if mf640 {
		opts = append(opts, database.WithMF640Grid())
	}

	v, err := database.NewSetup(myEnv, opts...)
	if err != nil {
		return err
	}

	v.Logger.Infof("Model Run Started at: %s", timeStart.Format(time.UnixDate))

	pterm.Info.Printf("Model Description: %s\n", mDesc)
	v.Logger.Infof("Model Description: %s", mDesc)
	v.SYear = AvgStart
	v.EYear = AvgEnd

	noteDb, err := database.ResultsNoteDB(v.SlDb)
	if err != nil {
		return err
	}

	sYearNote := fmt.Sprintf("Average Start Year: %d", v.SYear)
	eYearNote := fmt.Sprintf("Average End Year: %d", v.EYear)
	if err := noteDb.Add(database.Note{Nt: "Desc: " + mDesc}); err != nil {
		return err
	}
	if err := noteDb.Add(database.Note{Nt: sYearNote}); err != nil {
		return err
	}
	if err := noteDb.Add(database.Note{Nt: eYearNote}); err != nil {
		return err
	}

	if v.OldGrid {
		if err := noteDb.Add(database.Note{Nt: "grid=1"}); err != nil {
			return err
		}
		if err := database.AddCellsToOutput(v); err != nil {
			return err
		}
	} else {
		if err := noteDb.Add(database.Note{Nt: "grid=2"}); err != nil {
			return err
		}
	}

	// get cropsim files loaded
	spinnerSuccess, _ := pterm.DefaultSpinner.Start("Reading CropSim Results files")
	csResults, err := fileio.LoadTextFiles(CSDir, v.Logger)
	if err != nil {
		spinnerSuccess.Fail("Error in Loading Text Files")
		return err
	}

	v.Logger.Info("Averaging CS Results")
	avgCSResults, err := fileio.AverageStationResults(csResults, AvgStart, AvgEnd)
	if err != nil {
		spinnerSuccess.Fail("Error in Averaging CS Results")
		return err
	}

	spinnerSuccess.Success()

	// get Weather Stations
	spinnerSuccess, _ = pterm.DefaultSpinner.Start("Getting Weather Stations")
	v.Logger.Info("Getting Weather Stations")
	wStations, err := database.GetWeatherStations(v.PgDb)
	if err != nil {
		spinnerSuccess.Fail("Error Getting Weather Stations")
		return err
	}
	spinnerSuccess.Success()

	spinnerSuccess, _ = pterm.DefaultSpinner.Start("Getting Coefficients of Crops")
	v.Logger.Info("Getting Coefficients of Crops")
	cCoefficients, err := database.GetCoeffCrops(v)
	if err != nil {
		return err
	}
	spinnerSuccess.Success()

	// create average natural vegetation values for each month
	v.Logger.Info("Preforming Natural Vegetation Calculations")
	if err := rchFiles.NaturalVegSS(v, wStations, avgCSResults, cCoefficients); err != nil {
		v.Logger.Errorf("Error in Natural Vegetation: %s", err)
		return err
	}

	// parcel pumping
	v.Logger.Info("Preforming Parcel Pumping")
	irrParcels, err := parcels.ParcelPumpSS(v, avgCSResults, wStations, cCoefficients)
	if err != nil {
		v.Logger.Errorf("Error in Parcel Pumping: %s", err)
	}
	pterm.Success.Println("Successfully Completed Parcel Pumping")

	_ = irrParcels

	return nil
}