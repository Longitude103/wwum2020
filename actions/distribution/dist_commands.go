package distribution

import (
	"clibasic/color"
	"database/sql"
	"encoding/json"
	"fmt"
	"gisUtils"
	_ "github.com/lib/pq"
	"sort"
	"wwum2020/database"
	"wwum2020/fileio"
)

type actCell struct {
	Rw       int
	Clm      int
	SoilCode int
	CellId   int
	cor      coord
}

type coord struct {
	T           string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type weatherStation struct {
	Code string `json:"code"`
	Cor  coord  `json:"location"`
}

type stDistances struct {
	Station  string
	Distance float64
	weight   float64
}

func Distribution(debug *bool, startYr *int, endYr *int, CSDir string) {
	fmt.Println("Distribution")
	if *debug {
		fmt.Println(color.Red + "Debug Mode" + color.Reset)
	}

	stationData := fileio.LoadTextFiles(CSDir)
	//fmt.Println(stationData["AGAT"])
	_ = stationData

	fmt.Printf("Start Year: %d -> End Year %d\n", *startYr, *endYr)
	db := database.PgConn()

	cells := getCells(db)
	wStations := getWeatherStations(db) // weather station list

	for _, c := range cells[:5] {
		dist := distances(c, wStations)
		for _, v := range dist {
			fmt.Printf("Cell Address: %d, Distance to station %s is %.0f Meters and weight is %.4f\n",
				c.CellId, v.Station, v.Distance, v.weight)
		}
	}

}

func getCells(db *sql.DB) []actCell {
	rows, err := db.Query(`SELECT tfg_cellid as cellid, rw, clm, soil_code, 
       st_asgeojson(st_transform(st_centroid(geom), 4326)) as cent FROM public.act_cells;`)
	if err != nil {
		panic(err)
	}

	var cells []actCell // active cells list
	cell := actCell{}
	for rows.Next() {
		var cellid, rw, clm, soil int
		var c []byte
		var cor coord

		err = rows.Scan(&cellid, &rw, &clm, &soil, &c)
		if err != nil {
			panic(err.Error())
		}

		err := json.Unmarshal(c, &cor)
		if err != nil {
			fmt.Println("error", err)
		}

		cell.Rw = rw
		cell.Clm = clm
		cell.SoilCode = soil
		cell.CellId = cellid
		cell.cor = cor
		cells = append(cells, cell)
	}

	return cells
}

func getWeatherStations(db *sql.DB) []weatherStation {
	rows, err := db.Query(`SELECT code, st_asgeojson(st_transform(geom, 4326)) as location FROM public.weather_stations;`)
	if err != nil {
		panic(err)
	}

	var wStations []weatherStation
	for rows.Next() {
		station := weatherStation{}
		var code string
		var c []byte
		var cor coord

		err = rows.Scan(&code, &c)
		if err != nil {
			fmt.Println("error", err)
		}

		err := json.Unmarshal(c, &cor)
		if err != nil {
			fmt.Println("error", err)
		}

		station.Code = code
		station.Cor = cor
		wStations = append(wStations, station)
	}

	return wStations
}

// distances is a function that
func distances(cell actCell, wStations []weatherStation) []stDistances {
	var dist []stDistances
	var lenghts []float64
	for _, v := range wStations {
		var stDistance stDistances
		d := gisUtils.Distance(cell.cor.Coordinates[1], cell.cor.Coordinates[0], v.Cor.Coordinates[1], v.Cor.Coordinates[0])
		lenghts = append(lenghts, d)
		stDistance.Distance = d
		stDistance.Station = v.Code
		dist = append(dist, stDistance)
	}

	sort.Slice(dist, func(i, j int) bool {
		return dist[i].Distance < dist[j].Distance
	})

	sort.Float64s(lenghts)

	idw, err := gisUtils.InverseDW(lenghts[:3])
	if err != nil {
		fmt.Println("Error", err)
	}

	for i, v := range idw {
		dist[i].weight = v
	}

	return dist[:3]
}
