package parcels

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/fileio"
	"github.com/joho/godotenv"
)

var u1 = Usage{Yr: 2014, Nrd: "np", CertNum: "3456", UseAF: 100.0}
var u2 = Usage{Yr: 2014, Nrd: "np", CertNum: "3459", UseAF: 240.0}
var u3 = Usage{Yr: 2014, Nrd: "np", CertNum: "3457", UseAF: 100.0}
var u4 = Usage{Yr: 2014, Nrd: "np", CertNum: "3458", UseAF: 200.0}

var testUsageSlice = []Usage{u1, u2, u3, u4}

func dbConnection() *database.Setup {
	var myEnv map[string]string
	myEnv, err := godotenv.Read("../.env")
	if err != nil {
		fmt.Println("Cannot load Env Variables:", err)
		os.Exit(1)
	}

	var v *database.Setup
	v, err = database.NewSetup(myEnv, database.WithLogger(), database.WithNoSQLite(), database.WithDebug())
	if err != nil {
		fmt.Printf("Error in NewSetup: %s", err)
	}

	if err = v.SetYears(1997, 1997); err != nil {
		fmt.Println("error setting years")
	}

	return v
}

func Test_distUsage(t *testing.T) {
	for i := 0; i < 12; i++ {
		testParcelSlice[0].Pump[i] = 0
	}

	err := distUsage(testUsageSlice, &testParcelSlice)
	if err != nil {
		t.Error("Function returned an error:", err)
	}

	total := 0.0
	for _, f := range testParcelSlice[0].Pump {
		total += f
	}

	if total < 9.89 || total > 9.9 {
		t.Errorf("Total pumping should have been 9.897 but got %f", total)
	}
}

func Test_ParcelPump(t *testing.T) {
	v := dbConnection()

	csResults, _ := fileio.LoadTextFiles("/run/media/heath/T7 Touch/WWUMM2020/CropSim/Run005_WWUM2020/Output", v.Logger)
	wStations, _ := database.GetWeatherStations(v.PgDb)
	cCoefficients, _ := database.GetCoeffCrops(v)

	irrParcels, err := ParcelPump(v, csResults, wStations, cCoefficients)
	if err != nil {
		v.Logger.Errorf("Error in Parcel Pumping: %s", err)
	}

	for i, parcel := range irrParcels {
		v.Logger.Debug(parcel.String())
		v.Logger.Debug(parcel.NIRString())
		v.Logger.Debug(strings.Repeat("-", 100))
		v.Logger.Debug(parcel.SWString())
		v.Logger.Debug(parcel.pumpString())
		v.Logger.Debug(strings.Repeat("-", 100))
		v.Logger.Debug(parcel.RoString())
		v.Logger.Debug(parcel.DpString())
		v.Logger.Debug(strings.Repeat("=", 100))

		if i == 15 {
			break
		}
	}

}
