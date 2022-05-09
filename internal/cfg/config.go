package cfg

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HistoryFile string
}

// FromFile returns a configuration object from a given file path and errors
// if the file does not exists of it the file contexts can't be unmarshaled.
// expected format is yaml.
func FromFile(path string) (*Config, error) {
	cfgfile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(cfgfile, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// ToFile marshals a configuration object to the given filepath using the yaml
// format
func ToFile(filePath string, cfg *Config) error {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	d := path.Dir(filePath)
	if _, err := os.Stat(d); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	if err := ioutil.WriteFile(filePath, b, 0644); err != nil {
		return err
	}
	return nil
}
