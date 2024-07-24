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

func LinkSubstitute(oldPath, newPath string) error {
	backup, err := MakeTempFallback(oldPath)
	if err != nil {
		return err
	}
	defer backup.Close()
	backupInfo, err := backup.Stat()
	if err != nil {
		return err
	}
	backupName := backupInfo.Name()
	err = os.Rename(oldPath, newPath)
	if err != nil {
		log.Errorf("failed to move %s to %s. deleting backup and moving on...", oldPath, newPath)
		cleanupErr := os.RemoveAll(backupName)
		if cleanupErr != nil {
			log.Errorf("failed to clean up backup: %s", cleanupErr)
		}
		return err
	}
	err = errors.Join(
		os.Symlink(newPath, oldPath),
		os.Chmod(oldPath, 0o700),
	)
	if err == nil {
		err = os.RemoveAll(backupName)
		if err != nil {
			log.Errorf("backup not cleaned up: %s", err)
		}
		return nil
	}
	restoreBackupError := os.Rename(backupName, oldPath)
	err = errors.Join(err, restoreBackupError)
	if restoreBackupError != nil {
		log.Errorf("failed to restore %s from backup. backup is located at: %s", oldPath, backup)
	}
	return err
}
