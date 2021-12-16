package Utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type tm struct {
	Dt string `json:"datetime"`
}

func (t *tm) replaceColon() {
	t.Dt = strings.Replace(t.Dt, ":", "-", -1)
}

func (t tm) trimed() string {
	return "results" + t.Dt[:len(t.Dt)-13]
}

// MakeOutputDir makes the output directory, sends the full path and the subdir value for making files as return values for the function as well as
// any errors that might have occured since this hits a world time API.
func MakeOutputDir() (string, string, error) {
	resp, err := http.Get("http://worldtimeapi.org/api/timezone/America/Denver")
	if err != nil {
		return "", "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var responseTime tm
	json.Unmarshal(body, &responseTime)
	responseTime.replaceColon()

	wd, _ := os.Getwd()
	rsT := responseTime.trimed()
	subPath := filepath.Join("OutputFiles", rsT)
	path := filepath.Join(wd, subPath)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", "", err
	}

	return path, rsT, nil
}
