package database

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type WeatherStation struct {
	Code string `json:"code"`
	Cor  coord  `json:"location"`
}

func GetWeatherStations(db *sqlx.DB) []WeatherStation {
	rows, err := db.Query(`SELECT code, st_asgeojson(st_transform(geom, 4326)) as location FROM public.weather_stations;`)
	if err != nil {
		panic(err)
	}

	var wStations []WeatherStation
	for rows.Next() {
		station := WeatherStation{}
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
