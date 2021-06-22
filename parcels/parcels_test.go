package parcels

import (
	"database/sql"
	"fmt"
	"github.com/Longitude103/wwum2020/database"
	"math"
	"testing"
)

var p = Parcel{ParcelNo: 1234, AppEff: 0.85, Nir: [12]float64{0, 0, 0, 0, 0.2, 0.4, 0.8, 0.8, 0.5, 0, 0, 0},
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

var p2 = Parcel{ParcelNo: 1234, AppEff: 0.85, Nir: [12]float64{0, 0, 0, 0, 0.2, 0.4, 0.8, 0.8, 0.5, 0, 0, 0},
	DryEt: [12]float64{0, 0, 0, 0, 0.05, 0.1, 0.2, 0.2, 0.1, 0, 0, 0}, Et: [12]float64{0, 0, 0, 0, 1.2, 2.5, 4.5, 4.5, 3, 0, 0, 0},
	SWDel: [12]float64{0, 0, 0, 0, 0, 12, 15, 14, 5, 0, 0, 0}, Ro: [12]float64{0, 0, 0, 0, 0, 0, 0, 0.35, 0.87, 0, 0, 0},
	Dp: [12]float64{0, 0, 0, 0, 0, 0, 1.5, 1.1, 1.3, 0, 0, 0}, CoeffZone: 3, Crop1: sql.NullInt64{Int64: 8, Valid: true},
	Crop1Cov: sql.NullFloat64{Float64: 0.5, Valid: true}, Crop2: sql.NullInt64{Int64: 5, Valid: true},
	Crop2Cov: sql.NullFloat64{Float64: 0.5, Valid: true}, Crop3: sql.NullInt64{Int64: 0, Valid: false},
	Crop3Cov: sql.NullFloat64{Float64: 0, Valid: false}, Crop4: sql.NullInt64{Int64: 0, Valid: false},
	Crop4Cov: sql.NullFloat64{Float64: 0, Valid: false}, Nrd: "np", Yr: 2014,
	CertNum: sql.NullString{String: "3456", Valid: true}, PointX: 40.21, PointY: 103.0,
	Sw: sql.NullBool{Bool: true, Valid: true}, Gw: sql.NullBool{Bool: false, Valid: true}}

// p3 is the groundwater only cell made into a parcel from the TFG Example document
var p3 = Parcel{ParcelNo: 159988, AppEff: 0.65, Nir: [12]float64{0, 0, 0, 0, 0, 0, 4.98, 4.31, 1.65, 0, 0, 0},
	DryEt: [12]float64{0.24, 0.62, 0.39, 1.36, 1.82, 5.13, 4.55, 2.66, 1.16, 0.70, 0.66, 0.19},
	Et:    [12]float64{0.27, 0.33, 0.82, 1.36, 1.82, 5.13, 7.77, 7.21, 4.02, 0.44, 0.51, 0.23},
	Pump:  [12]float64{0, 0, 0, 0, 0, 0, 2.34, 2.32, 1.40, 0, 0, 0}, Ro: [12]float64{0, 0, 0, 1.04, 0.73, 1.81, 0, 0.11, 0, 0.03, 0, 0},
	Dp:        [12]float64{0, 0, 0, 0, 0, 0.39, 0, 0, 0, 0.01, 0, 0},
	CoeffZone: 2, SoilCode: 622, Area: 40.0, IrrType: sql.NullString{String: "FLOOD", Valid: true},
	Crop1: sql.NullInt64{Int64: 8, Valid: true}, Crop1Cov: sql.NullFloat64{Float64: 0.5, Valid: true},
	Crop2:    sql.NullInt64{Int64: 5, Valid: true},
	Crop2Cov: sql.NullFloat64{Float64: 0.5, Valid: true}, Crop3: sql.NullInt64{Int64: 0, Valid: false},
	Crop3Cov: sql.NullFloat64{Float64: 0, Valid: false}, Crop4: sql.NullInt64{Int64: 0, Valid: false},
	Crop4Cov: sql.NullFloat64{Float64: 0, Valid: false}, Nrd: "np", Yr: 2014,
	CertNum: sql.NullString{String: "3456", Valid: true}, PointX: 41.4, PointY: 102.5,
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

func TestParcel_setGirFact(t *testing.T) {
	eff := 0.85

	gir, fsl := setGirFact(eff)
	if gir != 1/0.95 || fsl != 0.02 {
		t.Errorf("Error in setGirFact eff of %g, gir got: %g, should be 1.05263..; fsl got: %g, should be 0.02", eff, gir, fsl)
	}

	eff = 0.65
	gir, fsl = setGirFact(eff)

	if gir != 1 || fsl != 0.05 {
		t.Errorf("Error in setGirFact with eff of %g, gir got: %g, should be 1.3333..; fsl got: %g, should be 0.05", eff, gir, fsl)
	}
}

func TestParcel_sumAnnual(t *testing.T) {
	data := [12]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	result := sumAnnual(data)

	if result != 78 {
		t.Errorf("should have gotten 78 but got %g instead", result)
	}
}

func TestParcel_setAppWat(t *testing.T) {
	sw := [12]float64{50, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0}
	gw := [12]float64{50, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0}
	fsl := 0.02

	appWat, sL, pslIrr := setAppWat(sw, gw, fsl)

	if appWat[0] != 100 || appWat[4] != 100 || appWat[5] != 100 {
		t.Errorf("All appWat values should be 100 but got 0: %g, 4:%g, 5:%g", appWat[0], appWat[4], appWat[5])
	}

	if sL[0] != 2 || sL[4] != 2 || sL[5] != 2 {
		t.Errorf("All sL values should be 2 but got 0: %g, 4:%g, 5:%g", sL[0], sL[4], sL[5])
	}

	if pslIrr[0] != 98 || pslIrr[4] != 98 || pslIrr[5] != 98 {
		t.Errorf("All pslIrr values should be 98 but got 0: %g, 4:%g, 5:%g", pslIrr[0], pslIrr[4], pslIrr[5])
	}
}

func TestParcel_setRoDpWt(t *testing.T) {
	ro := [12]float64{10, 75, 5, 100, 0, 0, 0, 0, 0, 0, 0, 0}
	dp := [12]float64{30, 25, 100, 5, 0, 0, 0, 0, 0, 0, 0, 0}

	result := setRoDpWt(ro, dp)

	if result[0] != 0.25 || result[1] != 0.75 || result[2] != 0.2 || result[3] != 0.8 || result[4] != 0.5 {
		t.Errorf("wieights shoudl be 0.25, 0.75, 0.2, 0.8, 0.5... and got %v", result)
	}
}

func TestParcel_setInitialRoDp(t *testing.T) {
	nir := [12]float64{0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	appWat := [12]float64{0, 5, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	psl := [12]float64{1, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0}
	roDpWt := [12]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}

	ro, dp := setInitialRoDp(nir, appWat, psl, roDpWt)

	if ro[0] != 0 || ro[1] != 2 || ro[2] != 0 || ro[3] != 0 {
		t.Errorf("incorrect initial values for RO, should be 0, 2, 0, 0, 0... and got %v", ro)
	}

	if dp[0] != 0 || dp[1] != 2 || dp[2] != 0 || dp[3] != 0 {
		t.Errorf("incorrect initial values for DP, should be 0, 2, 0, 0, 0... and got %v", dp)
	}
}

func TestParcel_setPreGain(t *testing.T) {
	et := [12]float64{12, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	dryEt := [12]float64{6, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	appWat := [12]float64{3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0}
	psl := [12]float64{2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0}

	gainApWat, gainPsl, gainIrrEt, gainDryEt := setPreGain(et, dryEt, appWat, psl)

	fmt.Println(gainApWat, gainPsl, gainIrrEt, gainDryEt)
	if gainApWat != 3 || gainPsl != 2 || gainIrrEt != 12 || gainDryEt != 6 {
		t.Errorf("error in setPreGain, gainApWat: %g, expecting 3; gainPsl: %g, expecting 2; gainIrrEt: %g, "+
			"expecting 12; gainDryEt: %g, expecting 6", gainApWat, gainPsl, gainIrrEt, gainDryEt)
	}
}

func TestParcel_roundTo(t *testing.T) {
	pi := math.Pi

	if roundTo(pi, 3) != 3.142 {
		t.Errorf("round function not working: got %g, expected 3.142", roundTo(pi, 3))
	}

}

func TestParcel_setEtGain(t *testing.T) {
	appWat := 6.07
	psl := 5.76
	cir := 10.22
	gir := 14.59
	eff := 0.65
	irrEt := 19.00
	dryEt := 8.78

	gain := setEtGain(cir, psl, gir, appWat, eff, irrEt, dryEt)

	if roundTo(gain, 2) != 3.95 {
		t.Errorf("error in gain function, got: %g, expected 3.95", roundTo(gain, 2))
	}
}

func TestParcel_distEtGain(t *testing.T) {
	etDry := [12]float64{0.24, 0.62, 0.39, 1.36, 1.82, 5.13, 4.55, 2.66, 1.16, 0.7, 0.66, 0.19}
	etIrr := [12]float64{0.27, 0.33, 0.82, 1.36, 1.82, 5.13, 7.77, 7.21, 4.02, 0.44, 0.51, 0.23}
	psl := [12]float64{0, 0, 0, 0, 0, 0, 2.22, 2.21, 1.33, 0, 0, 0}
	gain := 3.94

	dist := distEtGain(gain, psl, etIrr, etDry)

	fmt.Println(dist)

}

func TestParcel_setEtBase(t *testing.T) {
	etDry := [12]float64{0.24, 0.62, 0.39, 1.36, 1.82, 5.13, 4.55, 2.66, 1.16, 0.7, 0.66, 0.19}
	etIrr := [12]float64{0.27, 0.33, 0.82, 1.36, 1.82, 5.13, 7.77, 7.21, 4.02, 0.44, 0.51, 0.23}
	psl := [12]float64{0, 0, 0, 0, 0, 0, 2.22, 2.21, 1.33, 0, 0, 0}

	base := setEtBase(psl, etIrr, etDry)

	if base[5] != 5.13 || base[6] != 4.55 {
		t.Errorf("base calculated incorrect: June base: %g, expecting 5.13; July base: %g, expecting 4.55", base[5], base[6])
	}
}

func TestParcel_setET(t *testing.T) {
	etBase := [12]float64{0.27, 0.33, 0.82, 1.36, 1.82, 5.13, 4.55, 2.66, 1.16, 0.44, 0.51, 0.23}
	etGain := [12]float64{0, 0, 0, 0, 0, 0, 1.19, 1.69, 1.06, 0, 0, 0}

	et := setET(etBase, etGain)

	if et[5] != 5.13 || et[6] != 5.74 {
		t.Errorf("et calculated incorrect: June ET: %g, expecting 5.13; July ET: %g, expecting 5.74", et[5], et[6])
	}
}

func TestParcel_setDeltaET(t *testing.T) {
	etIrr := [12]float64{0.26, 0.31, 0.78, 1.29, 1.73, 4.87, 5.46, 4.13, 2.11, 0.42, 0.48, 0.22}
	factor := 0.95

	delta := setDeltaET(etIrr, factor)

	if roundTo(delta[0], 3) != 0.013 || roundTo(delta[5], 3) != 0.244 {
		t.Errorf("setDelta calculated incorrect: Jan delta: %g, expected 0.013; June delat: %g, expected 0.244", roundTo(delta[0], 3), roundTo(delta[5], 3))
	}
}
