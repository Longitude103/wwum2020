package database

import (
	"github.com/jmoiron/sqlx"
)

type WeatherStation struct {
	Code   string  `db:"code"`
	PointX float64 `db:"pointx"`
	PointY float64 `db:"pointy"`
}

// GetWeatherStations is a function to return a slice of WeatherStation that includes the weather station text id (code)
// and the X and Y coordinates of it's location.
func GetWeatherStations(db *sqlx.DB) (wStations []WeatherStation, err error) {
	//goland:noinspection ALL
	query := `SELECT code, st_x(st_transform(st_centroid(geom), 4326)) pointx, 
				st_y(st_transform(st_centroid(geom), 4326)) pointy FROM public.weather_stations;`

	err = db.Select(&wStations, query)
	if err != nil {
		return nil, err
	}

	return
}
