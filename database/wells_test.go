package database

import (
	"fmt"
	"github.com/joho/godotenv"
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
	if err := v.NewSetup(false, false, myEnv); err != nil {
		t.Error("Could not setup DB connection")
	}

	if err = v.SetYears(1997, 1998); err != nil {
		t.Error("Could not setup years")
	}

	mi, err := GetMIWells(v)
	if err != nil {
		t.Errorf("GetMIWells didn't work %s", err)
	}

	fmt.Println(mi)

	// remove results DB and LOG
}
