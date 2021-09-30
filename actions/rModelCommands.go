package actions

import (
	"fmt"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/fileio"
	"github.com/Longitude103/wwum2020/parcels"
	"github.com/Longitude103/wwum2020/rchFiles"
	"github.com/Longitude103/wwum2020/wells"
	"time"
	//"wwum2020/rchFiles"
)

func RunModel(debug bool, CSDir *string, sY int, eY int, eF bool, myEnv map[string]string) error {
	timeStart := time.Now()

	v := database.Setup{}
	if err := v.NewSetup(debug, eF, myEnv); err != nil {
		return err
	}
	v.Logger.Infof("Model Run Started at: %s", timeStart.Format(time.UnixDate))
	if err := v.SetYears(sY, eY); err != nil {
		v.Logger.Errorf("Error Setting Years Error: %s", err)
		return err
	}

	csResults, err := fileio.LoadTextFiles(*CSDir, v.Logger)
	if err != nil {
		fmt.Println("Error in Loading Text Files, check log file")
		return err
	}

	v.Logger.Info("Getting Weather Stations")
	wStations, err := database.GetWeatherStations(v.PgDb)
	if err != nil {
		v.Logger.Errorf("Error Getting Weather stations: %s", err)
	}

	v.Logger.Info("Getting Coefficients of Crops")
	cCoefficients, err := database.GetCoeffCrops(v)
	if err != nil {
		return err
	}

	// parcel pumping
	v.Logger.Info("Preforming Parcel Pumping")
	irrParcels, err := parcels.ParcelPump(v, csResults, wStations, cCoefficients)
	if err != nil {
		v.Logger.Errorf("Error in Parcel Pumping: %s", err)
	}

	if err := v.PNirDB.Flush(); err != nil {
		v.Logger.Errorf("Error in flush: %s", err)
	}

	v.Logger.Info("Preforming Dryland Parcel Operations")
	dryParcels, err := parcels.DryLandParcels(v, csResults, wStations, cCoefficients)
	if err != nil {
		v.Logger.Errorf("Error in Dry Land Parcels: %s", err)
	}

	if err := v.PNirDB.Close(); err != nil { // close doesn't close the db, that must be call explicitly so we can keep using it.
		return err
	}

	// Natural Veg 102
	v.Logger.Info("Preforming Natural Vegetation Operations")
	if err := rchFiles.NaturalVeg(v, wStations, csResults, cCoefficients); err != nil {
		v.Logger.Errorf("Error in Natural Vegatation: %s", err)
		return err
	}

	// Dryland 101
	if err := rchFiles.Dryland(v, dryParcels); err != nil {
		v.Logger.Errorf("Error in Dryland: %s", err)
		return err
	}

	if err := v.RchDb.Flush(); err != nil {
		return err
	}

	// Irr Cells
	v.Logger.Info("Preforming Irrigation RCH Operations")
	if err := rchFiles.IrrigationRCH(v, irrParcels); err != nil {
		v.Logger.Errorf("Error in Creating Irrigation RCH %s", err)
		return err
	}
	if err := v.RchDb.Close(); err != nil {
		return err
	}

	// write out WEL File to db
	if err := wells.WriteWELResults(v, &irrParcels); err != nil {
		return err
	}

	// run steady State Wells
	if err := wells.SteadyStateWells(v); err != nil {
		return err
	}

	// run external wells
	if err := wells.CreateExternalWells(v); err != nil {
		return err
	}

	// run external recharge
	if err := rchFiles.CreateExternalRecharge(v); err != nil {
		return err
	}

	// TODO: Add Municipal and Industrial Wells

	_ = v.SlDb.Close() // close the db before ending the program
	v.Logger.Infof("Model Run took: %s", time.Now().Sub(timeStart))
	v.Logger.Info("Model Completed without Error")
	fmt.Println("Model Completed without Error.")
	return nil

}
