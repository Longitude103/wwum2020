package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetMIWells(t *testing.T) {
	var myEnv map[string]string
	myEnv, err := godotenv.Read("../.env")
	if err != nil {
		fmt.Println("Cannot load Env Variables:", err)
		os.Exit(1)
	}

	var v Setup
	if err := v.NewSetup(false, false, myEnv, true); err != nil {
		t.Error("Could not setup DB connection")
	}

	if err = v.SetYears(1997, 1998); err != nil {
		t.Error("Could not setup years")
	}

	mi, err := GetMIWells(v)
	if err != nil {
		t.Errorf("GetMIWells didn't work %s", err)
	}

	var well17, well2273 MIWell
	for _, d := range mi {
		if d.WellId == 17 {
			well17 = d
		}
		if d.WellId == 2273 {
			well2273 = d
		}
	}

	// well 17 should have measured pumping
	if len(well17.Pumping) == 0 {
		t.Error("Well 17 should have pumping with it, but does not")
	}

	// well 2373 should not have any pumping
	if len(well2273.Pumping) > 0 {
		t.Error("Well 2273 shouldn't have pumping but does.")
	}

	// removes log file generated when you run v.NewSetup
	path, _ := os.Getwd()
	files, err := ioutil.ReadDir(path)

	for _, file := range files {
		if file.Name()[len(file.Name())-3:] == "log" {
			fmt.Println("removing file:" + file.Name())
			_ = os.Remove(file.Name())
		}
	}
}
