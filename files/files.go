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
	path := PathAppend(homeDir, ".figure")
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

func PathAppend(path string, names ...string) string {
	if len(names) == 0 {
		return path
	}
	segments := strings.Split(path, string(os.PathSeparator))
	for _, name := range names {
		segments = append(segments, name)
	}
	return strings.Join(segments, string(os.PathSeparator))
}
