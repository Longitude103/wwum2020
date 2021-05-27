package rchFiles

import (
	"fmt"
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
	"github.com/heath140/wwum2020/parcels"
)

func NaturalVeg(v database.Setup, wStations []database.WeatherStation, csResults map[string][]fileio.StationResults) error {
	v.Logger.Infow("Starting Natural Vegetation Ops.")
	cCoefficients := database.GetCoeffCrops(v.PgDb)

	for yr := v.SYear; yr < v.EYear+1; yr++ {
		cells, err := database.GetCellAreas(v, yr)
		if err != nil {
			return err
		}

		for i, v := range cells {
			fmt.Printf("Area of Node: %d Natural Veg is: %.2f, soil is %d\n", v.Node, v.VegArea(), v.Soil)

			if i == 10 {
				break
			}
		}

		for i := 0; i < len(cells); i++ {
			dist, err := database.Distances(cells[i], wStations)
			if err != nil {
				return err
			}

			for _, st := range dist {
				var annData []fileio.StationResults
				for _, data := range csResults[st.Station] {
					if data.Yr == yr && data.Soil == cells[i].Soil &&
						data.Irrigation == int(parcels.DryLand) && data.Crop == 13 {
						annData = append(annData, data)
					}
				}

				fmt.Printf("Annual Data is: %v\n", annData)
				fmt.Printf("Station is: %s, and weight is %.2f\n", st.Station, st.Weight)

				// TODO: Calculate RO at each cell for 102 (RO gets pushed into recharge)
				// TODO: Use Adjustment factors for DP and RO from DB.
				// TODO: Calculate DP at each cell for 102 (DP is recharge at the cell)

				_ = cCoefficients
				// TODO: Add them together and send to "Results" sqlite file, do not split them

			}

			return nil

		}

	}

	return nil
}
