package manifest

import (
	"path/filepath"
	"testing"
)

func TestCompress(t *testing.T) {
	testCases := []struct {
		name   string
		data   string
		expect string
	}{
		{"home directory compression", homePath, homePathAbbr},
		{"home subdirectory compression", filepath.Join(homePath, "testcase"), filepath.Join(homePathAbbr, "testcase")},
		{"config directory compression", configPath, configPathAbbr},
		{"config subdirectory compression", filepath.Join(configPath, "testcase"), filepath.Join(configPathAbbr, "testcase")},
		{"dotstash directory compression", dotstashPath, dotstashPathAbbr},
		{"dotstash subdirectory compression", filepath.Join(dotstashPath, "testcase"), filepath.Join(dotstashPathAbbr, "testcase")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := compressPath(tc.data)
			if result != tc.expect {
				t.Errorf("Expected %s, got %s", tc.expect, result)
			}
			t.Logf("compression successful: %s => %s, expect %s", tc.data, result, tc.expect)
		})
	}
}

func TestExpand(t *testing.T) {
	testCases := []struct {
		name   string
		data   string
		expect string
	}{
		{"home directory expansion", homePathAbbr, homePath},
		{"home subdirectory expansion", filepath.Join(homePathAbbr, "testcase"), filepath.Join(homePath, "testcase")},
		{"config directory expansion", configPathAbbr, configPath},
		{"config subdirectory expansion", filepath.Join(configPathAbbr, "testcase"), filepath.Join(configPath, "testcase")},
		{"dotstash directory expansion", dotstashPathAbbr, dotstashPath},
		{"dotstash subdirectory expansion", filepath.Join(dotstashPathAbbr, "testcase"), filepath.Join(dotstashPath, "testcase")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := expandPath(tc.data)
			if result != tc.expect {
				t.Errorf("Expected %s, got %s", tc.expect, result)
			}
			t.Logf("expansion successful: %s => %s, expect %s", tc.data, result, tc.expect)
		})
	}
}
