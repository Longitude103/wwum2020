package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/fileio"
	"github.com/Longitude103/wwum2020/parcels"
	"github.com/Longitude103/wwum2020/rchFiles"
	"github.com/Longitude103/wwum2020/wells"
	"github.com/pterm/pterm"
)

var runModelCmd = &cobra.Command{
	Use:   "runModel",
	Short: "Run the EscModel",
	Long:  `This command is the main execution command of EscModel. This function will run the model using the provided configuration. It will create a new output database file on your local filesystem in the same directory as the executable.`,
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		sYr, _ := cmd.Flags().GetInt("StartYr")
		eYr, _ := cmd.Flags().GetInt("EndYr")
		csDir, _ := cmd.Flags().GetString("CSDir")
		ef, _ := cmd.Flags().GetBool("noExcessFlow")
		desc, _ := cmd.Flags().GetString("Desc")
		post97, _ := cmd.Flags().GetBool("post97")
		oldGrid, _ := cmd.Flags().GetBool("oldGrid")
		mf6Grid40, _ := cmd.Flags().GetBool("mf6Grid40")

		if oldGrid && mf6Grid40 {
			pterm.Error.Println("oldGrid and mf6Grid40 are mutually exclusive")
			return
		}

		if err := runModel(debug, &csDir, desc, sYr, eYr, ef, post97, oldGrid, mf6Grid40, myEnv); err != nil {
			pterm.Error.Printf("Error in Application: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runModelCmd)
	runModelCmd.Flags().BoolP("debug", "b", false, "sets debugger for more log information")
	runModelCmd.Flags().IntP("StartYr", "s", 1997, "Sets the start year of Command")
	runModelCmd.Flags().IntP("EndYr", "e", 2020, "Sets the end year of Command")
	runModelCmd.Flags().StringP("CSDir", "c", "", "REQUIRED CropSim Directory path")
	runModelCmd.Flags().BoolP("noExcessFlow", "n", false, "Sets to use Excess Flow or Not")
	runModelCmd.Flags().StringP("Desc", "d", "", "REQUIRED Model Description")
	runModelCmd.Flags().BoolP("post97", "p", false, "If flag set, a post 97 run will be made")
	runModelCmd.Flags().BoolP("oldGrid", "o", false, "If flag set, the model will use the 40 acre grid, not USG as default")
	runModelCmd.Flags().BoolP("mf6Grid40", "m", false, "If flag set, the model will use the 40 acre grid but in MF6 Node Numbers")
	_ = runModelCmd.MarkFlagRequired("CSDir")
	_ = runModelCmd.MarkFlagRequired("Desc")
}

func runModel(debug bool, CSDir *string, mDesc string, sY int, eY int, eF bool, p97 bool, oldGrid bool, mf640 bool, myEnv map[string]string) error {
	timeStart := time.Now()

	pterm.Info.Println("Setting up results database")
	var opts []database.Option
	if debug {
		opts = append(opts, database.WithDebug())
	}

	if eF {
		opts = append(opts, database.WithNoExcessFlow())
	}

	if p97 {
		opts = append(opts, database.WithPost97())
	}

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
	if err := v.SetYears(sY, eY); err != nil {
		v.Logger.Errorf("Error Setting Years: %s", err)
		return err
	}

	noteDb, err := database.ResultsNoteDB(v.SlDb)
	if err != nil {
		return err
	}

	sYearNote := fmt.Sprintf("Start Year: %d", v.SYear)
	eYearNote := fmt.Sprintf("End Year: %d", v.EYear)
	if err := noteDb.Add(database.Note{Nt: "Desc: " + mDesc}); err != nil {
		return err
	}
	if err := noteDb.Add(database.Note{Nt: sYearNote}); err != nil {
		return err
	}
	if err := noteDb.Add(database.Note{Nt: eYearNote}); err != nil {
		return err
	}

	if v.Post97 {
		if err := noteDb.Add(database.Note{Nt: "In Post 97 Mode"}); err != nil {
			return err
		}
	}

	if v.OldGrid {
		if err := noteDb.Add(database.Note{Nt: "grid=1"}); err != nil {
			return err
		}
		if err := database.AddCellsToOutput(v); err != nil { // TODO: Check, may need run for all models
			return err
		}
	} else {
		if err := noteDb.Add(database.Note{Nt: "grid=2"}); err != nil {
			return err
		}
	}

	if eF {
		if err := noteDb.Add(database.Note{Nt: "Includes Excess Flow Diversions"}); err != nil {
			return err
		}
	}

	spinnerSuccess, _ := pterm.DefaultSpinner.Start("Reading CropSim Results files")
	csResults, err := fileio.LoadTextFiles(*CSDir, v.Logger)
	if err != nil {
		fmt.Println("Error in Loading Text Files, check log file")
		return err
	}
	spinnerSuccess.Success()

	spinnerSuccess, _ = pterm.DefaultSpinner.Start("Getting Weather Stations")
	v.Logger.Info("Getting Weather Stations")
	wStations, err := database.GetWeatherStations(v.PgDb)
	if err != nil {
		v.Logger.Errorf("Error Getting Weather stations: %s", err)
	}
	spinnerSuccess.Success()

	spinnerSuccess, _ = pterm.DefaultSpinner.Start("Getting Coefficients of Crops")
	v.Logger.Info("Getting Coefficients of Crops")
	cCoefficients, err := database.GetCoeffCrops(v)
	if err != nil {
		return err
	}
	spinnerSuccess.Success()

	// parcel pumping
	v.Logger.Info("Preforming Parcel Pumping")
	irrParcels, err := parcels.ParcelPump(v, csResults, wStations, cCoefficients)
	if err != nil {
		v.Logger.Errorf("Error in Parcel Pumping: %s", err)
	}
	pterm.Success.Println("Successfully Completed Parcel Pumping")

	if err := v.PNirDB.Flush(); err != nil {
		v.Logger.Errorf("Error in flush: %s", err)
	}

	v.Logger.Info("Preforming Dryland Parcel Operations")
	dryParcels, err := parcels.DryLandParcels(v, csResults, wStations, cCoefficients)
	if err != nil {
		v.Logger.Errorf("Error in Dry Land Parcels: %s", err)
	}
	pterm.Success.Println("Successfully Completed Dryland Parcel Ops")

	if err := v.PNirDB.Close(); err != nil { // close doesn't close the db, that must be call explicitly so we can keep using it.
		return err
	}

	// Natural Veg 102
	v.Logger.Info("Preforming Natural Vegetation Operations")
	if err := rchFiles.NaturalVeg(v, wStations, csResults, cCoefficients); err != nil {
		v.Logger.Errorf("Error in Natural Vegatation: %s", err)
		pterm.Warning.Println(fmt.Sprintf("Error in Natural Vegtation: %s", err))
		return err
	}
	pterm.Success.Println("Successfully Completed Natural Vegetation")

	// Dryland 101
	if err := rchFiles.Dryland(v, dryParcels, cCoefficients); err != nil {
		v.Logger.Errorf("Error in Dryland: %s", err)
		return err
	}
	pterm.Success.Println("Successfully Completed Dryland Results")

	if err := v.RchDb.Flush(); err != nil {
		return err
	}

	// Irr Cells
	v.Logger.Info("Preforming Irrigation RCH Operations")
	if err := rchFiles.IrrigationRCH(v, irrParcels, cCoefficients); err != nil {
		v.Logger.Errorf("Error in Creating Irrigation RCH %s", err)
		return err
	}
	pterm.Success.Println("Successfully Completed Irrigated Results")
	if err := v.RchDb.Close(); err != nil {
		return err
	}

	// write out WEL File to db
	if err := wells.WriteWELResults(v, &irrParcels); err != nil {
		return err
	}
	pterm.Success.Println("Successfully Completed WEL Results")

	resultsDB, _ := database.ResultsWelDB(v.SlDb)
	// run steady State Wells
	if err := wells.SteadyStateWells(v, resultsDB); err != nil {
		return err
	}
	pterm.Success.Println("Successfully Completed SS Results")

	// run external wells
	if err := wells.CreateExternalWells(v, resultsDB); err != nil {
		return err
	}
	pterm.Success.Println("Successfully Completed External Wells")

	// run external recharge
	if err := rchFiles.CreateExternalRecharge(v); err != nil {
		return err
	}
	pterm.Success.Println("Successfully Completed External RCH")

	if err := wells.MunicipalIndWells(v, resultsDB); err != nil {
		return err
	}
	pterm.Success.Println("Successfully Completed MI Well Ops")

	_ = noteDb.Close()
	_ = v.SlDb.Close() // close the db before ending the program
	v.Logger.Infof("Model Runtime: %s", time.Since(timeStart))
	v.Logger.Info("Model Completed Normally")
	pterm.Info.Println("Model Completed Normally, check logs for details of run")
	v.Logger.Close()
	return nil
}
