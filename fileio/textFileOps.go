package fileio

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type StationResults struct {
	Station     string
	soil        int
	monthlyData []monthlyValues
	yr          int
	crop        int
	tillage     int
	irrigation  int
}

type monthlyValues struct {
	et         float64
	eff_precip float64
	nir        float64
	dp         float64
	ro         float64
	precip     float64
}

func LoadTextFiles(filePath string) map[string][]StationResults {
	fmt.Println("File Path:", filePath)
	fls, err := os.ReadDir(filePath)
	if err != nil {
		fmt.Println("Error", err)
	}

	dataMap := make(map[string][]StationResults)
	for _, v := range fls {
		wStationId := v.Name()
		fmt.Printf("Reading station: %s\n", wStationId[:4])
		path := filepath.Join(filePath, v.Name())
		dataMap[wStationId[:4]] = getFileData(path)
	}

	return dataMap
}

func getFileData(filePath string) []StationResults {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	//fmt.Println("Line 1:", scanner.Text())
	//fmt.Println("Line 1 Split", strings.Fields(scanner.Text()))
	//fmt.Println("Line 1 length", len(strings.Fields(scanner.Text()))) // 78

	var stationData []StationResults
	for scanner.Scan() {
		elements := strings.Fields(scanner.Text())
		station := StationResults{Station: elements[0]}
		station.yr, err = strconv.Atoi(elements[1])
		station.soil, err = strconv.Atoi(elements[2])
		station.crop, err = strconv.Atoi(elements[3])
		station.tillage, err = strconv.Atoi(elements[4])
		station.irrigation, err = strconv.Atoi(elements[5])

		for i := 6; i < 78; i = i + 6 {
			var mvals monthlyValues
			mvals.et, err = strconv.ParseFloat(elements[i], 64)
			mvals.eff_precip, err = strconv.ParseFloat(elements[i+1], 64)
			mvals.nir, err = strconv.ParseFloat(elements[i+2], 64)
			mvals.dp, err = strconv.ParseFloat(elements[i+3], 64)
			mvals.ro, err = strconv.ParseFloat(elements[i+4], 64)
			mvals.precip, err = strconv.ParseFloat(elements[i+5], 64)

			station.monthlyData = append(station.monthlyData, mvals)
		}

		stationData = append(stationData, station)
	}

	return stationData
}
