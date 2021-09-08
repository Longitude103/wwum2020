package wells

import (
	"github.com/Longitude103/wwum2020/database"
	"github.com/schollz/progressbar/v3"
)

func CreateExternalWells(v database.Setup) error {
	v.Logger.Info("Starting External Wells Process")

	// get the external wells from "ext_pumping"
	v.Logger.Info("Getting External Wells Data.")
	extPump, err := database.GetExternalWells(v)
	if err != nil {
		return err
	}

	v.Logger.Info("Setting up results DB")
	welDb, err := database.ResultsWelDB(v.SlDb)
	if err != nil {
		return err
	}

	v.Logger.Info("Saving Data to results DB")
	bar := progressbar.Default(int64(len(extPump)), "records saved")
	for i := 0; i < len(extPump); i++ {
		if err := welDb.Add(database.WelResult{Wellid: i, Node: extPump[i].Node, Dt: extPump[i].Date(),
			FileType: extPump[i].FileType, Result: extPump[i].Pumping}); err != nil {
			return err
		}
		_ = bar.Add(1)
	}
	_ = bar.Close()

	v.Logger.Info("Finished adding External Wells to the results dataset")
	return nil
}
