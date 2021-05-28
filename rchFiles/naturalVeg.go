package rchFiles

import (
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
	"github.com/heath140/wwum2020/parcels"
)

// NaturalVeg is a function that calculates the area of each cell the is natural vegetation and applies the dryland pasture
// crop type to that area. It then calculates the RO and DP for that crop at that cell location and saves it out as a
// result value in the RCH file. It does use the Adjustment Factors used in previous models.
func NaturalVeg(v database.Setup, wStations []database.WeatherStation,
	csResults map[string][]fileio.StationResults, cCoefficients []database.CoeffCrop) error {
	v.Logger.Infow("Starting Natural Vegetation Ops.")

	for yr := v.SYear; yr < v.EYear+1; yr++ {
		var cellResults []database.NPastCellStruct
		_ = cellResults
		cells, err := database.GetCellAreas(v, yr)
		if err != nil {
			return err
		}

		for i := 0; i < len(cells); i++ {
			dist, err := database.Distances(cells[i], wStations)
			if err != nil {
				return err
			}

			_, _, aDp, aRo, err := database.FilterCCDryLand(cCoefficients, cells[i].CZone, 12)
			if err != nil {
				return err
			}

			cellResult := database.NPastCellStruct{Node: cells[i].Node, Yr: yr}
			for _, st := range dist {
				var annData fileio.StationResults
				for _, data := range csResults[st.Station] {
					if data.Yr == yr && data.Soil == cells[i].Soil &&
						data.Irrigation == int(parcels.DryLand) && data.Crop == 13 {
						annData = data
						break
					}
				}

				for m := 0; m < 12; m++ {
					cellResult.RO[m] = annData.MonthlyData[m].Ro * st.Weight * cells[i].CellArea / 12 * aRo
					cellResult.DP[m] = annData.MonthlyData[m].Dp * st.Weight * cells[i].CellArea / 12 * aDp
				}

			}

			if err := v.NatVegDB.Add(cellResult); err != nil {
				return err
			}
		}
	}
	return nil
}
