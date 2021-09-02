package wells

import "github.com/Longitude103/wwum2020/database"

// SteadyStateWells is a function to get the steady state wells from the DB and write them into the results db for the
// model year(s) in the proper nodes and with the proper amounts.
func SteadyStateWells(v database.Setup) error {
	// get wells from DB with their node number
	v.Logger.Info("Starting Steady State Wells Addition")
	ssWells, err := database.GetSSWells(v)
	if err != nil {
		return err
	}

	welDB, err := database.ResultsWelDB(v.SlDb)
	if err != nil {
		return err
	}

	for yr := v.SYear; yr < v.EYear+1; yr++ {
		// write them to the results DB
		for i := 0; i < len(ssWells); i++ {
			if err := welDB.Add(database.WelResult{Wellid: ssWells[i].WellName, Node: ssWells[i].Node, Yr: yr,
				FileType: 209, Result: ssWells[i].MVolume}); err != nil {
				return err
			}
		}
	}

	v.Logger.Info("Added all the steady state wells to results DB")
	return nil
}
