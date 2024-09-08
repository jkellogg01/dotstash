package manifest

import (
	"path/filepath"
	"testing"
)

func TestCompress(t *testing.T) {
	testSegment := "test_segment"
	testCases := []struct {
		name   string
		data   string
		expect string
	}{
		{"home directory compression", homePath, homePathAbbr},
		{"home subdirectory compression", filepath.Join(homePath, "testcase"), filepath.Join(homePathAbbr, "testcase")},
		{"config directory compression", configPath, configPathAbbr},
		{"config subdirectory compression", filepath.Join(configPath, "testcase"), filepath.Join(configPathAbbr, "testcase")},
		{"dotstash directory expansion", filepath.Join(dotstashPath, testSegment), dotstashPathAbbr},
		{"dotstash subdirectory expansion", filepath.Join(dotstashPath, testSegment, "testcase"), filepath.Join(dotstashPathAbbr, "testcase")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := compressPath(tc.data, testSegment)
			if result != tc.expect {
				t.Errorf("Expected %s, got %s", tc.expect, result)
			}
			t.Logf("compression successful: %s => %s, expect %s", tc.data, result, tc.expect)
		})
	}
}

func TestExpand(t *testing.T) {
	testSegment := "test_segment"
	testCases := []struct {
		name   string
		data   string
		expect string
	}{
		{"home directory expansion", homePathAbbr, homePath},
		{"home subdirectory expansion", filepath.Join(homePathAbbr, "testcase"), filepath.Join(homePath, "testcase")},
		{"config directory expansion", configPathAbbr, configPath},
		{"config subdirectory expansion", filepath.Join(configPathAbbr, "testcase"), filepath.Join(configPath, "testcase")},
		{"dotstash directory expansion", dotstashPathAbbr, filepath.Join(dotstashPath, testSegment)},
		{"dotstash subdirectory expansion", filepath.Join(dotstashPathAbbr, "testcase"), filepath.Join(dotstashPath, testSegment, "testcase")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := expandPath(tc.data, testSegment)
			if result != tc.expect {
				t.Errorf("Expected %s, got %s", tc.expect, result)
			}
			t.Logf("expansion successful: %s => %s, expect %s", tc.data, result, tc.expect)
		})
	}
}
