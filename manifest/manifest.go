package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/files"
)

const (
	figurePathAbbr = "#fig"
	homePathAbbr   = "#hom"
	configPathAbbr = "#cfg"
)

type pathfunc func() (string, error)

var pathMap = map[string]pathfunc{
	figurePathAbbr: files.GetFigurePath,
	homePathAbbr:   os.UserHomeDir,
	configPathAbbr: os.UserConfigDir,
}

type ConfigMetadata struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Targets   []ConfigTarget `json:"targets"`
}

type ConfigTarget struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

// NOTE: path should point to the config DIRECTORY; it will create manifest.json on its own.
func (d *ConfigMetadata) EmitManifest(path string) error {
	dst := filepath.Join(path, "manifest.json")
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	// err = dstFile.Chmod(0o664)
	// if err != nil {
	// 	log.Warnf("could not change permissions for manifest at %s", dst)
	// }
	if info, err := dstFile.Stat(); err == nil {
		log.Debugf("manifest permissions: %s", info.Mode().Perm().String())
	} else {
		log.Error(err)
	}
	return d.writeJson(dstFile)
}

func ReadManifest(path string) (*ConfigMetadata, error) {
	src := filepath.Join(path, "manifest.json")
	srcFile, err := os.Open(src)
	if errors.Is(err, fs.ErrPermission) {
		info, statErr := os.Stat(src)
		if statErr != nil {
			return nil, fmt.Errorf("Tried and failed to get more info about the problematic file. Original error: %w", err)
		}
		perms := info.Mode().Perm().String()
		return nil, fmt.Errorf("%w\n\tProblematic file permissions: %s", err, perms)
	} else if err != nil {
		return nil, err
	}
	result, err := readJson(srcFile)
	if err := srcFile.Close(); err != nil {
		log.Errorf("failed to close manifest file after reading: %s", err)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *ConfigMetadata) writeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(d)
}

func readJson(w io.Reader) (*ConfigMetadata, error) {
	decoder := json.NewDecoder(w)
	result := new(ConfigMetadata)
	err := decoder.Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// NOTE: in this case src and dst refer to the source inside of the config repository,
// and the destination in the user's configuration files.
// For example, if I was storing my neovim config:
// src = "nvim"
// dst = "#hom/.config/nvim"
func (d *ConfigMetadata) AppendTarget(src, dst string) error {
	var toAppend ConfigTarget

	figPath, err := files.GetFigurePath()
	if err != nil {
		return err
	}
	post, cut := strings.CutPrefix(src, figPath)
	if !cut {
		log.Warn("source name did not begin with figure's home path, this may be indicative of a problem with the app.", "src", src)
	}
	toAppend.Src = post

	// NOTE: we open the home path first even though we don't check for it later because
	// the os package docs say that a failure to locate the home directory is the most
	// common failure mode for os.UserConfigDir
	homPath, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cfgPath, err := os.UserConfigDir()
	if err != nil {
		log.Error("config directory could not be determined", "error", err)
		return err
	}
	var compressedDst string
	postConfig, cfgCut := strings.CutPrefix(dst, cfgPath)
	postHome, homCut := strings.CutPrefix(dst, homPath)
	switch {
	case cfgCut:
		compressedDst = filepath.Join(configPathAbbr, postConfig)
	case homCut:
		compressedDst = filepath.Join(homePathAbbr, postHome)
	default:
		compressedDst = dst
	}
	toAppend.Dst = compressedDst

	d.Targets = append(d.Targets, toAppend)
	return nil
}
