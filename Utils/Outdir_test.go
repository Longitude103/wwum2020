package Utils

import (
	"strings"
	"testing"
)

func TestMakeOutputDir(t *testing.T) {
	o, fn, err := MakeOutputDir("results")
	if err != nil {
		t.Error("error creating directory name: ", err)
	}

	if !strings.Contains(o, "/results") {
		t.Error("should be a path with /results in it but got: ", o)
	}

	if !strings.Contains(fn, "results") {
		t.Error("should be a file name with results in it but got: ", fn)
	}
}
