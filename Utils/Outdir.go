package Utils

import (
	"encoding/json"
	"fmt"
	"io"
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

func (t *tm) trimmed() string {
	return t.Dt[:len(t.Dt)-13]
}

// MakeOutputDir makes the output directory, sends the full path and the sub dir value for making files as return values for the function as well as
// any errors that might have occurred since this hits a world time API.
func MakeOutputDir(folderTitle string) (string, string, error) {
	resp, err := http.Get("https://worldtimeapi.org/api/timezone/America/Denver")
	if err != nil {
		return "", "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var responseTime tm
	if err := json.Unmarshal(body, &responseTime); err != nil {
		fmt.Print("Error getting time for directory")
	}
	responseTime.replaceColon()

	wd, _ := os.Getwd()
	rsT := folderTitle + responseTime.trimmed()
	subPath := filepath.Join("OutputFiles", rsT)
	path := filepath.Join(wd, subPath)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", "", err
	}

	return path, responseTime.trimmed(), nil
}
