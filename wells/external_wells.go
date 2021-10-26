package wells

import (
	"github.com/Longitude103/wwum2020/database"
	"github.com/pterm/pterm"
)

func CreateExternalWells(v database.Setup) error {
	v.Logger.Info("Starting External Wells Process")

	spin, _ := pterm.DefaultSpinner.Start("Getting External Wells and results DB")
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
	spin.Success()

	v.Logger.Info("Saving Data to results DB")
	pBar, _ := pterm.DefaultProgressbar.WithTotal(len(extPump)).WithTitle("Well Results Parcels").WithRemoveWhenDone(true).Start()
	for i := 0; i < pBar.Total; i++ {
		pBar.Increment()
		if err := welDb.Add(database.WelResult{Wellid: i, Node: extPump[i].Node, Dt: extPump[i].Date(),
			FileType: extPump[i].FileType, Result: extPump[i].Pmp()}); err != nil {
			return err
		}
	}

	v.Logger.Info("Finished adding External Wells to the results dataset")
	return nil
}
