package files

import (
	"errors"
	"io/fs"
	"os"
	"strings"
)

func GetFigurePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	pathSlice := strings.Split(homeDir, string(os.PathSeparator))
	pathSlice = append(pathSlice, ".figure")
	path := strings.Join(pathSlice, string(os.PathSeparator))
	err = os.Mkdir(path, fs.ModeDir|fs.ModePerm)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return "", err
	}
	return path, nil
}
