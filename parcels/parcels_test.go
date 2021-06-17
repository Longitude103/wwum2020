package parcels

import (
	"database/sql"
	"fmt"
	"github.com/Longitude103/wwum2020/database"
	"testing"
)

var p = Parcel{ParcelNo: 1234, AppEff: 0.8, Nir: [12]float64{0, 0, 0, 0, 0.2, 0.4, 0.8, 0.8, 0.5, 0, 0, 0},
	DryEt: [12]float64{0, 0, 0, 0, 0.05, 0.1, 0.2, 0.2, 0.1, 0, 0, 0}, Et: [12]float64{0, 0, 0, 0, 1.2, 2.5, 4.5, 4.5, 3, 0, 0, 0},
	Pump: [12]float64{0, 0, 0, 0, 0, 12.3, 21.4, 18.9, 0, 0, 0, 0}, Ro: [12]float64{0, 0, 0, 0, 0, 0, 0, 0.35, 0.87, 0, 0, 0},
	Dp: [12]float64{0, 0, 0, 0, 0, 0, 1.5, 1.1, 1.3, 0, 0, 0}, SWDel: [12]float64{0, 0, 0, 0, 1, 13, 15, 14, 5, 0, 0, 0},
	CoeffZone: 3, Crop1: sql.NullInt64{Int64: 8, Valid: true},
	Crop1Cov: sql.NullFloat64{Float64: 0.5, Valid: true}, Crop2: sql.NullInt64{Int64: 5, Valid: true},
	Crop2Cov: sql.NullFloat64{Float64: 0.5, Valid: true}, Crop3: sql.NullInt64{Int64: 0, Valid: false},
	Crop3Cov: sql.NullFloat64{Float64: 0, Valid: false}, Crop4: sql.NullInt64{Int64: 0, Valid: false},
	Crop4Cov: sql.NullFloat64{Float64: 0, Valid: false}, Nrd: "np", Yr: 2014,
	CertNum: sql.NullString{String: "3456", Valid: true}, PointX: 40.21, PointY: 103.0,
	Sw: sql.NullBool{Bool: true, Valid: true}, Gw: sql.NullBool{Bool: true, Valid: true}}

var p2 = Parcel{ParcelNo: 1234, AppEff: 0.8, Nir: [12]float64{0, 0, 0, 0, 0.2, 0.4, 0.8, 0.8, 0.5, 0, 0, 0},
	DryEt: [12]float64{0, 0, 0, 0, 0.05, 0.1, 0.2, 0.2, 0.1, 0, 0, 0}, Et: [12]float64{0, 0, 0, 0, 1.2, 2.5, 4.5, 4.5, 3, 0, 0, 0},
	SWDel: [12]float64{0, 0, 0, 0, 0, 12, 15, 14, 5, 0, 0, 0}, Ro: [12]float64{0, 0, 0, 0, 0, 0, 0, 0.35, 0.87, 0, 0, 0},
	Dp: [12]float64{0, 0, 0, 0, 0, 0, 1.5, 1.1, 1.3, 0, 0, 0}, CoeffZone: 3, Crop1: sql.NullInt64{Int64: 8, Valid: true},
	Crop1Cov: sql.NullFloat64{Float64: 0.5, Valid: true}, Crop2: sql.NullInt64{Int64: 5, Valid: true},
	Crop2Cov: sql.NullFloat64{Float64: 0.5, Valid: true}, Crop3: sql.NullInt64{Int64: 0, Valid: false},
	Crop3Cov: sql.NullFloat64{Float64: 0, Valid: false}, Crop4: sql.NullInt64{Int64: 0, Valid: false},
	Crop4Cov: sql.NullFloat64{Float64: 0, Valid: false}, Nrd: "np", Yr: 2014,
	CertNum: sql.NullString{String: "3456", Valid: true}, PointX: 40.21, PointY: 103.0,
	Sw: sql.NullBool{Bool: true, Valid: true}, Gw: sql.NullBool{Bool: false, Valid: true}}

