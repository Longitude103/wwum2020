package distribution

import (
	"clibasic/color"
	"fmt"
	_ "github.com/lib/pq"
	"wwum2020/database"
	"wwum2020/fileio"
)

func Distribution(debug *bool, startYr *int, endYr *int, CSDir string) {
	fmt.Println("Distribution")
	if *debug {
		fmt.Println(color.Red + "Debug Mode" + color.Reset)
	}

	stationData := fileio.LoadTextFiles(CSDir)
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
