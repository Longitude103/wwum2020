package wells

import (
	"fmt"
	"testing"

	"github.com/Longitude103/wwum2020/database"
)

func TestCreateExternalWells(t *testing.T) {
	sR := &sqlResults{}

	dbconn := dbConnection()
	dbconn.SetYears(1968, 1968)
	dbconn.OldGrid = true
	dbconn.AppDebug = false

	if err := CreateExternalWells(dbconn, sR); err != nil {
		t.Error("Error in function")
	}

	for _, d := range sR.data {
		e := d.(database.WelResult)
		if e.Node == 119612 {
			fmt.Printf("%+v\n", d)
		}
	}

}
