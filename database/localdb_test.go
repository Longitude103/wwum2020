package database

import (
	"fmt"
	"os"
	"testing"
)

type fakelog struct {
	m string
}

func (l *fakelog) Errorf(template string, args ...interface{}) {
	l.m = fmt.Sprintf(template, args...)
}

func (l *fakelog) Infof(template string, args ...interface{}) {
	l.m = fmt.Sprintf(template, args...)
}

func TestGetSqlite(t *testing.T) {
	l := &fakelog{}

	_, err := GetSqlite(l, "./", "testfile")
	if err != nil {
		t.Fatalf("Error Creating SQLite: %s", err)
	}

	if _, err := os.Stat("./OutputFiles/"); os.IsNotExist(err) {
		t.Error("Didn't create the SQLite Database")
	}

	// clean up and remove the test db
	os.RemoveAll("./OutputFiles/")

}
