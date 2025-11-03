package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

// TODO: validate?

// TODO: c -> vc || vc -> c

type LowerVersionComparison int

const (
	ComparisonGreater LowerVersionComparison = iota + 1
	ComparisonGreaterOrEqual
)

func (c LowerVersionComparison) String() string {
	switch c {
	case ComparisonGreater:
		return ">"
	case ComparisonGreaterOrEqual:
		return ">="
	default:
		panic(fmt.Sprintf("unknown lower version comparison %d", c))
	}
}

type UpperVersionComparison int

const (
	ComparisonLess UpperVersionComparison = iota + 1
	ComparisonLessOrEqual
)

func (c UpperVersionComparison) String() string {
	switch c {
	case ComparisonLess:
		return "<"
	case ComparisonLessOrEqual:
		return "<="
	default:
		panic(fmt.Sprintf("unknown upper version comparison %d", c))
	}
}

type ExactVersionConstraint struct {
	Version Version
}

func (c *ExactVersionConstraint) String() string {
	// XXX
	// if c == nil {
	// 	return ""
	// }
	return fmt.Sprintf("=%s", c.Version)
}

func (c *ExactVersionConstraint) Match(v Version) bool {
	return c != nil && v == c.Version
}

type LowerVersionConstraint struct {
	Comparison LowerVersionComparison
	Version    Version
}

func (c *LowerVersionConstraint) String() string {
	// XXX
	// if c == nil {
	// 	return ""
	// }
	return fmt.Sprintf("%s%s", c.Comparison, c.Version)
}

func (c *LowerVersionConstraint) Match(v Version) bool {
	if c == nil {
		return true
	}
	switch {
	case v.Major > c.Version.Major:
		return true
	case v.Major < c.Version.Major:
		return false
	case v.Minor > c.Version.Minor:
		return true
	case v.Minor < c.Version.Minor:
		return false
	case c.Comparison == ComparisonGreater:
		return false
	default:
		return true
	}
}

type UpperVersionConstraint struct {
	Comparison UpperVersionComparison
	Version    Version
}

func (c *UpperVersionConstraint) String() string {
	// XXX
	// if c == nil {
	// 	return ""
	// }
	return fmt.Sprintf("%s%s", c.Comparison, c.Version)
}

func (c *UpperVersionConstraint) Match(v Version) bool {
	if c == nil {
		return true
	}
	switch {
	case v.Major > c.Version.Major:
		return false
	case v.Major < c.Version.Major:
		return true
	case v.Minor > c.Version.Minor:
		return false
	case v.Minor < c.Version.Minor:
		return true
	case c.Comparison == ComparisonLess:
		return false
	default:
		return true
	}
}

type VersionConstraint struct {
	Exact *ExactVersionConstraint
	Lower *LowerVersionConstraint
	Upper *UpperVersionConstraint
}

func (c VersionConstraint) String() string {
	if c.Exact != nil {
		return c.Exact.String()
	}
	if c.Lower != nil && c.Upper != nil {
		return fmt.Sprintf("%s %s", c.Lower, c.Upper)
	}
	if c.Lower != nil {
		return c.Lower.String()
	}
	if c.Upper != nil {
		return c.Upper.String()
	}
	// TODO: panic msg
	panic("unreachable")
}

func (c *VersionConstraint) Match(v Version) bool {
	return c.Exact.Match(v) && c.Lower.Match(v) && c.Upper.Match(v)
}

func (c *VersionConstraint) Equal(o VersionConstraint) bool {
	if (c.Exact == nil) != (o.Exact == nil) {
		return false
	}
	if (c.Lower == nil) != (o.Lower == nil) {
		return false
	}
	if (c.Upper == nil) != (o.Upper == nil) {
		return false
	}

	if c.Exact != nil {
		return *c.Exact == *o.Exact
	}
	if c.Lower != nil {
		if *c.Lower != *o.Lower {
			return false
		}
	}
	if c.Upper != nil {
		if *c.Upper != *o.Upper {
			return false
		}
	}

	return true
}

