package conveyLoss

import "testing"

func Test_Conveyance(t *testing.T) {
	v := dbConnection()

	v.SYear = 1953
	v.EYear = 1953

	err := Conveyance(v)
	if err != nil {
		t.Errorf("Error in Conveyance: %s", err)
	}
}
