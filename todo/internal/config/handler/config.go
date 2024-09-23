package config

import "github.com/BurntSushi/toml"

type PaginationConfig struct {
	DefaultLimit  int `toml:"default_limit"`
	DefaultOffset int `toml:"default_offset"`
}

type Config struct {
	Pagination PaginationConfig `toml:"pagination"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
