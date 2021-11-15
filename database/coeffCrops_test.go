package database

import (
	"testing"
)

func TestFilterCCDryLand(t *testing.T) {
	// create cSlice, only need dryETadj
	c1 := CoeffCrop{Crop: 8, Zone: 1, DryEtAdj: 1.1, DryEtToRo: 0.5}
	c2 := CoeffCrop{Crop: 7, Zone: 1, DryEtAdj: 1.2, DryEtToRo: 0.6}
	c3 := CoeffCrop{Crop: 15, Zone: 1, DryEtAdj: 1.3, DryEtToRo: 0.7}
	c4 := CoeffCrop{Crop: 7, Zone: 2, DryEtAdj: 1.4, DryEtToRo: 0.8}
	c5 := CoeffCrop{Crop: 8, Zone: 2, DryEtAdj: 1.5, DryEtToRo: 0.9}

	var cS = []CoeffCrop{c1, c2, c3, c4, c5}

	dEA, _, _, _, _, err := FilterCCDryLand(cS, 1, 8)
	if err != nil {
		t.Error("Error with crop 8, zone 1")
	}

	if dEA != 1.1 {
		t.Errorf("should have returned 1.1, but got %f", dEA)
	}

	dEA, _, _, _, _, err = FilterCCDryLand(cS, 1, 7)
	if err != nil {
		t.Error("Error with crop 7, zone 1")
	}

	if dEA != 1.2 {
		t.Errorf("should have returned 1.2, but got %f", dEA)
	}

	dEA, _, _, _, _, err = FilterCCDryLand(cS, 1, 15)
	if err != nil {
		t.Error("Error with crop 15, zone 1")
	}
	if dEA != 1.3 {
		t.Errorf("should have returned 1.3, but got %f", dEA)
	}

	dEA, _, _, _, _, err = FilterCCDryLand(cS, 2, 7)
	if err != nil {
		t.Error("Error with crop 7, zone 2")
	}
	if dEA != 1.4 {
		t.Errorf("should have returned 1.4, but got %f", dEA)
	}

	dEA, _, _, _, _, err = FilterCCDryLand(cS, 2, 15)
	if err != nil {
		t.Error("Error with crop 15, zone 2")
	}
	if dEA != 1.4 {
		t.Errorf("this is crop 15 (not found) and defaults to crop 7 and should have returned 1.4, but got %f", dEA)
	}

	_, _, _, _, _, err = FilterCCDryLand(cS, 2, 9)
	if err == nil {
		t.Error("Should have produced an error since crop isn't found, but it didn't")
	}

}
