package wells

import "github.com/Longitude103/wwum2020/database"

func createExternalWells(v database.Setup) error {
	// get the external wells from "ext_pumping"
	extPump, err := database.GetExternalWells(v)
	if err != nil {
		return err
	}

	_ = extPump
	// modify the data

	// save to results DB

	return nil
}
