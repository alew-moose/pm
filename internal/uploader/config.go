package uploader

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name    string         `json:"name" yaml:"name"`
	Version PackageVersion `json:"ver" yaml:"ver"`
	Targets []Target       `json:"targets" yaml:"targets"`
}

// TODO XXX
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
	return &config, nil
}

func fromYAML(b []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

var packageNameRe = regexp.MustCompile(`^[\w-]+$`)

func (c *Config) Validate() error {
	if !packageNameRe.MatchString(c.Name) {
		return fmt.Errorf("invalid package name %q", c.Name)
	}
	if c.Version.Major == 0 && c.Version.Minor == 0 {
		return fmt.Errorf("invalid version %d.%d", c.Version.Major, c.Version.Minor)
	}
	for _, target := range c.Targets {
		// TODO: forbid absolute paths?
		if target.Path == "" {
			return errors.New("invalid target: empty path")
		}
	}
	return nil
}

type PackageVersion struct {
	Major uint64
	Minor uint64
}

func (v *PackageVersion) UnmarshalJSON(b []byte) error {
	var verStr string
	if err := json.Unmarshal(b, &verStr); err != nil {
		return fmt.Errorf("failed to parse version from %q: %s", b, err)
	}

	parts := strings.SplitN(verStr, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid version %q", verStr)
	}

	verMajor, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid major version %q", parts[0])
	}
	verMinor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid minor version %q", parts[1])
	}

	v.Major = verMajor
	v.Minor = verMinor

	return nil
}

func (v *PackageVersion) UnmarshalYAML(value *yaml.Node) error {
	parts := strings.SplitN(value.Value, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid version %q", value.Value)
	}

	verMajor, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid major version %q", parts[0])
	}
	verMinor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid minor version %q", parts[1])
	}

	v.Major = verMajor
	v.Minor = verMinor

	return nil
}

func (v PackageVersion) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

type Target struct {
	Path    string
	Exclude string
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
			return fmt.Errorf("failed to parse target %q: %s", b, err)
		}
	default:
		return fmt.Errorf("failed to parse target %q: unsupported type", b)
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
			return fmt.Errorf("failed to parse target: %s", err)
		}
		*t = target
	default:
		return fmt.Errorf("failed to parse target: unsupported kind %d", kind)
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
