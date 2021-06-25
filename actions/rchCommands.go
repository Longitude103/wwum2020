package actions

import (
	"fmt"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/fileio"
	"github.com/Longitude103/wwum2020/parcels"
	"github.com/Longitude103/wwum2020/rchFiles"
	//"wwum2020/rchFiles"
)

func RunModel(debug bool, CSDir *string, sY int, eY int, eF bool, myEnv map[string]string) error {
	v := database.Setup{}
	if err := v.NewSetup(debug, eF, myEnv); err != nil {
		return err
	}
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

	// TODO: Write Dryland Parcel RCH Values to Results DB
	_ = dryParcels

	if err := v.PNirDB.Close(); err != nil { // close doesn't close the db, that must be call explicitly so we can keep using it.
		return err
	}

	//// load up data with cell acres
	//cells, err := database.GetCells(v)
	//if err != nil {
	//	v.Logger.Errorf("Error getting cells from DB: %s", err)
	//	return err
	//}

	// Natural Veg 102
	v.Logger.Info("Preforming Natural Vegetation Operations")
	if err := rchFiles.NaturalVeg(v, wStations, csResults, cCoefficients); err != nil {
		v.Logger.Errorf("Error in Natural Vegatation: %s", err)
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

	//fmt.Println("First Irrigated Cell:")
	//fmt.Println(irrCells[0])
	_ = v.SlDb.Close() // close the db before ending the program
	return nil
	// Dry Cells

}
