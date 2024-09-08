package files

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

type ErrNoClobber struct {
	path string
}

func (e ErrNoClobber) Error() string {
	return fmt.Sprintf("encountered a file or directory at %s, with clobbering disabled", e.path)
}

func GetDotstashPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(homeDir, ".dotstash")
	err = os.Mkdir(p, 0o700)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return "", err
	}
	return p, nil
}

// remove the file or directory at `dst`, replacing it with the file or directory at `src`
func Substitute(src, dst string) error {
	// NOTE: os.RemoveAll DOES NOT return fs.ErrNotExist if the path does not exist
	err := os.RemoveAll(dst)
	if err != nil {
		return err
	}
	err = os.Rename(src, dst)
	if err != nil {
		err = errors.Join(
			err,
			os.Symlink(src, dst),
			os.Chmod(dst, 0o700),
		)
		return err
	}
	return nil
}

// remove the file or directory at `src`, moving it to `dst` and placing a symlink to the original file at `src`
func SubstituteForSymlink(src, dst string) error {
	backup, err := MakeTempFallback(src)
	if err != nil {
		return err
	}
	defer backup.Close()
	backupInfo, err := backup.Stat()
	if err != nil {
		return err
	}
	backupName := backupInfo.Name()
	err = os.Rename(src, dst)
	if err != nil {
		log.Errorf("failed to move %s to %s. deleting backup and moving on...", src, dst)
		cleanupErr := os.RemoveAll(backupName)
		if cleanupErr != nil {
			log.Errorf("failed to clean up backup: %s", cleanupErr)
		}
		return err
	}
	err = errors.Join(
		os.Symlink(dst, src),
		os.Chmod(src, 0o700),
	)
	if err == nil {
		err = os.RemoveAll(backupName)
		if err != nil {
			log.Errorf("backup not cleaned up: %s", err)
		}
		return nil
	}
	restoreBackupError := os.Rename(backupName, src)
	err = errors.Join(err, restoreBackupError)
	if restoreBackupError != nil {
		log.Errorf("failed to restore %s from backup. backup is located at: %s", src, backup)
	}
	return err
}

func Link(src, dst string, clobber bool) error {
	dfi, err := os.Stat(dst)
	if errors.Is(err, fs.ErrNotExist) {
		log.Debug("destination file does not exist, nothing to do")
	} else if err != nil {
		return err
	} else if !clobber && (dfi.Mode().IsDir() || dfi.Mode().IsRegular()) {
		// TODO: clobber check
		return ErrNoClobber{dst}
	}

	// TODO: possibly integrate with clobber checking to take an alternate action in that case that a clobber conflict occurs
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if err := os.Symlink(src, dst); err != nil {
		return err
	}
	log.Debug("link successful!", "source", src, "destination", dst)
	return nil
}

func Unlink(dst string) error {
	dfi, err := os.Stat(dst)
	if errors.Is(err, fs.ErrNotExist) {
		// successfully deleted nothing. Everybody's happy!
		return nil
	}
	if dfi.IsDir() || dfi.Mode().IsRegular() {
		return fmt.Errorf("%s does not appear to be a link...", dst)
	}
	return os.Remove(dst)
}
