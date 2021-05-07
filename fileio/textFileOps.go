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
	Soil        int
	MonthlyData []monthlyValues
	Yr          int
	Crop        int
	Tillage     int
	Irrigation  int
}

type monthlyValues struct {
	Et         float64
	Eff_precip float64
	Nir        float64
	Dp         float64
	Ro         float64
	Precip     float64
}

// LoadTextFiles loads cropsim text files from a location and returns a map of the results to use in further processing.
// This should include all the files and results for each station and each year.
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

// getFileData is an function that breaks down the station data and puts it into a struct to work with.
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
		station.Yr, err = strconv.Atoi(elements[1])
		station.Soil, err = strconv.Atoi(elements[2])
		station.Crop, err = strconv.Atoi(elements[3])
		station.Tillage, err = strconv.Atoi(elements[4])
		station.Irrigation, err = strconv.Atoi(elements[5])

		for i := 6; i < 78; i = i + 6 {
			var mvals monthlyValues
			mvals.Et, err = strconv.ParseFloat(elements[i], 64)
			mvals.Eff_precip, err = strconv.ParseFloat(elements[i+1], 64)
			mvals.Nir, err = strconv.ParseFloat(elements[i+2], 64)
			mvals.Dp, err = strconv.ParseFloat(elements[i+3], 64)
			mvals.Ro, err = strconv.ParseFloat(elements[i+4], 64)
			mvals.Precip, err = strconv.ParseFloat(elements[i+5], 64)

			station.MonthlyData = append(station.MonthlyData, mvals)
		}

		stationData = append(stationData, station)
	}

	return stationData
}
