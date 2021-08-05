package wells

import "github.com/Longitude103/wwum2020/database"

func writeWELFile(v database.Setup) error {
	wellParcels, err := database.GetWellParcels(v)
	if err != nil {
		return err
	}

	_ = wellParcels

	return nil
}
