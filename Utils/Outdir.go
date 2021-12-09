package Utils

import (
	"os"
	"path/filepath"
)

func MakeOutputDir(fileName string) (string, error) {
	wd, _ := os.Getwd()
	subPath := filepath.Join("OutputFiles", fileName[:len(fileName)-7])
	path := filepath.Join(wd, subPath)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}

	return path, nil
}
