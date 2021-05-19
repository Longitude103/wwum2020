package distribution

import (
	"fmt"
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func Distribution(debug *bool, startYr *int, endYr *int, CSDir string, logger *zap.SugaredLogger) {
	fmt.Println("Distribution")
	if *debug {
		fmt.Println("Debug Mode")
	}

	stationData := fileio.LoadTextFiles(CSDir, logger)
	//fmt.Println(stationData["AGAT"])
	_ = stationData

	fmt.Printf("Start Year: %d -> End Year %d\n", *startYr, *endYr)
	db := database.PgConnx()

	cells := database.GetCells(db)
	wStations := database.GetWeatherStations(db) // weather station list

	for _, c := range cells[:5] {
		dist := database.Distances(c, wStations)
		for _, v := range dist {
			fmt.Printf("Cell Address: %d, Distance to station %s is %.0f Meters and weight is %.4f\n",
				c.CellId, v.Station, v.Distance, v.Weight)
		}
	}

}
