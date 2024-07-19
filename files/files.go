package files

import (
	"errors"
	"io/fs"
	"os"
	"strings"
)

func GotoFigureRoot() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	pathSlice := strings.Split(homeDir, string(os.PathSeparator))
	pathSlice = append(pathSlice, ".figure")
	path := strings.Join(pathSlice, string(os.PathSeparator))
	err = os.Mkdir(path, fs.ModeDir|fs.ModePerm)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}
	err = os.Chdir(path)
	if err != nil {
		return err
	}
	return nil
}
