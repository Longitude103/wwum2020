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

func MakeOutputDir(fileName string) (string, error) {
	resp, err := http.Get("http://worldtimeapi.org/api/timezone/America/Denver")
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseTime tm
	json.Unmarshal(body, &responseTime)
	responseTime.replaceColon()

	wd, _ := os.Getwd()
	subPath := filepath.Join("OutputFiles", responseTime.trimed())
	path := filepath.Join(wd, subPath)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}

	return path, nil
}
