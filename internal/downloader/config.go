package downloader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alew-moose/pm/internal/version"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Packages []PackageVersionSpec `json:"packages" yaml:"packages"`
}

type PackageVersionSpec struct {
	Name        string              `json:"name" yaml:"name"`
	VersionSpec version.VersionSpec `json:"ver" yaml:"ver"`
}

// красота!
func (pvs *PackageVersionSpec) Match(pv PackageVersion) bool {
	return pvs.Name == pv.Name && pvs.VersionSpec.Match(pv.Version)
}

func (pvs PackageVersionSpec) String() string {
	return fmt.Sprintf("%s(ver %s)", pvs.Name, pvs.VersionSpec)
}

// XXX: copypasted from uploader
var packageNameRe = regexp.MustCompile(`^[\w-]+$`)

func (c *Config) Validate() error {
	packages := make(map[PackageVersionSpec]struct{}, len(c.Packages))
	for _, p := range c.Packages {
		if !packageNameRe.MatchString(p.Name) {
			return fmt.Errorf("invalid package name %q", p.Name)
		}
		if err := p.VersionSpec.Validate(); err != nil {
			return fmt.Errorf("invalid version spec %#v: %q", p.VersionSpec, err)
		}
		if _, ok := packages[p]; ok {
			// TODO: pretty print
			return fmt.Errorf("duplicate package %#v", p)
		}
		packages[p] = struct{}{}
	}
	return nil
}

func ConfigFromFile(path string) (*Config, error) {
	var config *Config
	var err error

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switch ext := filepath.Ext(path); ext {
	case ".json":
		config, err = fromJSON(b)
		if err != nil {
			return nil, fmt.Errorf("parse json: %s", err)
		}
	case ".yaml", ".yml":
		config, err = fromYAML(b)
		if err != nil {
			return nil, fmt.Errorf("parse yaml: %s", err)
		}
	default:
		return nil, fmt.Errorf("%q format is not supported", ext)
	}

	return config, nil
}

func fromJSON(b []byte) (*Config, error) {
	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	fillDefaultVersionSpecs(&config)
	return &config, nil
}

func fromYAML(b []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	fillDefaultVersionSpecs(&config)
	return &config, nil
}

var emptyVersionSpec = version.VersionSpec{}
var defaultVersionSpec = version.VersionSpec{
	Version: version.Version{
		Major: 0,
		Minor: 1,
	},
	Comparison: version.ComparisonGreaterOrEqual,
}

func fillDefaultVersionSpecs(config *Config) {
	for i := range config.Packages {
		p := &config.Packages[i]
		if p.VersionSpec == emptyVersionSpec {
			p.VersionSpec = defaultVersionSpec
		}
	}
}
