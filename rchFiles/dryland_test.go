package rchFiles

import (
	"github.com/Longitude103/wwum2020/database"
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
