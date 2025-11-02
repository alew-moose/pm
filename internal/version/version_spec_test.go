package version

import (
	"testing"
)

func TestVersionSpecFromString(t *testing.T) {
	tests := []struct {
		str             string
		wantVersionSpec VersionSpec
		wantErr         bool
	}{
		{str: "0.0", wantErr: true},
		{str: "<>1.2", wantErr: true},
		{str: "=1.2", wantErr: true},
		{str: ">>1.2", wantErr: true},
		{str: "?1.2", wantErr: true},
		{str: "x1.2", wantErr: true},
		{str: "1.2", wantVersionSpec: VersionSpec{Comparison: ComparisonEqual, Version: Version{Major: 1, Minor: 2}}},
		{str: "<1.2", wantVersionSpec: VersionSpec{Comparison: ComparisonLess, Version: Version{Major: 1, Minor: 2}}},
		{str: "<=1.2", wantVersionSpec: VersionSpec{Comparison: ComparisonLessOrEqual, Version: Version{Major: 1, Minor: 2}}},
		{str: ">1.2", wantVersionSpec: VersionSpec{Comparison: ComparisonGreater, Version: Version{Major: 1, Minor: 2}}},
		{str: ">=1.2", wantVersionSpec: VersionSpec{Comparison: ComparisonGreaterOrEqual, Version: Version{Major: 1, Minor: 2}}},
	}
	for ti, tt := range tests {
		versionSpec, err := VersionSpecFromString(tt.str)
		if err != nil && !tt.wantErr {
			t.Errorf("failed test #%d: VersionSpecFromString(%q) returned error %q", ti, tt.str, err)
		}
		if err == nil && tt.wantErr {
			t.Errorf("failed test #%d: expected error from VersionSpecFromString(%q)", ti, tt.str)
		}
		if !tt.wantErr && versionSpec != tt.wantVersionSpec {
			t.Errorf("failed test #%d: VersionSpecFromString(%q): got %#v, want %#v", ti, tt.str, versionSpec, tt.wantVersionSpec)
		}
	}
}

func TestVersionSpecMatch(t *testing.T) {
	tests := []struct {
		versionStr     string
		versionSpecStr string
		wantMatch      bool
	}{
		{versionStr: "1.1", versionSpecStr: "1.1", wantMatch: true},
		{versionStr: "1.1", versionSpecStr: "<1.1", wantMatch: false},
		{versionStr: "1.1", versionSpecStr: "<=1.1", wantMatch: true},
		{versionStr: "1.1", versionSpecStr: ">1.1", wantMatch: false},
		{versionStr: "1.1", versionSpecStr: ">=1.1", wantMatch: true},

		{versionStr: "2.1", versionSpecStr: "1.1", wantMatch: false},
		{versionStr: "2.1", versionSpecStr: ">=1.1", wantMatch: true},
		{versionStr: "2.1", versionSpecStr: ">1.1", wantMatch: true},
		{versionStr: "2.1", versionSpecStr: "<1.1", wantMatch: false},
		{versionStr: "2.1", versionSpecStr: "<=1.1", wantMatch: false},

		{versionStr: "1.1", versionSpecStr: "2.1", wantMatch: false},
		{versionStr: "1.1", versionSpecStr: ">=2.1", wantMatch: false},
		{versionStr: "1.1", versionSpecStr: ">2.1", wantMatch: false},
		{versionStr: "1.1", versionSpecStr: "<2.1", wantMatch: true},
		{versionStr: "1.1", versionSpecStr: "<=2.1", wantMatch: true},

		{versionStr: "1.2", versionSpecStr: "1.1", wantMatch: false},
		{versionStr: "1.2", versionSpecStr: ">=1.1", wantMatch: true},
		{versionStr: "1.2", versionSpecStr: ">1.1", wantMatch: true},
		{versionStr: "1.2", versionSpecStr: "<1.1", wantMatch: false},
		{versionStr: "1.2", versionSpecStr: "<=1.1", wantMatch: false},

		{versionStr: "1.1", versionSpecStr: "1.2", wantMatch: false},
		{versionStr: "1.1", versionSpecStr: ">=1.2", wantMatch: false},
		{versionStr: "1.1", versionSpecStr: ">1.2", wantMatch: false},
		{versionStr: "1.1", versionSpecStr: "<1.2", wantMatch: true},
		{versionStr: "1.1", versionSpecStr: "<=1.2", wantMatch: true},
	}

	for ti, tt := range tests {
		version, err := VersionFromString(tt.versionStr)
		if err != nil {
			t.Errorf("failed test #%d: VersionFromString(%q) returned error %q", ti, tt.versionStr, err)
		}

		versionSpec, err := VersionSpecFromString(tt.versionSpecStr)
		if err != nil {
			t.Errorf("failed test #%d: VersionSpecFromString(%q) returned error %q", ti, tt.versionSpecStr, err)
		}

		match := versionSpec.Match(version)
		if match != tt.wantMatch {
			t.Errorf("failed test #%d: version %q match versionSpec %q: got %t, wanted %t", ti, tt.versionStr, tt.versionSpecStr, match, tt.wantMatch)
		}
	}
}
