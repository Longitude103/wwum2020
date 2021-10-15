package parcels

import (
	"database/sql"
	"github.com/Longitude103/wwum2020/database"
	"testing"
)

var (
	cCrop1 = database.CoeffCrop{
		Zone:         1,
		Crop:         1,
		DryEtAdj:     0.1,
		IrrEtAdj:     0.2,
		NirAdjFactor: 0.3,
		FslGW:        0.4,
		DryEtToRo:    0.5,
		FslSW:        0.6,
		PerToRch:     0.7,
		DpAdj:        0.8,
		RoAdj:        0.9,
	}
	cCrop2 = database.CoeffCrop{
		Zone:         2,
		Crop:         2,
		DryEtAdj:     1.1,
		IrrEtAdj:     1.2,
		NirAdjFactor: 1.3,
		FslGW:        1.4,
		DryEtToRo:    1.5,
		FslSW:        1.6,
		PerToRch:     1.7,
		DpAdj:        1.8,
		RoAdj:        1.9,
	}
	cCrop3 = database.CoeffCrop{
		Zone:         2,
		Crop:         1,
		DryEtAdj:     2.1,
		IrrEtAdj:     2.2,
		NirAdjFactor: 0.95,
		FslGW:        2.4,
		DryEtToRo:    2.5,
		FslSW:        2.6,
		PerToRch:     2.7,
		DpAdj:        2.8,
		RoAdj:        2.9,
	}

	cCrops = []database.CoeffCrop{cCrop1, cCrop2, cCrop3}
)

func TestPumping_adjFactor(t *testing.T) {
	result, err := adjFactor(cCrops, 1, 1, database.DryET)
	if result != 0.1 {
		t.Errorf("adjFactor error: got DryETAdj %f, expected 0.1, error: %s", result, err)
	}

	result, err = adjFactor(cCrops, 2, 2, database.NirEt)
	if result != 1.3 {
		t.Errorf("adjFactor error: got NirETAdj %f, expected 1.3, error: %s", result, err)
	}
}

func TestPumping_adjustmentFactor(t *testing.T) {
	result, err := adjustmentFactor(&p3, cCrops, database.NirEt)

	if result != 0.95 {
		t.Errorf("adjustmentFactor error: got %f, expected 0.95; error: %s", result, err)
	}

	result, err = adjustmentFactor(&p2, cCrops, database.NirEt)

	if result != 0 || err == nil {
		t.Errorf("adjustmentFactor should return a zero with error")
	}
}

func TestPumping_estimatePumping(t *testing.T) {
	//zero pumping
	p1.Pump = [12]float64{}
	// alter NIR
	p1.Nir = [12]float64{0, 0, 0, 0, 3, 25, 40, 45, 20, 0, 0, 0}
	p1.Subarea = sql.NullString{String: "FA", Valid: true}

	err := p1.estimatePumping(cCrops)
	if err != nil {
		t.Errorf("Should not return error")
	}

	mayPump := roundTo(p1.Pump[4], 2)
	junePump := roundTo(p1.Pump[5], 2)
	if mayPump != 2.58 || junePump != 17.06 {
		t.Errorf("Pumping is not calculated correctly: May got %f, expected 2.58; June got %f, expected 17.06", mayPump, junePump)
	}

	err = p2.estimatePumping(cCrops)
	if err == nil {
		t.Errorf("should return an error with this test parcel")
	}
}
