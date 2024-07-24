package files

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

func MakeTempFallback(src string) (fs.ReadDirFile, error) {
	log.Debugf("making a fallback for %s", src)
	info, err := os.Stat(src)
	if err != nil {
		log.Errorf("failed to get info about %s", src)
		return nil, err
	}
	splitPath := strings.Split(strings.Trim(filepath.ToSlash(src), "/"), "/")
	tempNamePattern := strings.Join(splitPath, "_")
	mode := info.Mode()
	switch {
	case mode.IsRegular():
		log.Debugf("%s is a regular file", src)
		f, err := os.CreateTemp("", tempNamePattern)
		if err != nil {
			return nil, err
		}
		err = copyFile(src, f.Name())
		return wrappedFile{f}, err
	case mode.IsDir():
		log.Debugf("%s is a directory", src)
		dirName, err := os.MkdirTemp("", tempNamePattern)
		if err != nil {
			return nil, err
		}
		dir, err := os.Open(dirName)
		if err != nil {
			panic("created a temp dir but could not get a pointer to it!")
		}
		err = deepCopy(src, dirName)
		return dir, err
	default:
		return nil, fmt.Errorf("source file %s is neither a directory nor a regular file. making fallbacks of other types of files is not supported at this time.", src)
	}
}

func deepCopy(src, dst string) error {
	// log.Debugf("starting deep copy: %s => %s", src, dst)
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if sfi.Mode().IsRegular() {
		return copyFile(src, dst)
	} else if !sfi.IsDir() {
		return fmt.Errorf("%s is neiter a directory nor a regular file", src)
	}
	err = os.Mkdir(dst, 0o700)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}
	srcDir, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range srcDir {
		err := deepCopy(
			filepath.Join(src, entry.Name()),
			filepath.Join(dst, entry.Name()),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	log.Debugf("starting file copy: %s => %s", src, dst)
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	} else if !sfi.Mode().IsRegular() {
		return fmt.Errorf("[copyFile] source file %s is not regular", src)
	}
	dfi, err := os.Stat(dst)
	if errors.Is(err, fs.ErrNotExist) {
		// noop
	} else if err != nil {
		return err
	} else if !dfi.Mode().IsRegular() {
		return fmt.Errorf("[copyFile] destination file %s is not regular", dst)
	}

	if os.SameFile(sfi, dfi) {
		return nil
	}

	if err := os.Link(src, dst); err == nil {
		return nil
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

type wrappedFile struct {
	*os.File
}

func (f wrappedFile) ReadDir(int) ([]fs.DirEntry, error) {
	return nil, errors.New("this file is not a directory")
}
