package sftp

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Host string
	Port string
	User string
	Path string // Path to packages dir
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
