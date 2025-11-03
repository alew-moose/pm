package version

import (
	"testing"
)

func TestVersionConstraintFromString(t *testing.T) {
	tests := []struct {
		str                   string
		wantVersionConstraint VersionConstraint
		wantErr               bool
	}{
		{str: "0.0", wantErr: true},
		{str: "<>1.2", wantErr: true},
		{str: ">>1.2", wantErr: true},
		{str: "?1.2", wantErr: true},
		{str: "x1.2", wantErr: true},
		{str: "<=1.2 =>1.3", wantErr: true},
		{str: "<1.2 >1.3", wantErr: true},
		{str: ">1.3 <1.2", wantErr: true},
		{str: "1.2", wantVersionConstraint: VersionConstraint{
			Exact: &ExactVersionConstraint{
				Version: Version{
					Major: 1,
					Minor: 2,
				},
			},
		}},
		{str: "=1.2", wantVersionConstraint: VersionConstraint{
			Exact: &ExactVersionConstraint{
				Version: Version{
					Major: 1,
					Minor: 2,
				},
			},
		}},
		{str: "<1.2", wantVersionConstraint: VersionConstraint{
			Upper: &UpperVersionConstraint{
				Comparison: ComparisonLess,
				Version: Version{
					Major: 1,
					Minor: 2,
				},
			},
		}},
		{str: "<=1.2", wantVersionConstraint: VersionConstraint{
			Upper: &UpperVersionConstraint{
				Comparison: ComparisonLessOrEqual,
				Version: Version{
					Major: 1,
					Minor: 2,
				},
			},
		}},
		{str: ">1.2", wantVersionConstraint: VersionConstraint{
			Lower: &LowerVersionConstraint{
				Comparison: ComparisonGreater,
				Version: Version{
					Major: 1,
					Minor: 2,
				},
			},
		}},
		{str: ">=1.2", wantVersionConstraint: VersionConstraint{
			Lower: &LowerVersionConstraint{
				Comparison: ComparisonGreaterOrEqual,
				Version: Version{
					Major: 1,
					Minor: 2,
				},
			},
		}},
	}
	for ti, tt := range tests {
		versionConstraint, err := VersionConstraintFromString(tt.str)

		if err != nil && !tt.wantErr {
			t.Errorf("failed test #%d: VersionConstraintFromString(%q) returned error %q", ti, tt.str, err)
			continue
		}
		if err == nil && tt.wantErr {
			t.Errorf("failed test #%d: expected error from VersionConstraintFromString(%q)", ti, tt.str)
			continue
		}
		if !tt.wantErr && !versionConstraint.Equal(tt.wantVersionConstraint) {
			t.Errorf("failed test #%d: VersionConstraintFromString(%q): got %s, want %s", ti, tt.str, versionConstraint, tt.wantVersionConstraint)
		}
	}
}

// func TestVersionConstraintMatch(t *testing.T) {
// 	tests := []struct {
// 		versionStr           string
// 		versionConstraintStr string
// 		wantMatch            bool
// 	}{
// 		{versionStr: "1.1", versionConstraintStr: "1.1", wantMatch: true},
// 		{versionStr: "1.1", versionConstraintStr: "<1.1", wantMatch: false},
// 		{versionStr: "1.1", versionConstraintStr: "<=1.1", wantMatch: true},
// 		{versionStr: "1.1", versionConstraintStr: ">1.1", wantMatch: false},
// 		{versionStr: "1.1", versionConstraintStr: ">=1.1", wantMatch: true},

// 		{versionStr: "2.1", versionConstraintStr: "1.1", wantMatch: false},
// 		{versionStr: "2.1", versionConstraintStr: ">=1.1", wantMatch: true},
// 		{versionStr: "2.1", versionConstraintStr: ">1.1", wantMatch: true},
// 		{versionStr: "2.1", versionConstraintStr: "<1.1", wantMatch: false},
// 		{versionStr: "2.1", versionConstraintStr: "<=1.1", wantMatch: false},

// 		{versionStr: "1.1", versionConstraintStr: "2.1", wantMatch: false},
// 		{versionStr: "1.1", versionConstraintStr: ">=2.1", wantMatch: false},
// 		{versionStr: "1.1", versionConstraintStr: ">2.1", wantMatch: false},
// 		{versionStr: "1.1", versionConstraintStr: "<2.1", wantMatch: true},
// 		{versionStr: "1.1", versionConstraintStr: "<=2.1", wantMatch: true},

// 		{versionStr: "1.2", versionConstraintStr: "1.1", wantMatch: false},
// 		{versionStr: "1.2", versionConstraintStr: ">=1.1", wantMatch: true},
// 		{versionStr: "1.2", versionConstraintStr: ">1.1", wantMatch: true},
// 		{versionStr: "1.2", versionConstraintStr: "<1.1", wantMatch: false},
// 		{versionStr: "1.2", versionConstraintStr: "<=1.1", wantMatch: false},

// 		{versionStr: "1.1", versionConstraintStr: "1.2", wantMatch: false},
// 		{versionStr: "1.1", versionConstraintStr: ">=1.2", wantMatch: false},
// 		{versionStr: "1.1", versionConstraintStr: ">1.2", wantMatch: false},
// 		{versionStr: "1.1", versionConstraintStr: "<1.2", wantMatch: true},
// 		{versionStr: "1.1", versionConstraintStr: "<=1.2", wantMatch: true},
// 	}

// 	for ti, tt := range tests {
// 		version, err := VersionFromString(tt.versionStr)
// 		if err != nil {
// 			t.Errorf("failed test #%d: VersionFromString(%q) returned error %q", ti, tt.versionStr, err)
// continue
// 		}

// 		versionConstraint, err := VersionConstraintFromString(tt.versionConstraintStr)
// 		if err != nil {
// 			t.Errorf("failed test #%d: VersionConstraintFromString(%q) returned error %q", ti, tt.versionConstraintStr, err)
// continue
// 		}

// 		match := versionConstraint.Match(version)
// 		if match != tt.wantMatch {
// 			t.Errorf("failed test #%d: version %q match versionConstraint %q: got %t, wanted %t", ti, tt.versionStr, tt.versionConstraintStr, match, tt.wantMatch)
// 		}
// 	}
// }
