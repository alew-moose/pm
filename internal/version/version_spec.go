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

type VersionSpec struct {
	Version    Version
	Comparison Comparison
}

func (vs *VersionSpec) Match(v Version) bool {
	switch vs.Comparison {
	case ComparisonEqual:
		return v.Major == vs.Version.Major && v.Minor == vs.Version.Minor
	case ComparisonLess:
		if v.Major < vs.Version.Major {
			return true
		}
		if v.Major > vs.Version.Major {
			return false
		}
		if v.Minor < vs.Version.Minor {
			return true
		}
		if v.Minor > vs.Version.Minor {
			return false
		}
		return false
	case ComparisonLessOrEqual:
		if v.Major < vs.Version.Major {
			return true
		}
		if v.Major > vs.Version.Major {
			return false
		}
		if v.Minor < vs.Version.Minor {
			return true
		}
		if v.Minor > vs.Version.Minor {
			return false
		}
		return true
	case ComparisonGreater:
		if v.Major < vs.Version.Major {
			return false
		}
		if v.Major > vs.Version.Major {
			return true
		}
		if v.Minor < vs.Version.Minor {
			return false
		}
		if v.Minor > vs.Version.Minor {
			return true
		}
		return false
	case ComparisonGreaterOrEqual:
		if v.Major < vs.Version.Major {
			return false
		}
		if v.Major > vs.Version.Major {
			return true
		}
		if v.Minor < vs.Version.Minor {
			return false
		}
		if v.Minor > vs.Version.Minor {
			return true
		}
		return true
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

	return versionSpec, nil
}
