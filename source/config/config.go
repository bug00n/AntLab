package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	World WorldConfig `json:"world"`
	Ant   AntConfig   `json:"ant"`
	TPS   int         `json:"tps"`
}

type WorldConfig struct {
	Width  	 int `json:"width"`
	Height 	 int `json:"height"`
	GridSize int `json:"grid_size"`
}

type AntConfig struct {
	Count 	 int `json:"count"`
	Capacity int `json:"capacity"`
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
	if cfg.World.GridSize <= 0 {
		return Config{}, fmt.Errorf("config %q: world.grid_size must be positive", path)
	}

	if cfg.Ant.Count <= 0 {
		return Config{}, fmt.Errorf("config %q: ant.count must be positive", path)
	}
	if cfg.Ant.Capacity <= 0 {
		return Config{}, fmt.Errorf("config %q: ant.capacity must be positive", path)
	}
	if cfg.TPS <= 0 {
		return Config{}, fmt.Errorf("config %q: tps must be positive", path)
	}

	return cfg, nil
}
