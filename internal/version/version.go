package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Version struct {
	Major uint64
	Minor uint64
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

func (v *Version) Validate() error {
	if v.Major == 0 && v.Minor == 0 {
		return errors.New("invalid version 0.0")
	}
	return nil
}

func (v *Version) GreaterThan(otherVersion Version) bool {
	switch {
	case v.Major > otherVersion.Major:
		return true
	case v.Major < otherVersion.Major:
		return false
	case v.Minor > otherVersion.Minor:
		return true
	case v.Minor < otherVersion.Minor:
		return false
	default:
		return false
	}
}

func (v *Version) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("parse version from %q: %s", b, err)
	}
	version, err := VersionFromString(s)
	if err != nil {
		return err
	}
	*v = version
	return nil
}

func (v *Version) UnmarshalYAML(node *yaml.Node) error {
	version, err := VersionFromString(node.Value)
	if err != nil {
		return err
	}
	*v = version
	return nil
}

func VersionFromString(s string) (Version, error) {
	var version Version
	var err error

	parts := strings.SplitN(s, ".", 2)
	if len(parts) != 2 {
		return version, fmt.Errorf("invalid version %q", s)
	}

	version.Major, err = strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return version, fmt.Errorf("invalid major version %q", parts[0])
	}
	version.Minor, err = strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return version, fmt.Errorf("invalid minor version %q", parts[1])
	}

	if err := version.Validate(); err != nil {
		return version, err
	}

	return version, nil
}
