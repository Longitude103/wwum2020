package rchFiles

import (
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/fileio"
	"github.com/Longitude103/wwum2020/parcels"
	"github.com/schollz/progressbar/v3"
	"time"
)

// NaturalVeg is a function that calculates the area of each cell the is natural vegetation and applies the dryland pasture
// crop type to that area. It then calculates the RO and DP for that crop at that cell location and saves it out as a
// result value in the RCH file. It does use the Adjustment Factors used in previous models.
func NaturalVeg(v database.Setup, wStations []database.WeatherStation,
	csResults map[string][]fileio.StationResults, cCoefficients []database.CoeffCrop) error {
	v.Logger.Infow("Starting Natural Vegetation Ops.")

	nVegBarYears := progressbar.Default(int64(v.EYear-v.SYear), "Years of Natural Veg")
	for yr := v.SYear; yr < v.EYear+1; yr++ {
		_ = nVegBarYears.Add(1)
		cells, err := database.GetCellAreas(v, yr)
		if err != nil {
			return err
		}

		nVegBarCells := progressbar.Default(int64(len(cells)), "Natural Veg Cells")
		for i := 0; i < len(cells); i++ {
			_ = nVegBarCells.Add(1)
			var cellResult database.RchResult
			dist, err := database.Distances(cells[i], wStations)
			if err != nil {
				v.Logger.Errorf("error in distance calculation for cell: %v", cells[i])
				return err
			}

			_, _, aDp, aRo, err := database.FilterCCDryLand(cCoefficients, cells[i].CZone, 13)
			if err != nil {
				v.Logger.Errorf("error in getting FilterCCDryLand Function for cell: %v and crop 13", cells[i].CZone)
				return err
			}

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
					cellResult = database.RchResult{Node: cells[i].Node,
						Dt: time.Date(yr, time.Month(m+1), 1, 0, 0, 0, 0, time.UTC), FileType: 102,
						Result: annData.MonthlyData[m].Ro*st.Weight*cells[i].VegArea()/12*aRo + annData.MonthlyData[m].Dp*st.Weight*cells[i].VegArea()/12*aDp}
				}

			}

			if err := v.RchDb.Add(cellResult); err != nil {
				v.Logger.Errorf("Error Adding Result to RchDB Buffer, Result: %+v", cellResult)
				return err
			}
		}
		_ = nVegBarCells.Close()
	}
	_ = nVegBarYears.Close()

	v.Logger.Info("finished natural vegetation function.")
	return nil
}
