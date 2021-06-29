package rchFiles

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/parcels"
	"github.com/jmoiron/sqlx"
	"testing"
	"time"
)

var (
	r1 = database.RchResult{Node: 123456, Dt: time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC),
		FileType: 101, Result: 1}

	r2 = database.RchResult{Node: 234567, Dt: time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC),
		FileType: 101, Result: 2}

	r3 = database.RchResult{Node: 123456, Dt: time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC),
		FileType: 101, Result: 2}

	sliceR = []database.RchResult{r1, r2, r3}

	dp1 = parcels.Parcel{ParcelNo: 101, Nrd: "sp", Ro: [12]float64{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}, Dp: [12]float64{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}}
	dp2 = parcels.Parcel{ParcelNo: 102, Nrd: "sp", Ro: [12]float64{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}, Dp: [12]float64{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}}
	dp3 = parcels.Parcel{ParcelNo: 101, Nrd: "np", Ro: [12]float64{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}, Dp: [12]float64{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}}
	dp4 = parcels.Parcel{ParcelNo: 103, Nrd: "sp", Ro: [12]float64{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}, Dp: [12]float64{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}}

	sliceDryParcels = []parcels.Parcel{dp1, dp2, dp3, dp4}
)

func Test_findResult(t *testing.T) {
	found, location := findResult(sliceR, 234567, time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC))
	if found != true || location != 1 {
		t.Error("should have found record but it didn't")
	}

	found, location = findResult(sliceR, 1, time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC))
	if found == true || location != 0 {
		t.Error("found a result that there was none to find")
	}

	found, location = findResult(sliceR, 234567, time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC))
	if found == true || location != 0 {
		t.Error("found a result that there was none to find")
	}
}

func Test_inGrouped(t *testing.T) {
	result := inGrouped(sliceR, 234567)
	if result != true {
		t.Error("Should be false, but got true")
	}

	result = inGrouped(sliceR, 1)
	if result == true {
		t.Error("Should be false, but got true")
	}
}

func Test_groupResults(t *testing.T) {
	result := groupResults(sliceR)

	if result[0].Result != 3 {
		t.Errorf("Grouping not working, expected 3 but got %f", result[0].Result)
	}

	newSliceR := sliceR[:len(sliceR)-1]
	result = groupResults(newSliceR)

	if result[0].Result != 1 {
		t.Errorf("Grouping not working with slice of nothing to group, expected 1 got %f", result[0].Result)
	}
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func Test_dryland(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := sqlmock.NewRows([]string{"node", "c_area", "d_area", "parcel_id", "nrd"}).AddRow(1, 40, 5, 101, "sp").AddRow(2, 40, 5, 102, "sp")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	mockSqliteDB, mockSqlite, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockSqliteDB.Close()
	sqlitexDB := sqlx.NewDb(mockSqliteDB, "sqlmock")
	mockSqlite.ExpectPrepare("INSERT INTO results")
	rchDB, err1 := database.ResultsRchDB(sqlitexDB)
	if err1 != nil {
		t.Fatalf("an error has occured in rchDB: %s", err1)
	}

	mockSqlite.ExpectBegin()
	mockSqlite.ExpectCommit()
	//mockSqlite.ExpectExec("INSERT INTO results").WithArgs(1, AnyTime{}, 101, 0)

	v := database.Setup{SYear: 2014, EYear: 2014, RchDb: rchDB, PgDb: sqlxDB}
	if err2 := dryland(v, sliceDryParcels); err2 != nil {
		t.Errorf("Error in dryland function: %s", err2)
	}

	if err3 := rchDB.Flush(); err3 != nil {
		t.Errorf("error flushing rchdb: %s", err3)
	}

	// we make sure that all expectations were met
	if err4 := mock.ExpectationsWereMet(); err4 != nil {
		t.Errorf("there were unfulfilled expectations: %s", err4)
	}

	if err5 := mockSqlite.ExpectationsWereMet(); err5 != nil {
		t.Errorf("there were unfulfilled expectations: %s", err5)
	}
}