var p3 = Parcel{ParcelNo: 1234, AppEff: 0.8, Nir: [12]float64{0, 0, 0, 0, 0.2, 0.4, 0.8, 0.8, 0.5, 0, 0, 0},
	DryEt: [12]float64{0, 0, 0, 0, 0.05, 0.1, 0.2, 0.2, 0.1, 0, 0, 0}, Et: [12]float64{0, 0, 0, 0, 1.2, 2.5, 4.5, 4.5, 3, 0, 0, 0},
	Pump: [12]float64{0, 0, 0, 0, 0, 12.3, 21.4, 18.9, 0, 0, 0, 0}, Ro: [12]float64{0, 0, 0, 0, 0, 0, 0, 0.35, 0.87, 0, 0, 0},
	Dp:        [12]float64{0, 0, 0, 0, 0, 0, 1.5, 1.1, 1.3, 0, 0, 0},
	CoeffZone: 3, Crop1: sql.NullInt64{Int64: 8, Valid: true}, Crop1Cov: sql.NullFloat64{Float64: 0.5, Valid: true},
	Crop2:    sql.NullInt64{Int64: 5, Valid: true},
	Crop2Cov: sql.NullFloat64{Float64: 0.5, Valid: true}, Crop3: sql.NullInt64{Int64: 0, Valid: false},
	Crop3Cov: sql.NullFloat64{Float64: 0, Valid: false}, Crop4: sql.NullInt64{Int64: 0, Valid: false},
	Crop4Cov: sql.NullFloat64{Float64: 0, Valid: false}, Nrd: "np", Yr: 2014,
	CertNum: sql.NullString{String: "3456", Valid: true}, PointX: 40.21, PointY: 103.0,
	Sw: sql.NullBool{Bool: false, Valid: true}, Gw: sql.NullBool{Bool: true, Valid: true}}

func TestParcel_String(t *testing.T) {
	if p.String() != "Parcel No: 1234, NRD: np, Year: 2014" {
		t.Error("string doesn't produce correct result")
	}
}

func TestParcel_WaterBalanceWWSP(t *testing.T) {
	cCrops := []database.CoeffCrop{{Zone: 3, Crop: 8, NirAdjFactor: 0.95}, {Zone: 3, Crop: 5, NirAdjFactor: 0.95}}

	fmt.Println("RO is:", p.Ro)
	fmt.Println("Dp is:", p.Dp)

	var err error
	err = p.waterBalanceWSPP(cCrops)
	if err != nil {
		t.Errorf("Error in WSPP Method of %s", err)
	}

	fmt.Println("Post Process RO is:", p.Ro)
	fmt.Println("Post Process Dp is:", p.Dp)
}

func TestParcel_WaterBalanceWWSP_SWOnly(t *testing.T) {
	cCrops := []database.CoeffCrop{{Zone: 3, Crop: 8, NirAdjFactor: 0.95}, {Zone: 3, Crop: 5, NirAdjFactor: 0.95}}

	fmt.Println("RO is:", p2.Ro)
	fmt.Println("Dp is:", p2.Dp)

	var err error
	err = p2.waterBalanceWSPP(cCrops)
	if err != nil {
		t.Errorf("Error in WSPP Method of %s", err)
	}

	fmt.Println("Post Process RO is:", p2.Ro)
	fmt.Println("Post Process Dp is:", p2.Dp)
}

func TestParcel_WaterBalanceWWSP_GWOnly(t *testing.T) {
	cCrops := []database.CoeffCrop{{Zone: 3, Crop: 8, NirAdjFactor: 0.95}, {Zone: 3, Crop: 5, NirAdjFactor: 0.95}}

	fmt.Println("RO is:", p3.Ro)
	fmt.Println("Dp is:", p3.Dp)

	var err error
	err = p3.waterBalanceWSPP(cCrops)
	if err != nil {
		t.Errorf("Error in WSPP Method of %s", err)
	}

	fmt.Println("Post Process RO is:", p3.Ro)
	fmt.Println("Post Process Dp is:", p3.Dp)
}

func TestFilterParcelByCert(t *testing.T) {
	sliceP := []Parcel{p}
	fp := filterParcelByCert(&sliceP, "3456")

	if fp[0].ParcelNo != 1234 {
		t.Error("is not returning correct parcel")
	}
}

func TestFilterParcelByCertNoneFound(t *testing.T) {
	sliceP := []Parcel{p}
	fp := filterParcelByCert(&sliceP, "6789")

	if fp != nil {
		t.Error("parcel returned when none should be")
	}
}

func TestParcel_GetXY(t *testing.T) {
	x, y := p.GetXY()

	if x != 40.21 || y != 103.0 {
		t.Error("not returning the correct X, Y")
	}
}
