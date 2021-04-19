package distribution

import (
	"clibasic/color"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "long103-wwum.clmtjoquajav.us-east-2.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "rQ!461k&Rk8J"
	dbname   = "wwum"
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

func Distribution(debug *bool, startYr *int, endYr *int) {
	fmt.Println("Distribution")
	if *debug {
		fmt.Println(color.Red + "Debug Mode" + color.Reset)
	}

	fmt.Printf("Start Year: %d -> End Year %d\n", *startYr, *endYr)

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")

	rows, err := db.Query(`SELECT tfg_cellid as cellid, rw, clm, soil_code, 
       st_asgeojson(st_transform(st_centroid(geom), 4326)) as cent FROM public.act_cells;`)

	if err != nil {
		panic(err)
	}

	var Cells []actCell
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
		Cells = append(Cells, cell)
	}

	wStations := getWeatherStations(db)

	for _, c := range Cells[:5] {
		fmt.Println(c)
	}

	for _, w := range wStations {
		fmt.Println(w)
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
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
