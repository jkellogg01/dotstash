package manifest

import (
	"os"
	"testing"

	"github.com/jkellogg01/dotstash/files"
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
	homeDir, err = os.UserHomeDir()
	if err != nil {
		panic("could not retrieve home directory for testing")
	}
	cfgDir, err = os.UserConfigDir()
	if err != nil {
		panic("could not retrieve config directory for testing")
	}
	dsDir, err = files.GetDotstashPath()
	if err != nil {
		panic("could not retrieve dotstash directory for testing")
	}
}
