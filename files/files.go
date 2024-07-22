package files

import (
	"errors"
	"io/fs"
	"os"
	"path"
)

func GetFigurePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := path.Join(homeDir, ".figure")
	err = os.Mkdir(p, 0750)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return "", err
	}
	return p, nil
}

func DirExists(p string) (bool, error) {
	stat, err := os.Stat(p)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return false, err
	} else if err != nil {
		return false, nil
	}
	return stat.IsDir(), nil
}
