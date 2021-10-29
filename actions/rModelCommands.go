package actions

import (
	"fmt"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/fileio"
	"github.com/Longitude103/wwum2020/parcels"
	"github.com/Longitude103/wwum2020/rchFiles"
	"github.com/Longitude103/wwum2020/wells"
	"github.com/pterm/pterm"
	"time"
)

func RunModel(debug bool, CSDir *string, mDesc string, sY int, eY int, eF bool, myEnv map[string]string) error {
	timeStart := time.Now()

	pterm.Info.Println("Setting up results database")
	v := database.Setup{}
	if err := v.NewSetup(debug, eF, myEnv, false, mDesc); err != nil {
		return err
	}
	v.Logger.Infof("Model Run Started at: %s", timeStart.Format(time.UnixDate))

	pterm.Info.Printf("Model Description: %s\n", mDesc)
	v.Logger.Infof("Model Description: %s", mDesc)
	if err := v.SetYears(sY, eY); err != nil {
		v.Logger.Errorf("Error Setting Years Error: %s", err)
		return err
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
	if err := rchFiles.Dryland(v, dryParcels); err != nil {
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

	// run steady State Wells
	if err := wells.SteadyStateWells(v); err != nil {
		return err
	}
	pterm.Success.Println("Successfully Completed SS Results")

	// run external wells
	if err := wells.CreateExternalWells(v); err != nil {
		return err
	}
	pterm.Success.Println("Successfully Completed External Wells")

	// run external recharge
	if err := rchFiles.CreateExternalRecharge(v); err != nil {
		return err
	}
	pterm.Success.Println("Successfully Completed External RCH")

	if err := wells.MunicipalIndWells(v); err != nil {
		return err
	}
	pterm.Success.Println("Successfully Completed MI Well Ops")

	_ = v.SlDb.Close() // close the db before ending the program
	v.Logger.Infof("Model Runtime: %s", time.Now().Sub(timeStart))
	v.Logger.Info("Model Completed Normally")
	pterm.Info.Println("Model Completed Normally, check logs for details of run")
	return nil
}
