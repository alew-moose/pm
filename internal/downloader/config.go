package downloader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alew-moose/pm/internal/version"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Packages []Package `json:"packages" yaml:"packages"`
}

type Package struct {
	Name        string              `json:"name" yaml:"name"`
	VersionSpec version.VersionSpec `json:"ver" yaml:"ver"`
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
			return nil, fmt.Errorf("failed to parse json: %s", err)
		}
	case ".yaml", ".yml":
		config, err = fromYAML(b)
		if err != nil {
			return nil, fmt.Errorf("failed to parse yaml: %s", err)
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
