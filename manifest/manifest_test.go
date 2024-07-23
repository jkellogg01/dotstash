package manifest

import (
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/files"
)

var (
	homeDir string
	cfgDir  string
	dsDir   string
)

func TestCompress(t *testing.T) {
	testCases := []struct {
		name   string
		data   string
		expect string
	}{
		{"home directory compression", homeDir, homePathAbbr},
		{"config directory compression", cfgDir, configPathAbbr},
		{"dotstash directory compression", dsDir, dotstashPathAbbr},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := compressPath(tc.data)
			if result != tc.expect {
				t.Errorf("Expected %s, got %s", tc.expect, result)
			}
		})
	}
}

func TestExpand(t *testing.T) {
	testCases := []struct {
		name   string
		data   string
		expect string
	}{
		{"home directory expansion", homePathAbbr, homeDir},
		{"config directory expansion", configPathAbbr, cfgDir},
		{"dotstash directory expansion", dotstashPathAbbr, dsDir},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := expandPath(tc.data)
			if result != tc.expect {
				t.Errorf("Expected %s, got %s", tc.expect, result)
			}
		})
	}
}

func init() {
	var err error
	tlog := log.NewWithOptions(os.Stderr, log.Options{
		Prefix: "[TESTING]",
	})
	homeDir, err = os.UserHomeDir()
	if err != nil {
		tlog.Fatal("could not retrieve home directory for testing")
	}
	cfgDir, err = os.UserConfigDir()
	if err != nil {
		tlog.Fatal("could not retrieve config directory for testing")
	}
	dsDir, err = files.GetFigurePath()
	if err != nil {
		tlog.Fatal("could not retrieve dotstash directory for testing")
	}
}
