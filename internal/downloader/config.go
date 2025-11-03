package downloader

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alew-moose/pm/internal/pkg"
	"github.com/alew-moose/pm/internal/version"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Packages []pkg.PackageVersionConstraint `json:"packages" yaml:"packages"`
}

func (c *Config) Validate() error {
	packages := make(map[pkg.PackageVersionConstraint]struct{}, len(c.Packages))
	for _, p := range c.Packages {
		if err := p.Validate(); err != nil {
			return err
		}
		if _, ok := packages[p]; ok {
			return fmt.Errorf("duplicate package %s", p)
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

	FillDefaultVersionConstraints(config.Packages)

	return config, nil
}

func fromJSON(b []byte) (*Config, error) {
	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func fromYAML(b []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// TODO: remove?
var emptyVersionConstraint = version.VersionConstraint{}
var defaultVersionConstraint = version.VersionConstraint{
	// TODO: fill?
}

// TODO: remove?
func FillDefaultVersionConstraints(packages []pkg.PackageVersionConstraint) {
	for i := range packages {
		p := &packages[i]
		if p.VersionConstraint == emptyVersionConstraint {
			log.Printf("using default version constraint %s for package %s\n", defaultVersionConstraint, p.Name)
			p.VersionConstraint = defaultVersionConstraint
		}
	}
}