// _______________________________________________1__________2________3____4__________5____6_________7_______________8
var versionConstraintRe = regexp.MustCompile(`^(?:([><]=?|=?)([\d.]+)|(>=?)([\d.]+)\s+(<=?)([\d.]+)|([\d.]+)\s*-\s*([\d.]+))$`)

func VersionConstraintFromString(s string) (VersionConstraint, error) {
	matches := versionConstraintRe.FindStringSubmatch(s)
	fmt.Printf(">>> matches: %#v\n", matches)
	if len(matches) == 0 {
		return VersionConstraint{}, fmt.Errorf("invalid version spec %q", s)
	}

	var exact *ExactVersionConstraint
	var lower *LowerVersionConstraint
	var upper *UpperVersionConstraint

	switch {
	case matches[2] != "":
		fmt.Println(">>> match 2")
		// fmt.Println(">>> HERE 1")
		ver, err := VersionFromString(matches[2])
		if err != nil {
			return VersionConstraint{}, err
		}
		switch matches[1] {
		case "", "=":
			exact = &ExactVersionConstraint{Version: ver}
		case "<":
			upper = &UpperVersionConstraint{
				Comparison: ComparisonLess,
				Version:    ver,
			}
		case "<=":
			upper = &UpperVersionConstraint{
				Comparison: ComparisonLessOrEqual,
				Version:    ver,
			}
		case ">":
			lower = &LowerVersionConstraint{
				Comparison: ComparisonGreater,
				Version:    ver,
			}
		case ">=":
			lower = &LowerVersionConstraint{
				Comparison: ComparisonGreaterOrEqual,
				Version:    ver,
			}
		default:
			// TODO: panic msg
			panic("XXX")
		}
		// fmt.Println(">>> HERE 2")
	case matches[3] != "":
		fmt.Println(">>> match 3")
		verLower, err := VersionFromString(matches[4])
		if err != nil {
			return VersionConstraint{}, err
		}
		verUpper, err := VersionFromString(matches[6])
		if err != nil {
			return VersionConstraint{}, err
		}
		if verLower.GreaterThan(verUpper) {
			return VersionConstraint{}, errors.New("invalid version")
		}
		lower = &LowerVersionConstraint{Version: verLower}
		upper = &UpperVersionConstraint{Version: verUpper}
		switch matches[3] {
		case ">":
			lower.Comparison = ComparisonGreater
		case ">=":
			lower.Comparison = ComparisonGreaterOrEqual
		}
		switch matches[5] {
		case "<":
			upper.Comparison = ComparisonLess
		case "<=":
			upper.Comparison = ComparisonLessOrEqual
		}
	case matches[7] != "":
		fmt.Println(">>> match 7")
		verLower, err := VersionFromString(matches[7])
		if err != nil {
			return VersionConstraint{}, err
		}
		verUpper, err := VersionFromString(matches[8])
		if err != nil {
			return VersionConstraint{}, err
		}
		if verLower.GreaterThan(verUpper) {
			return VersionConstraint{}, errors.New("invalid version")
		}
		lower = &LowerVersionConstraint{
			Comparison: ComparisonGreaterOrEqual,
			Version:    verLower,
		}
		upper = &UpperVersionConstraint{
			Comparison: ComparisonLessOrEqual,
			Version:    verUpper,
		}
	default:
		// TODO: panic msg
		panic("unreachable XXX")
	}

	// fmt.Println(">>> HERE")

	return VersionConstraint{
		Exact: exact,
		Lower: lower,
		Upper: upper,
	}, nil
}

func (vc *VersionConstraint) UnmarshalJSON(b []byte) error {
	var s string
	var err error
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*vc, err = VersionConstraintFromString(s)
	if err != nil {
		return err
	}
	return nil
}

func (vc *VersionConstraint) UnmarshalYAML(node *yaml.Node) error {
	var err error
	*vc, err = VersionConstraintFromString(node.Value)
	if err != nil {
		return err
	}
	return nil
}
