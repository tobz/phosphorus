package server

import "fmt"
import "strings"
import "github.com/tobz/phosphorus/log"
import "github.com/kylelemons/go-gypsy/yaml"
import "github.com/howeyc/fsnotify"

type Config struct {
	configPath string
	yamlConf   *yaml.File
	watcher    *fsnotify.Watcher
}

func NewConfig(path string) (*Config, error) {
	config := &Config{
		configPath: path,
		yamlConf:   nil,
		watcher:    nil,
	}

	// Make sure we can load the config.
	yamlConf, err := config.loadConfig(path)
	if err != nil {
		return nil, err
	}

	config.yamlConf = yamlConf

	// Start watching for changes.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Watch(path)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsModify() {
					conf, err := config.loadConfig(config.configPath)
					if err != nil {
						log.Server.Info("config", "Configuration change detected but unable to parse new configuration.")
						continue
					}

					config.yamlConf = conf
				}
			}
		}
	}()

	return config, nil
}

func (c *Config) loadConfig(path string) (*yaml.File, error) {
	file, err := yaml.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (c *Config) GetAsInteger(spec string) (int64, error) {
	yamlSpec := getYamlSpec(spec)

	// Make sure this is a single item.
	ok := c.assertScalar(yamlSpec)
	if !ok {
		return 0, fmt.Errorf("scalar requested, but spec points to list or doesn't exist")
	}

	return c.yamlConf.GetInt(yamlSpec)

}

func (c *Config) GetAsFloat(spec string) (float64, error) {
	return 0.0, nil
}

func (c *Config) GetAsString(spec string) (string, error) {
	yamlSpec := getYamlSpec(spec)

	// Make sure this is a single item.
	ok := c.assertScalar(yamlSpec)
	if !ok {
		return "", fmt.Errorf("scalar requested, but spec points to list or doesn't exist")
	}

	return c.yamlConf.Get(yamlSpec)
}

func (c *Config) GetAsManyIntegers(spec string) ([]int64, error) {
	return []int64{}, nil
}

func (c *Config) GetAsManyFloats(spec string) ([]float64, error) {
	return []float64{}, nil
}

func (c *Config) GetAsManyStrings(spec string) ([]string, error) {
	return []string{}, nil
}

func getYamlSpec(spec string) string {
	// We may eventually want a different syntax or more advanced stuff, but for now
	// I just want to use forward slashes instead of dots.
	return strings.Replace(spec, "/", ".", -1)
}

func (c *Config) assertScalar(spec string) bool {
	_, err := c.yamlConf.Count(spec)

	return err != nil
}

func (c *Config) assertList(spec string) bool {
	i, err := c.yamlConf.Count(spec)

	return i > 0 && err == nil
}
