package sftp

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Host string
	Port string
	User string
	Path string // Path to packages dir
}

func (c *Config) Validate() error {
	if c.Host == "" {
		return errors.New("host is empty")
	}
	if c.Port == "" {
		return errors.New("port is empty")
	}
	if c.User == "" {
		return errors.New("user is empty")
	}
	if c.Path == "" {
		return errors.New("path is empty")
	}
	return nil
}

func ConfigFromFile(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := json.Unmarshal(b, &conf); err != nil {
		return nil, fmt.Errorf("unmarshal json: %s", err)
	}

	return &conf, nil
}
