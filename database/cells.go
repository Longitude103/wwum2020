package database

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/heath140/gisUtils"
	"github.com/jmoiron/sqlx"
)

type ActCell struct {
	Rw       int
	Clm      int
	SoilCode int
	CellId   int
	Cor      coord
}

type coord struct {
	T           string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type StDistances struct {
	Station  string
	Distance float64
	Weight   float64
}

func GetCells(db *sqlx.DB) []ActCell {
	rows, err := db.Query(`SELECT tfg_cellid as cellid, rw, clm, soil_code, 
       st_asgeojson(st_transform(st_centroid(geom), 4326)) as cent FROM public.act_cells;`)
	if err != nil {
		panic(err)
	}

	var cells []ActCell // active cells list
	cell := ActCell{}
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
		cell.Cor = cor
		cells = append(cells, cell)
	}

	return cells
}

// Distances is a function that
func Distances(cell ActCell, wStations []WeatherStation) []StDistances {
	var dist []StDistances
	var lenghts []float64
	for _, v := range wStations {
		var stDistance StDistances
		d := gisUtils.Distance(cell.Cor.Coordinates[1], cell.Cor.Coordinates[0], v.Cor.Coordinates[1], v.Cor.Coordinates[0])
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
		dist[i].Weight = v
	}

	return dist[:3]
}
