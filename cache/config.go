package cache

import (
	"fmt"
	"github.com/dradtke/go-allegro/allegro"
)

var _configs = make(map[string]*allegro.Config)

type ConfigNotFound struct {
	Key string
}

func (e *ConfigNotFound) Error() string {
	return fmt.Sprintf("config %s not found", e.Key)
}

// ClearConfigs() removes all configurations from the cache.
func ClearConfigs() {
	for key, val := range _configs {
		val.Destroy()
		delete(_configs, key)
	}
}

// LoadConfig() loads a configuration into the cache.
func LoadConfig(path, key string) error {
	cfg, err := allegro.LoadConfig(path)
	if err != nil {
		return err
	}
	if key == "" {
		key = path
	}
	_configs[key] = cfg
	return nil
}

// FindConfig() finds a configuration in the cache. If it
// doesn't exist, an error of type ConfigNotFound is returned.
func FindConfig(key string) (*allegro.Config, error) {
	if cfg, ok := _configs[key]; ok {
		return cfg, nil
	}
	return nil, &ConfigNotFound{key}
}

// Config() gets a configuration from the cache using FindConfig(),
// panicking if it isn't found.
func Config(key string) *allegro.Config {
	cfg, err := FindConfig(key)
	if err != nil {
		panic(err)
	}
	return cfg
}
