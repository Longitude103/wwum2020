package rchFiles

import (
	"fmt"
	"testing"
)

func Test_cellRCH(t *testing.T) {

	// check error
	_, err := cellRCH(2.5, 3.5, 0, 40, 0.9, 0.5)
	if err == nil {
		t.Error("Parcel should produce error")
	}

	r, err := cellRCH(2.5, 3.5, 40, 20, 0.9, 0.5)
	if err != nil {
		t.Errorf("Error running function: %s", err)
	}

	fmt.Println(r)

}
