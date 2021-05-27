package database

import (
	"fmt"
	"github.com/heath140/gisUtils"
	"github.com/jmoiron/sqlx"
	"sort"
)

type ModelCell struct {
	Node      int     `db:"node"`
	SoilCode  int     `db:"soil_code"`
	CoeffZone int     `db:"coeff_zone"`
	Zone      int     `db:"zone"`
	Mtg       float64 `db:"mtg"`
	PointX    float64 `db:"pointx"`
	PointY    float64 `db:"pointy"`
}

func (m ModelCell) GetXY() (x float64, y float64) {
	return m.PointX, m.PointY
}

type StDistances struct {
	Station  string
	Distance float64
	Weight   float64
}

func GetCells(db *sqlx.DB) (cells []ModelCell, err error) {

	const query = `select node, st_x(st_transform(st_centroid(geom), 4326)) pointx, 
				st_y(st_transform(st_centroid(geom), 4326)) pointy,
				soil_code, coeff_zone, zone, mtg from public.model_cells;`

	err = db.Select(&cells, query)
	if err != nil {
		return nil, err
	}

	return
}

type XyPoints interface {
	GetXY() (x float64, y float64)
}

// Distances is a function that that returns the top three weather stations from the list with the appropriate weighting
// factor. Used to make CSResults Distribution.
func Distances(points XyPoints, wStations []WeatherStation) []StDistances {
	var dist []StDistances
	var lenghts []float64
	for _, v := range wStations {
		var stDistance StDistances
		pX, pY := points.GetXY()
		d := gisUtils.Distance(pY, pX, v.PointY, v.PointX)
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
