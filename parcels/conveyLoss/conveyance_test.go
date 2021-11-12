package conveyLoss

import "testing"

func Test_Conveyance(t *testing.T) {
	v := dbConnection()

	err := Conveyance(v)
	if err != nil {
		t.Errorf("Error in Conveyance: %s", err)
	}
}
