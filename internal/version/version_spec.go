package version

import (
	"encoding/json"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Comparison int

const (
	ComparisonEqual Comparison = iota
	ComparisonLess
	ComparisonLessOrEqual
	ComparisonGreater
	ComparisonGreaterOrEqual
)

func (c Comparison) String() string {
	switch c {
	case ComparisonEqual:
		return ""
	case ComparisonLess:
		return "<"
	case ComparisonLessOrEqual:
		return "<="
	case ComparisonGreater:
		return ">"
	case ComparisonGreaterOrEqual:
		return ">="
	default:
		panic(fmt.Sprintf("unknown comparison %d", c))
	}
}

type VersionSpec struct {
	Comparison Comparison
	Version    Version
}

func (vs VersionSpec) String() string {
	return fmt.Sprintf("%s%s", vs.Comparison, vs.Version)
}

func (vs *VersionSpec) Validate() error {
	if vs.Comparison < ComparisonEqual || vs.Comparison > ComparisonGreaterOrEqual {
		return fmt.Errorf("invalid comparison %d", vs.Comparison)
	}
	if err := vs.Version.Validate(); err != nil {
		return err
	}
	return nil
}

func (vs *VersionSpec) Match(v Version) bool {
	switch vs.Comparison {
	case ComparisonEqual:
		return v.Major == vs.Version.Major && v.Minor == vs.Version.Minor
	case ComparisonLess:
		switch {
		case v.Major < vs.Version.Major:
			return true
		case v.Major > vs.Version.Major:
			return false
		case v.Minor < vs.Version.Minor:
			return true
		case v.Minor > vs.Version.Minor:
			return false
		default:
			return false
		}
	case ComparisonLessOrEqual:
		switch {
		case v.Major < vs.Version.Major:
			return true
		case v.Major > vs.Version.Major:
			return false
		case v.Minor < vs.Version.Minor:
			return true
		case v.Minor > vs.Version.Minor:
			return false
		default:
			return true
		}
	case ComparisonGreater:
		switch {
		case v.Major < vs.Version.Major:
			return false
		case v.Major > vs.Version.Major:
			return true
		case v.Minor < vs.Version.Minor:
			return false
		case v.Minor > vs.Version.Minor:
			return true
		default:
			return false
		}
	case ComparisonGreaterOrEqual:
		switch {
		case v.Major < vs.Version.Major:
			return false
		case v.Major > vs.Version.Major:
			return true
		case v.Minor < vs.Version.Minor:
			return false
		case v.Minor > vs.Version.Minor:
			return true
		default:
			return true
		}
	default:
		panic(fmt.Sprintf("unknown comparison %d", vs.Comparison))
	}
}

func (vs *VersionSpec) UnmarshalJSON(b []byte) error {
	var s string
	var err error
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*vs, err = VersionSpecFromString(s)
	if err != nil {
		return err
	}
	return nil
}

func (vs *VersionSpec) UnmarshalYAML(node *yaml.Node) error {
	var err error
	*vs, err = VersionSpecFromString(node.Value)
	if err != nil {
		return err
	}
	return nil
}

var versionSpecRe = regexp.MustCompile(`^([><]=?)?(.+)$`)

func VersionSpecFromString(s string) (VersionSpec, error) {
	var versionSpec VersionSpec
	var err error

	matches := versionSpecRe.FindStringSubmatch(s)
	if len(matches) == 0 {
		return versionSpec, fmt.Errorf("invalid version spec %q", s)
	}

	versionSpec.Comparison = ComparisonEqual
	if len(matches) == 3 {
		switch matches[1] {
		case "<":
			versionSpec.Comparison = ComparisonLess
		case "<=":
			versionSpec.Comparison = ComparisonLessOrEqual
		case ">":
			versionSpec.Comparison = ComparisonGreater
		case ">=":
			versionSpec.Comparison = ComparisonGreaterOrEqual
		}
		matches = matches[2:]
	}

	versionSpec.Version, err = VersionFromString(matches[0])
	if err != nil {
		return versionSpec, err
	}

	if err := versionSpec.Validate(); err != nil {
		return versionSpec, fmt.Errorf("invalid version spec %q: %s", s, err)
	}

	return versionSpec, nil
}
