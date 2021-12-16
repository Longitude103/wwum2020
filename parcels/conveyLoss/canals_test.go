package conveyLoss

import (
	"fmt"
	"os"
	"testing"

	"github.com/Longitude103/wwum2020/database"
	"github.com/joho/godotenv"
)

func dbConnection() *database.Setup {
	var myEnv map[string]string
	myEnv, err := godotenv.Read("../../.env")
	if err != nil {
		fmt.Println("Cannot load Env Variables:", err)
		os.Exit(1)
	}

	var v *database.Setup
	v, err = database.NewSetup(myEnv, database.WithNoSQLite(), database.WithDebug())
	if err != nil {
		fmt.Println("Error in NewSetup")
	}

	if err = v.SetYears(1997, 1997); err != nil {
		fmt.Println("error setting years")
	}

	return v
}

func Test_getCanals(t *testing.T) {
	v := dbConnection()
	v.SYear = 1953
	v.EYear = 1953

	c, err := getCanals(v)
	if err != nil {
		t.Errorf("Error getting canals: %s", err)
	}

	for _, canal := range c {
		v.Logger.Debugf("Canal: %+v\n", canal)
	}

	if len(c) == 0 {
		t.Error("Didn't return any canals")
	}

}

func Test_getCanalCells(t *testing.T) {
	v := dbConnection()

	cc, err := getCanalCells(v)
	if err != nil {
		t.Errorf("Error getting canal cells: %s", err)
	}

	for _, cell := range cc {
		if cell.CLinkId == 25 {
			v.Logger.Debugf("Canal 25: %+v", cell)
		}

	}

	if len(cc) == 0 {
		t.Error("Didn't return any canal cells")
	}
}
