package manifest

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/files"
)

const (
	figurePathAbbr = "#fig"
	homePathAbbr   = "#hom"
	configPathAbbr = "#cfg"
)

type ConfigMetadata struct {
	Targets []ConfigTarget `json:"targets"`
}

type ConfigTarget struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

// NOTE: d.Emit should point to the config DIRECTORY; it will create manifest.json on its own.
func (d *ConfigMetadata) EmitManifest(path string) error {
	dst := filepath.Join(path, "manifest.json")
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := dstFile.Close()
		if err == nil {
			err = closeErr
		}
	}()
	err = d.writeJson(dstFile)
	return err
}

func (d *ConfigMetadata) writeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(d)
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
