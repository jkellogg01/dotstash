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

func DirExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return false, err
	} else if err != nil {
		return false, nil
	}
	return stat.IsDir(), nil
}
