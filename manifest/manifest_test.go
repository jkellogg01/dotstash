package manifest

import (
	"testing"
)

func TestCompress(t *testing.T) {
	testCases := []struct {
		name   string
		data   string
		expect string
	}{
		{"home directory compression", homePath, homePathAbbr},
		{"config directory compression", configPath, configPathAbbr},
		{"dotstash directory compression", dotstashPath, dotstashPathAbbr},
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
		{"home directory expansion", homePathAbbr, homePath},
		{"config directory expansion", configPathAbbr, configPath},
		{"dotstash directory expansion", dotstashPathAbbr, dotstashPath},
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
