package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/heath140/gisUtils"
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

type CellIntersect struct {
	Node      int             `db:"node"`
	Soil      int             `db:"soil_code"`
	CellArea  float64         `db:"cell_area"`
	NpIrrArea sql.NullFloat64 `db:"nip_area"`
	NpDryArea sql.NullFloat64 `db:"ndp_area"`
	SpIrrArea sql.NullFloat64 `db:"sip_area"`
	SpDryArea sql.NullFloat64 `db:"sdp_area"`
	PointX    float64         `db:"pointx"`
	PointY    float64         `db:"pointy"`
}

func (m ModelCell) GetXY() (x float64, y float64) {
	return m.PointX, m.PointY
}

type StDistances struct {
	Station  string
	Distance float64
	Weight   float64
}

func GetCells(v Setup) (cells []ModelCell, err error) {

	const query = `select node, st_x(st_transform(st_centroid(geom), 4326)) pointx, 
				st_y(st_transform(st_centroid(geom), 4326)) pointy,
				soil_code, coeff_zone, zone, mtg from public.model_cells;`

	if err = v.PgDb.Select(&cells, query); err != nil {
		return nil, err
	}

	if v.AppDebug {
		return cells[:50], nil
	}

	return
}

func GetCellAreas(v Setup, y int) (cells []CellIntersect, err error) {
	query := fmt.Sprintf(`select m.node, m.soil_code, st_area(geom)/43560 cell_area, st_x(st_transform(st_centroid(geom), 4326)) pointx,
       st_y(st_transform(st_centroid(geom), 4326)) pointy, nip_area, ndp_area, sip_area, sdp_area from model_cells m
           left join (select node, sum(st_area(st_intersection(c.geom, ni.geom))/43560) nip_area from public.model_cells c 
               inner join np.t%d_irr ni on st_intersects(c.geom, ni.geom) group by node) ni on m.node = ni.node
           left join (select node, sum(st_area(st_intersection(c.geom, nd.geom))/43560) ndp_area from public.model_cells c 
               inner join np.t%d_dry nd on st_intersects(c.geom, nd.geom) group by node) nd on m.node = nd.node
           left join (select node, sum(st_area(st_intersection(c.geom, si.geom))/43560) sip_area from public.model_cells c 
               inner join sp.t%d_irr si on st_intersects(c.geom, si.geom) group by node) si on m.node = si.node
           left join (select node, sum(st_area(st_intersection(c.geom, sd.geom))/43560) sdp_area from public.model_cells c 
               inner join sp.t%d_dry sd on st_intersects(c.geom, sd.geom) group by node) sd on m.node = sd.node;`, y, y, y, y)

	if err = v.PgDb.Select(&cells, query); err != nil {
		return nil, err
	}

	if v.AppDebug {
		return cells[6600:6800], nil
	}

	return cells, nil
}

type XyPoints interface {
	GetXY() (x float64, y float64)
}

// Distances is a function that that returns the top three weather stations from the list with the appropriate weighting
// factor. Used to make CSResults Distribution.
func Distances(points XyPoints, wStations []WeatherStation) (dist []StDistances, err error) {
	var lengths []float64
	for _, v := range wStations {
		var stDistance StDistances
		pX, pY := points.GetXY()
		d := gisUtils.Distance(pY, pX, v.PointY, v.PointX)
		lengths = append(lengths, d)
		stDistance.Distance = d
		stDistance.Station = v.Code
		dist = append(dist, stDistance)
	}

	sort.Slice(dist, func(i, j int) bool {
		return dist[i].Distance < dist[j].Distance
	})

	sort.Float64s(lengths)

	var idw []float64
	if len(lengths) >= 3 {
		idw, err = gisUtils.InverseDW(lengths[:3])
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("less then three stations present")
	}

	for i, v := range idw {
		dist[i].Weight = v
	}

	return dist[:3], nil
}

func (c CellIntersect) VegArea() float64 {
	return c.CellArea - returnF64(c.NpIrrArea) - returnF64(c.SpIrrArea) - returnF64(c.NpDryArea) - returnF64(c.SpDryArea)
}

func (c CellIntersect) GetXY() (x float64, y float64) {
	return c.PointX, c.PointY
}

func returnF64(v sql.NullFloat64) float64 {
	if v.Valid {
		return v.Float64
	}
	return 0.0
}
