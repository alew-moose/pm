package version

import "testing"

func TestVersionFromString(t *testing.T) {
	tests := []struct {
		str         string
		wantVersion Version
		wantErr     bool
	}{
		{str: "", wantErr: true},
		{str: "1", wantErr: true},
		{str: "-1", wantErr: true},
		{str: "abcde", wantErr: true},
		{str: "0.0", wantErr: true},
		{str: " 1 . 2 ", wantErr: true},
		{str: "0.1", wantVersion: Version{Major: 0, Minor: 1}},
		{str: "001.0001", wantVersion: Version{Major: 1, Minor: 1}},
		{str: "1000.1000", wantVersion: Version{Major: 1000, Minor: 1000}},
	}

	for ti, tt := range tests {
		version, err := VersionFromString(tt.str)
		if err != nil && !tt.wantErr {
			t.Errorf("failed test #%d: VersionFromString(%q) returned error %q", ti, tt.str, err)
		}
		if err == nil && tt.wantErr {
			t.Errorf("failed test #%d: expected error from VersionFromString(%q)", ti, tt.str)
		}
		if !tt.wantErr && version != tt.wantVersion {
			t.Errorf("failed test #%d: VersionFromString(%q): got %#v, want %#v", ti, tt.str, version, tt.wantVersion)
		}
	}
}
