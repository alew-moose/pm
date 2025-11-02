package uploader

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/alew-moose/pm/internal/downloader"
	"github.com/alew-moose/pm/internal/pkg"
	"github.com/alew-moose/pm/internal/version"
)

type Config struct {
	// TODO Name -> PackageName ?
	// TODO string -> PackageName
	Name         pkg.PackageName          `json:"name" yaml:"name"`
	Version      version.Version          `json:"ver" yaml:"ver"`
	Targets      []Target                 `json:"targets" yaml:"targets"`
	Dependencies []pkg.PackageVersionSpec `json:"packets" yaml:"packets"`
}

// TODO: rename (package name?)
func (c *Config) FileName() string {
	return fmt.Sprintf("%s-%s", c.Name, c.Version)
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

	downloader.FillDefaultVersionSpecs(config.Dependencies)
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

func (c *Config) Validate() error {
	if err := c.Name.Validate(); err != nil {
		return err
	}
	if err := c.Version.Validate(); err != nil {
		return err
	}
	for _, target := range c.Targets {
		if err := target.Validate(); err != nil {
			return err
		}
	}
	for _, dep := range c.Dependencies {
		if err := dep.Validate(); err != nil {
			return fmt.Errorf("invalid dependency: %s", err)
		}
	}
	return nil
}

type Target struct {
	Path    string
	Exclude string
}

func (t Target) Validate() error {
	// TODO: forbid absolute paths?
	if t.Path == "" {
		return errors.New("invalid target: empty path")
	}
	return nil
}

func (t *Target) FromMap(m map[string]any) error {
	path, ok := m["path"]
	if !ok {
		return errors.New("no path")
	}
	pathStr, ok := path.(string)
	if !ok {
		return errors.New("path is not a string")
	}
	t.Path = pathStr

	exclude, ok := m["exclude"]
	if ok {
		excludeStr, ok := exclude.(string)
		if !ok {
			return errors.New("exclude is not a string")
		}
		t.Exclude = excludeStr
	}

	return nil
}

func (t *Target) UnmarshalJSON(b []byte) error {
	var targetAny any
	if err := json.Unmarshal(b, &targetAny); err != nil {
		return err
	}
	switch target := targetAny.(type) {
	case string:
		t.Path = target
	case map[string]any:
		if err := t.FromMap(target); err != nil {
			return fmt.Errorf("parse target %q: %s", b, err)
		}
	default:
		return fmt.Errorf("parse target %q: unsupported type", b)
	}
	return nil
}

func (t *Target) UnmarshalYAML(node *yaml.Node) error {
	switch kind := node.Kind; kind {
	case yaml.ScalarNode:
		t.Path = node.Value
	case yaml.MappingNode:
		target, err := yamlNodesToTarget(node.Content)
		if err != nil {
			return fmt.Errorf("parse target: %s", err)
		}
		*t = target
	default:
		return fmt.Errorf("parse target: unsupported kind %d", kind)
	}
	return nil
}

func yamlNodesToTarget(nodes []*yaml.Node) (Target, error) {
	var target Target
	i := 0
	for i < len(nodes) {
		key := nodes[i].Value
		val := nodes[i+1].Value
		switch key {
		case "path":
			target.Path = val
		case "exclude":
			target.Exclude = val
		default:
			return target, fmt.Errorf("unknown field %q", key)
		}
		i += 2
	}
	return target, nil
}
