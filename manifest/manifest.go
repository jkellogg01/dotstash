package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/dotstash/files"
)

var (
	dotstashPath,
	configPath,
	homePath string
)

const (
	dotstashPathAbbr = "#dst"
	homePathAbbr     = "#hom"
	configPathAbbr   = "#cfg"
)

type ConfigMetadata struct {
	Author  string         `json:"author"`
	URL     string         `json:"url"`
	Branch  string         `json:"branch"`
	Targets []ConfigTarget `json:"targets"`
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
	err = dstFile.Chmod(0o640)
	if err != nil {
		log.Warnf("could not change permissions for manifest at %s", dst)
	}
	if info, err := dstFile.Stat(); err == nil {
		log.Debugf("manifest permissions: %s", info.Mode().Perm().String())
	} else {
		log.Error(err)
	}
	return d.writeJson(dstFile)
}

// NOTE: path should not include 'manifest.json'
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
	if err != nil {
		return nil, err
	}
	if err := srcFile.Close(); err != nil {
		log.Errorf("failed to close manifest file after reading: %s", err)
	}
	return result, nil
}

func (d ConfigMetadata) Link(clobber bool) error {
	var err error
	for _, t := range d.ExpandTargets() {
		err = errors.Join(err, files.Link(t.Src, t.Dst, clobber))
	}
	return err
}

func (d ConfigMetadata) Unlink() error {
	var err error
	for _, t := range d.ExpandTargets() {
		err = errors.Join(err, files.Unlink(t.Dst))
	}
	return err
}

// RemoveTarget takes the basename of the target as it appears in the config repository.
func (d *ConfigMetadata) RemoveTarget(base string) {
	targets := d.ExpandTargets()
	for i, target := range targets {
		if filepath.Base(target.Src) != base {
			continue
		}
		targets = slices.Concat(targets[:i], targets[i+1:])
	}
	d.Targets = targets
}

// THIS IS ALMOST ALWAYS THE CORRECT WAY TO READ A MANIFEST'S TARGETS
func (d ConfigMetadata) ExpandTargets() []ConfigTarget {
	result := make([]ConfigTarget, 0, len(d.Targets))
	for _, t := range d.Targets {
		result = append(result, t.expand())
	}
	return result
}

// NOTE: in this case src and dst refer to the source inside of the config repository,
// and the destination in the user's configuration files.
// For example, if I was storing my neovim config:
// src = "nvim"
// dst = "#hom/.config/nvim"
func (d *ConfigMetadata) AppendTarget(src, dst string) {
	newTarget := ConfigTarget{
		Src: compressPath(src),
		Dst: compressPath(dst),
	}
	log.Debug("appending target", "target", newTarget.String())
	d.Targets = append(d.Targets, newTarget)
}

// NOTE: unlike expansion, compression must happen in order from deepest to shallowest path.
func compressPath(p string) string {
	var (
		post string
		cut  bool
	)
	post, cut = strings.CutPrefix(p, dotstashPath)
	if cut {
		return filepath.Join(dotstashPathAbbr, post)
	}
	post, cut = strings.CutPrefix(p, configPath)
	if cut {
		return filepath.Join(configPathAbbr, post)
	}
	post, cut = strings.CutPrefix(p, homePath)
	if cut {
		return filepath.Join(homePathAbbr, post)
	}
	return p
}

func (t ConfigTarget) expand() ConfigTarget {
	return ConfigTarget{
		Src: expandPath(t.Src),
		Dst: expandPath(t.Dst),
	}
}

func expandPath(p string) string {
	var (
		post string
		cut  bool
	)
	post, cut = strings.CutPrefix(p, dotstashPathAbbr)
	if cut {
		return filepath.Join(dotstashPath, post)
	}
	post, cut = strings.CutPrefix(p, configPathAbbr)
	if cut {
		return filepath.Join(configPath, post)
	}
	post, cut = strings.CutPrefix(p, homePathAbbr)
	if cut {
		return filepath.Join(homePath, post)
	}
	return p
}

func (d *ConfigMetadata) writeJson(w io.Writer) error {
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	n, err := w.Write(data)
	if n != len(data) {
		return errors.New("did not write enough data!")
	}
	return err
}

func readJson(w io.Reader) (*ConfigMetadata, error) {
	decoder := json.NewDecoder(w)
	result := ConfigMetadata{}
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func init() {
	var err error
	dotstashPath, err = files.GetDotstashPath()
	if err != nil {
		panic(err)
	}
	configPath, err = os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	homePath, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
}

func (t ConfigTarget) String() string {
	return fmt.Sprintf("%s => %s", t.Src, t.Dst)
}
