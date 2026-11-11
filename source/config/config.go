package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	World WorldConfig `json:"world"`
	Ant   AntConfig   `json:"ant"`
}

type WorldConfig struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type AntConfig struct {
	Count int `json:"count"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", path, err)
	}

	if cfg.World.Width <= 0 || cfg.World.Height <= 0 {
		return Config{}, fmt.Errorf("config %q: world.width and world.height must be positive", path)
	}
	if cfg.Ant.Count <= 0 {
		return Config{}, fmt.Errorf("config %q: ant.count must be positive", path)
	}

	return cfg, nil
}
