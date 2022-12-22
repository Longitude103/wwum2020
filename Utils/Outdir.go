package Utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// MakeOutputDir makes the output directory, sends the full path and the sub dir value for making files as return values for the function as well as
// any errors that might have occurred since this hits a world time API.
func MakeOutputDir(folderTitle string) (string, string, error) {
	goTime := time.Now()
	timeString := fmt.Sprint(goTime.Year()) + "-" + fmt.Sprint(goTime.Month()) + "-" + fmt.Sprint(goTime.Day())
	ts2 := timeString + "T" + fmt.Sprint(goTime.Hour()) + ":" + fmt.Sprint(goTime.Minute()) + ":" + fmt.Sprint(goTime.Second())
	timeString += "T" + fmt.Sprint(goTime.Hour()) + "-" + fmt.Sprint(goTime.Minute()) + "-" + fmt.Sprint(goTime.Second())

	wd, _ := os.Getwd()
	rsT := folderTitle + timeString
	subPath := filepath.Join("OutputFiles", rsT)
	path := filepath.Join(wd, subPath)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", "", err
	}

	return path, ts2, nil
}
