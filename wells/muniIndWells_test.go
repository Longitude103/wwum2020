package wells

import (
	"fmt"
	"os"
	"testing"

	"github.com/Longitude103/wwum2020/database"
	"github.com/joho/godotenv"
)

func dbConnection() *database.Setup {
	var myEnv map[string]string
	myEnv, err := godotenv.Read("../.env")
	if err != nil {
		fmt.Println("Cannot load Env Variables:", err)
		os.Exit(1)
	}

	var v *database.Setup
	v, err = database.NewSetup(myEnv, database.WithNoSQLite(), database.WithDebug(), database.WithOldGrid())
	if err != nil {
		fmt.Printf("Error in NewSetup: %s", err)
	}

	if err = v.SetYears(1953, 1955); err != nil {
		fmt.Println("error setting years")
	}

	return v
}

type sqlResults struct {
	data []interface{}
}

func (s *sqlResults) Add(value interface{}) error {
	// fmt.Println(value)
	s.data = append(s.data, value)
	return nil
}

func Test_MunicipalIndWells(t *testing.T) {
	sR := &sqlResults{}

	if err := MunicipalIndWells(dbConnection(), sR); err != nil {
		t.Error("Error in function")
	}

	for _, d := range sR.data {
		e := d.(database.WelResult)
		if e.Wellid == 2237 {
			fmt.Printf("%+v\n", d)
		}
	}

	// fmt.Println(sR.data...)
}
