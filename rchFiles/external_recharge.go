package rchFiles

import (
	"github.com/Longitude103/wwum2020/database"
	"github.com/schollz/progressbar/v3"
)

// CreateExternalRecharge is a function that creates the external recharge values in the results' database by getting
// the data from postgres DB and saving those values into the results DB.
func CreateExternalRecharge(v database.Setup) error {
	v.Logger.Info("Starting External Recharge function")

	v.Logger.Info("Getting information from database")
	eRch, err := database.GetExtRecharge(v)
	if err != nil {
		return err
	}

	rchDB, err := database.ResultsRchDB(v.SlDb)
	if err != nil {
		return err
	}

	v.Logger.Info("Saving the recharge information to results database")
	bar := progressbar.Default(int64(len(eRch)), "External Recharge records saved")
	for i := 0; i < len(eRch); i++ {
		if err := rchDB.Add(database.RchResult{Node: eRch[i].Node, Dt: eRch[i].Date(), FileType: eRch[i].FileType,
			Result: eRch[i].Rch}); err != nil {
			return err
		}
		_ = bar.Add(1)
	}
	_ = bar.Close()

	v.Logger.Info("External Recharge function completed.")
	return nil
}
