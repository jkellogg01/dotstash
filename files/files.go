package files

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

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

// remove the file or directory at `dst`, replacing it the file or directory at `src`
func Substitute(src, dst string) error {
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
