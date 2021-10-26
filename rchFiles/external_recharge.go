package rchFiles

import (
	"github.com/Longitude103/wwum2020/database"
	"github.com/pterm/pterm"
)

// CreateExternalRecharge is a function that creates the external recharge values in the results' database by getting
// the data from postgres DB and saving those values into the results DB.
func CreateExternalRecharge(v database.Setup) error {
	v.Logger.Info("Starting External Recharge function")

	v.Logger.Info("Getting information from database")
	spin, _ := pterm.DefaultSpinner.Start("Getting External Recharge and results DB")
	eRch, err := database.GetExtRecharge(v)
	if err != nil {
		return err
	}

	rchDB, err := database.ResultsRchDB(v.SlDb)
	if err != nil {
		return err
	}
	spin.Success()

	v.Logger.Info("Saving the recharge information to results database")
	pBar, _ := pterm.DefaultProgressbar.WithTotal(len(eRch)).WithTitle("Well Results Parcels").WithRemoveWhenDone(true).Start()
	for i := 0; i < pBar.Total; i++ {
		pBar.Increment()
		if err := rchDB.Add(database.RchResult{Node: eRch[i].Node, Size: eRch[i].Size, Dt: eRch[i].Date(),
			FileType: eRch[i].FileType, Result: eRch[i].Rch}); err != nil {
			return err
		}
	}

	v.Logger.Info("External Recharge function completed.")
	return nil
}
